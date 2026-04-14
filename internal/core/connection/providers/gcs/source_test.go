package gcs

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
			name: "valid config with project ID and credentials file",
			config: map[string]interface{}{
				"project_id":       "my-gcp-project",
				"credentials_file": "/path/to/key.json",
				"bucket_name":      "my-bucket",
			},
			wantErr: false,
		},
		{
			name: "valid config with project ID and credentials JSON",
			config: map[string]interface{}{
				"project_id":       "my-gcp-project",
				"credentials_json": `{"type":"service_account"}`,
				"bucket_name":      "my-bucket",
			},
			wantErr: false,
		},
		{
			name: "valid config with project ID only",
			config: map[string]interface{}{
				"project_id": "my-gcp-project",
			},
			wantErr: false,
		},
		{
			name: "missing project ID",
			config: map[string]interface{}{
				"credentials_file": "/path/to/key.json",
				"bucket_name":      "my-bucket",
			},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "empty config",
			config: map[string]interface{}{},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "valid with custom endpoint",
			config: map[string]interface{}{
				"project_id": "my-gcp-project",
				"endpoint":   "https://storage.googleapis.com",
			},
			wantErr: false,
		},
		{
			name: "valid with disable auth",
			config: map[string]interface{}{
				"project_id":    "my-gcp-project",
				"disable_auth":  true,
				"bucket_name":   "public-bucket",
			},
			wantErr: false,
		},
		{
			name: "valid with all fields",
			config: map[string]interface{}{
				"project_id":       "my-gcp-project",
				"bucket_name":      "my-bucket",
				"credentials_file": "/path/to/key.json",
				"endpoint":         "https://storage.googleapis.com",
				"disable_auth":     false,
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

func TestGCSConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *GCSConfig
		valid  bool
	}{
		{
			name: "valid config",
			config: &GCSConfig{
				ProjectID:       "my-project",
				CredentialsFile: "/path/to/key.json",
			},
			valid: true,
		},
		{
			name: "valid with credentials JSON",
			config: &GCSConfig{
				ProjectID:       "my-project",
				CredentialsJSON: `{"type":"service_account"}`,
			},
			valid: true,
		},
		{
			name: "valid with bucket name",
			config: &GCSConfig{
				ProjectID:  "my-project",
				BucketName: "my-bucket",
			},
			valid: true,
		},
		{
			name: "missing project ID",
			config: &GCSConfig{
				CredentialsFile: "/path/to/key.json",
			},
			valid: false,
		},
		{
			name: "valid with disable auth",
			config: &GCSConfig{
				ProjectID:   "my-project",
				DisableAuth: true,
			},
			valid: true,
		},
		{
			name: "valid with custom endpoint",
			config: &GCSConfig{
				ProjectID:  "my-project",
				Endpoint:   "https://storage.googleapis.com",
			},
			valid: true,
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
