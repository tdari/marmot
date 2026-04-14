package connections

import (
	"net/http"

	"github.com/marmotdata/marmot/internal/api/v1/common"
	"github.com/marmotdata/marmot/internal/config"
	authService "github.com/marmotdata/marmot/internal/core/auth"
	"github.com/marmotdata/marmot/internal/core/connection"
	userService "github.com/marmotdata/marmot/internal/core/user"
	"github.com/marmotdata/marmot/internal/plugin"
)

// Handler handles connection-related HTTP requests
type Handler struct {
	connectionService connection.Service
	userService       userService.Service
	authService       authService.Service
	registry          *plugin.Registry
	config            *config.Config
}

// NewHandler creates a new connection handler
func NewHandler(connectionService connection.Service, userService userService.Service, authService authService.Service, registry *plugin.Registry, cfg *config.Config) *Handler {
	return &Handler{
		connectionService: connectionService,
		userService:       userService,
		authService:       authService,
		registry:          registry,
		config:            cfg,
	}
}

// Routes returns the routes for the connections handler
func (h *Handler) Routes() []common.Route {
	return []common.Route{
		{
			Path:    "/api/v1/connections",
			Method:  http.MethodPost,
			Handler: h.createConnection,
			Middleware: []func(http.HandlerFunc) http.HandlerFunc{
				common.WithAuth(h.userService, h.authService, h.config),
				common.RequirePermission(h.userService, "connections", "manage"),
			},
		},
		{
			Path:    "/api/v1/connections/{id}",
			Method:  http.MethodGet,
			Handler: h.getConnection,
			Middleware: []func(http.HandlerFunc) http.HandlerFunc{
				common.WithAuth(h.userService, h.authService, h.config),
				common.RequirePermission(h.userService, "connections", "view"),
			},
		},
		{
			Path:    "/api/v1/connections/types",
			Method:  http.MethodGet,
			Handler: h.listConnectionTypes,
			Middleware: []func(http.HandlerFunc) http.HandlerFunc{
				common.WithAuth(h.userService, h.authService, h.config),
				common.RequirePermission(h.userService, "connections", "view"),
			},
		},
		{
			Path:    "/api/v1/connections",
			Method:  http.MethodGet,
			Handler: h.listConnections,
			Middleware: []func(http.HandlerFunc) http.HandlerFunc{
				common.WithAuth(h.userService, h.authService, h.config),
				common.RequirePermission(h.userService, "connections", "view"),
			},
		},
		{
			Path:    "/api/v1/connections/{id}",
			Method:  http.MethodPut,
			Handler: h.updateConnection,
			Middleware: []func(http.HandlerFunc) http.HandlerFunc{
				common.WithAuth(h.userService, h.authService, h.config),
				common.RequirePermission(h.userService, "connections", "manage"),
			},
		},
		{
			Path:    "/api/v1/connections/{id}",
			Method:  http.MethodDelete,
			Handler: h.deleteConnection,
			Middleware: []func(http.HandlerFunc) http.HandlerFunc{
				common.WithAuth(h.userService, h.authService, h.config),
				common.RequirePermission(h.userService, "connections", "manage"),
			},
		},
	}
}
