package iceberg

import (
	"testing"

	"github.com/marmotdata/marmot/internal/core/connection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSource_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config with URI and credential",
			config: map[string]interface{}{
				"uri":        "http://localhost:8181",
				"warehouse":  "my-warehouse",
				"credential": "client-id:client-secret",
			},
			wantErr: false,
		},
		{
			name: "valid config with URI and token",
			config: map[string]interface{}{
				"uri":       "https://catalog.example.com",
				"warehouse": "my-warehouse",
				"token":     "my-bearer-token",
			},
			wantErr: false,
		},
		{
			name: "missing URI",
			config: map[string]interface{}{
				"warehouse": "my-warehouse",
				"credential": "client-id:client-secret",
			},
			wantErr: true,
			errMsg:  "must be a valid URL",
		},
		{
			name: "missing warehouse",
			config: map[string]interface{}{
				"uri":        "http://localhost:8181",
				"credential": "client-id:client-secret",
			},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "invalid URI format",
			config: map[string]interface{}{
				"uri":        "not-a-valid-url",
				"warehouse":  "my-warehouse",
				"credential": "client-id:client-secret",
			},
			wantErr: true,
			errMsg:  "must be a valid URL",
		},
		{
			name: "URI without scheme",
			config: map[string]interface{}{
				"uri":        "localhost:8181",
				"warehouse":  "my-warehouse",
				"credential": "client-id:client-secret",
			},
			wantErr: true,
			errMsg:  "must be a valid URL",
		},
		{
			name: "missing authentication - no credential or token",
			config: map[string]interface{}{
				"uri":       "http://localhost:8181",
				"warehouse": "my-warehouse",
			},
			wantErr: true,
			errMsg:  "authentication required",
		},
		{
			name: "both credential and token provided",
			config: map[string]interface{}{
				"uri":        "http://localhost:8181",
				"warehouse":  "my-warehouse",
				"credential": "client-id:client-secret",
				"token":      "bearer-token",
			},
			wantErr: false,
		},
		{
			name: "empty config",
			config: map[string]interface{}{},
			wantErr: true,
			errMsg:  "must be a valid URL",
		},
		{
			name: "valid with all optional fields",
			config: map[string]interface{}{
				"uri":        "https://catalog.example.com",
				"warehouse":  "prod-warehouse",
				"credential": "client-id:client-secret",
				"properties": map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
				"prefix": "v1",
			},
			wantErr: false,
		},
		{
			name: "valid HTTP URI",
			config: map[string]interface{}{
				"uri":        "http://192.168.1.1:8181",
				"warehouse":  "my-warehouse",
				"credential": "client-id:client-secret",
			},
			wantErr: false,
		},
		{
			name: "valid HTTPS URI",
			config: map[string]interface{}{
				"uri":       "https://secure-catalog.example.com",
				"warehouse": "my-warehouse",
				"token":     "my-token",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Source{}
			err := s.Validate(tt.config)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestIcebergRESTConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *IcebergRESTConfig
		valid  bool
	}{
		{
			name: "valid with credential",
			config: &IcebergRESTConfig{
				URI:        "http://localhost:8181",
				Warehouse:  "my-warehouse",
				Credential: "client-id:client-secret",
			},
			valid: true,
		},
		{
			name: "valid with token",
			config: &IcebergRESTConfig{
				URI:       "http://localhost:8181",
				Warehouse: "my-warehouse",
				Token:     "my-token",
			},
			valid: true,
		},
		{
			name: "missing URI",
			config: &IcebergRESTConfig{
				Warehouse:  "my-warehouse",
				Credential: "client-id:client-secret",
			},
			valid: false,
		},
		{
			name: "missing warehouse",
			config: &IcebergRESTConfig{
				URI:        "http://localhost:8181",
				Credential: "client-id:client-secret",
			},
			valid: false,
		},
		{
			name: "invalid URI",
			config: &IcebergRESTConfig{
				URI:        "not-a-url",
				Warehouse:  "my-warehouse",
				Credential: "client-id:client-secret",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := connection.ValidateConfig(tt.config)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
