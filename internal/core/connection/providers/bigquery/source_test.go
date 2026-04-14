package bigquery

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
			name: "valid config with credentials path",
			config: map[string]interface{}{
				"project_id":         "my-gcp-project",
				"credentials_path":   "/path/to/key.json",
				"dataset":            "my_dataset",
				"location":           "US",
			},
			wantErr: false,
		},
		{
			name: "valid config with credentials JSON",
			config: map[string]interface{}{
				"project_id":       "my-gcp-project",
				"credentials_json": `{"type":"service_account"}`,
				"dataset":          "my_dataset",
			},
			wantErr: false,
		},
		{
			name: "valid config with use_default_credentials",
			config: map[string]interface{}{
				"project_id":             "my-gcp-project",
				"use_default_credentials": true,
				"dataset":                "my_dataset",
			},
			wantErr: false,
		},
		{
			name: "missing project_id",
			config: map[string]interface{}{
				"credentials_path": "/path/to/key.json",
			},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "no authentication method",
			config: map[string]interface{}{
				"project_id": "my-gcp-project",
				"dataset":    "my_dataset",
			},
			wantErr: true,
			errMsg:  "at least one authentication method",
		},
		{
			name: "multiple authentication methods",
			config: map[string]interface{}{
				"project_id":              "my-gcp-project",
				"credentials_path":        "/path/to/key.json",
				"credentials_json":        `{"type":"service_account"}`,
				"use_default_credentials": true,
			},
			wantErr: true,
			errMsg:  "only one authentication method",
		},
		{
			name: "credentials_path and credentials_json both provided",
			config: map[string]interface{}{
				"project_id":        "my-gcp-project",
				"credentials_path":  "/path/to/key.json",
				"credentials_json":  `{"type":"service_account"}`,
			},
			wantErr: true,
			errMsg:  "only one authentication method",
		},
		{
			name: "credentials_path and use_default_credentials both provided",
			config: map[string]interface{}{
				"project_id":              "my-gcp-project",
				"credentials_path":        "/path/to/key.json",
				"use_default_credentials": true,
			},
			wantErr: true,
			errMsg:  "only one authentication method",
		},
		{
			name: "empty config",
			config: map[string]interface{}{},
			wantErr: true,
			errMsg:  "at least one authentication method",
		},
		{
			name: "valid config with all optional fields",
			config: map[string]interface{}{
				"project_id":       "my-gcp-project",
				"credentials_json": `{"type":"service_account"}`,
				"dataset":          "my_dataset",
				"location":         "europe-west1",
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

func TestBigQueryConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *BigQueryConfig
		valid  bool
	}{
		{
			name: "valid config with credentials path",
			config: &BigQueryConfig{
				ProjectID:       "my-project",
				CredentialsPath: "/path/to/key.json",
			},
			valid: true,
		},
		{
			name: "valid config with credentials JSON",
			config: &BigQueryConfig{
				ProjectID:       "my-project",
				CredentialsJSON: `{"type":"service_account"}`,
			},
			valid: true,
		},
		{
			name: "valid config with use default credentials",
			config: &BigQueryConfig{
				ProjectID:             "my-project",
				UseDefaultCredentials: true,
			},
			valid: true,
		},
		{
			name: "missing project_id",
			config: &BigQueryConfig{
				CredentialsPath: "/path/to/key.json",
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
