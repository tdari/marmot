package connection

import (
	"time"
)

type Connection struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"` // Connection type (postgresql, mysql, s3, bigquery, snowflake, etc.)
	Description *string                `json:"description,omitempty"`
	Config      map[string]interface{} `json:"config"` // Configuration as flat map - sensitive fields encrypted based on type's ConfigSpec
	Tags        []string               `json:"tags,omitempty"`
	CreatedBy   string                 `json:"created_by"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type CreateInput struct {
	Name        string                 `json:"name" validate:"required"`
	Type        string                 `json:"type" validate:"required"`
	Description *string                `json:"description"`
	Config      map[string]interface{} `json:"config"`
	Tags        []string               `json:"tags"`
	CreatedBy   string                 `json:"created_by" validate:"required"`
}

type UpdateInput struct {
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	Config      map[string]interface{} `json:"config"`
	Tags        []string               `json:"tags"`
}
