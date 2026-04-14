package connections

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/marmotdata/marmot/internal/api/v1/common"
	"github.com/marmotdata/marmot/internal/core/connection"
	"github.com/marmotdata/marmot/internal/core/user"
	"github.com/marmotdata/marmot/internal/plugin"
	"github.com/rs/zerolog/log"
)

// CreateRequest represents a request to create a connection
type CreateRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Type        string                 `json:"type" validate:"required"`
	Description *string                `json:"description"`
	Config      map[string]interface{} `json:"config" validate:"required"` // Configuration as flat map - sensitive fields encrypted based on type's ConfigSpec
	Tags        []string               `json:"tags"`
}

// UpdateRequest represents a request to update a connection
type UpdateRequest struct {
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	Config      map[string]interface{} `json:"config"` // Configuration as flat map - sensitive fields encrypted based on type's ConfigSpec
	Tags        []string               `json:"tags"`
}

// ConnectionResponse represents a connection in API responses (with credentials redacted)
type ConnectionResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Description *string                `json:"description,omitempty"`
	Config      map[string]interface{} `json:"config"` // Config with sensitive values redacted
	Tags        []string               `json:"tags,omitempty"`
	CreatedBy   string                 `json:"created_by"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

// ListResponse represents a paginated list of connections
type ListResponse struct {
	Connections []*ConnectionResponse `json:"connections"`
	Total       int                   `json:"total"`
	Limit       int                   `json:"limit"`
	Offset      int                   `json:"offset"`
}

// @Summary Create a new connection
// @Description Create a new data source connection
// @Tags connections
// @Accept json
// @Produce json
// @Param connection body CreateRequest true "Connection creation request"
// @Success 201 {object} ConnectionResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 409 {object} common.ErrorResponse
// @Router /connections [post]
func (h *Handler) createConnection(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Name == "" {
		common.RespondError(w, http.StatusBadRequest, "Connection name is required")
		return
	}
	if req.Type == "" {
		common.RespondError(w, http.StatusBadRequest, "Connection type is required")
		return
	}

	// Get the connection type's config spec from the registry
	typeMeta, err := connection.GetRegistry().GetMeta(req.Type)
	if err != nil {
		common.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid connection type: %s", req.Type))
		return
	}

	// Validate config against ConfigSpec
	if err := connection.ValidateConfigAgainstSpec(req.Config, typeMeta.ConfigSpec); err != nil {
		common.RespondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid configuration: %s", err.Error()))
		return
	}

	usr, ok := r.Context().Value(common.UserContextKey).(*user.User)
	if !ok {
		common.RespondError(w, http.StatusUnauthorized, "User context required")
		return
	}

	input := connection.CreateInput{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		Config:      req.Config,
		Tags:        req.Tags,
		CreatedBy:   usr.Name,
	}

	conn, err := h.connectionService.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, connection.ErrAlreadyExists) {
			common.RespondError(w, http.StatusConflict, "Connection with this name already exists")
		} else if errors.Is(err, connection.ErrInvalidInput) {
			common.RespondError(w, http.StatusBadRequest, err.Error())
		} else {
			log.Error().Err(err).Msg("Failed to create connection")
			common.RespondError(w, http.StatusInternalServerError, "Failed to create connection")
		}
		return
	}

	response := toConnectionResponse(conn, h.registry)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// @Summary Get a connection
// @Description Get a connection by ID
// @Tags connections
// @Produce json
// @Param id path string true "Connection ID"
// @Success 200 {object} ConnectionResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /connections/{id} [get]
func (h *Handler) getConnection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		common.RespondError(w, http.StatusBadRequest, "Connection ID is required")
		return
	}

	conn, err := h.connectionService.Get(r.Context(), id)
	if err != nil {
		if err == connection.ErrNotFound {
			common.RespondError(w, http.StatusNotFound, "Connection not found")
		} else {
			log.Error().Err(err).Str("id", id).Msg("Failed to get connection")
			common.RespondError(w, http.StatusInternalServerError, "Failed to get connection")
		}
		return
	}

	response := toConnectionResponse(conn, h.registry)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary List connections
// @Description List all connections with optional filtering
// @Tags connections
// @Produce json
// @Param type query string false "Filter by plugin type"
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} ListResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /connections [get]
func (h *Handler) listConnections(w http.ResponseWriter, r *http.Request) {
	opts := &connection.ListOptions{}

	if connType := r.URL.Query().Get("type"); connType != "" {
		opts.Type = connType
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit := 0
		if _, err := fmt.Sscanf(limitStr, "%d", &limit); err == nil && limit > 0 {
			opts.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset := 0
		if _, err := fmt.Sscanf(offsetStr, "%d", &offset); err == nil && offset >= 0 {
			opts.Offset = offset
		}
	}

	if query := r.URL.Query().Get("query"); query != "" {
		opts.Search = query
	}

	connections, total, err := h.connectionService.List(r.Context(), opts)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list connections")
		common.RespondError(w, http.StatusInternalServerError, "Failed to list connections")
		return
	}

	responses := make([]*ConnectionResponse, len(connections))
	for i, conn := range connections {
		responses[i] = toConnectionResponse(conn, h.registry)
	}

	response := ListResponse{
		Connections: responses,
		Total:       total,
		Limit:       opts.Limit,
		Offset:      opts.Offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary List available connection types
// @Description Get all available connection types with their configuration specifications
// @Tags connections
// @Produce json
// @Success 200 {array} connection.ConnectionTypeMeta
// @Failure 500 {object} common.ErrorResponse
// @Router /connections/types [get]
func (h *Handler) listConnectionTypes(w http.ResponseWriter, r *http.Request) {
	types := connection.GetRegistry().List()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(types); err != nil {
		log.Error().Err(err).Msg("Failed to encode connection types")
		common.RespondError(w, http.StatusInternalServerError, "Failed to encode response")
		return
	}
}

// @Summary Update a connection
// @Description Update an existing connection
// @Tags connections
// @Accept json
// @Produce json
// @Param id path string true "Connection ID"
// @Param connection body UpdateRequest true "Connection update request"
// @Success 200 {object} ConnectionResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /connections/{id} [put]
func (h *Handler) updateConnection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		common.RespondError(w, http.StatusBadRequest, "Connection ID is required")
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate that name is not empty if provided
	if req.Name != nil && *req.Name == "" {
		common.RespondError(w, http.StatusBadRequest, "Connection name is required")
		return
	}

	input := connection.UpdateInput{
		Name:        req.Name,
		Description: req.Description,
		Config:      req.Config,
		Tags:        req.Tags,
	}

	conn, err := h.connectionService.Update(r.Context(), id, input)
	if err != nil {
		if errors.Is(err, connection.ErrNotFound) {
			common.RespondError(w, http.StatusNotFound, "Connection not found")
		} else if errors.Is(err, connection.ErrInvalidInput) {
			common.RespondError(w, http.StatusBadRequest, err.Error())
		} else if errors.Is(err, connection.ErrAlreadyExists) {
			common.RespondError(w, http.StatusConflict, "Connection with this name already exists")
		} else {
			log.Error().Err(err).Str("id", id).Msg("Failed to update connection")
			common.RespondError(w, http.StatusInternalServerError, "Failed to update connection")
		}
		return
	}

	response := toConnectionResponse(conn, h.registry)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Delete a connection
// @Description Delete a connection by ID, optionally tearing down associated schedules
// @Tags connections
// @Param id path string true "Connection ID"
// @Param teardown query boolean false "Delete all schedules using this connection"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /connections/{id} [delete]
func (h *Handler) deleteConnection(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		common.RespondError(w, http.StatusBadRequest, "Connection ID is required")
		return
	}

	// Check if user wants to teardown all schedules using this connection
	teardown := r.URL.Query().Get("teardown") == "true"

	opts := &connection.DeleteOptions{
		TeardownSchedules: teardown,
	}

	err := h.connectionService.DeleteWithOptions(r.Context(), id, opts)
	if err != nil {
		if err == connection.ErrNotFound {
			common.RespondError(w, http.StatusNotFound, "Connection not found")
		} else {
			log.Error().Err(err).Str("id", id).Msg("Failed to delete connection")
			common.RespondError(w, http.StatusInternalServerError, "Failed to delete connection")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// toConnectionResponse converts a connection to API response, redacting sensitive fields
func toConnectionResponse(conn *connection.Connection, registry *plugin.Registry) *ConnectionResponse {
	// Redact sensitive fields in config
	config := redactSensitiveFields(conn.Config, conn.Type, registry)

	return &ConnectionResponse{
		ID:          conn.ID,
		Name:        conn.Name,
		Type:        conn.Type,
		Description: conn.Description,
		Config:      config,
		Tags:        conn.Tags,
		CreatedBy:   conn.CreatedBy,
		CreatedAt:   conn.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   conn.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// redactSensitiveFields replaces sensitive field values with redacted placeholders
func redactSensitiveFields(config map[string]interface{}, pluginType string, registry *plugin.Registry) map[string]interface{} {
	redacted := make(map[string]interface{})
	for k, v := range config {
		redacted[k] = v
	}

	// Get plugin metadata
	entry, err := registry.Get(pluginType)
	if err != nil {
		return redacted
	}

	// Identify sensitive fields
	for _, field := range entry.Meta.ConfigSpec {
		if field.Sensitive {
			if _, exists := redacted[field.Name]; exists {
				redacted[field.Name] = "***encrypted***"
			}
		}
	}

	return redacted
}
