package airflow

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
			name: "valid config with username and password",
			config: map[string]interface{}{
				"host":     "http://localhost:8080",
				"username": "admin",
				"password": "secret",
			},
			wantErr: false,
		},
		{
			name: "valid config with API token",
			config: map[string]interface{}{
				"host":      "http://localhost:8080",
				"api_token": "my-api-token",
			},
			wantErr: false,
		},
		{
			name: "missing host",
			config: map[string]interface{}{
				"username": "admin",
				"password": "secret",
			},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "invalid host - not a URL",
			config: map[string]interface{}{
				"host":     "not-a-url",
				"username": "admin",
				"password": "secret",
			},
			wantErr: true,
			errMsg:  "must be a valid URL",
		},
		{
			name: "no authentication method",
			config: map[string]interface{}{
				"host": "http://localhost:8080",
			},
			wantErr: true,
			errMsg:  "authentication required",
		},
		{
			name: "empty config",
			config: map[string]interface{}{},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "valid with trailing slash in URL",
			config: map[string]interface{}{
				"host":     "http://localhost:8080/",
				"username": "admin",
				"password": "secret",
			},
			wantErr: false, // Trailing slash is trimmed
		},
		{
			name: "valid with username only",
			config: map[string]interface{}{
				"host":     "http://localhost:8080",
				"username": "admin",
			},
			wantErr: false,
		},
		{
			name: "valid with both auth methods",
			config: map[string]interface{}{
				"host":      "http://localhost:8080",
				"username":  "admin",
				"password":  "secret",
				"api_token": "my-api-token",
			},
			wantErr: false,
		},
		{
			name: "valid HTTPS URL",
			config: map[string]interface{}{
				"host":     "https://airflow.example.com",
				"username": "admin",
				"password": "secret",
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

func TestAirflowConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *AirflowConfig
		valid  bool
	}{
		{
			name: "valid with basic auth",
			config: &AirflowConfig{
				Host:     "http://localhost:8080",
				Username: "admin",
				Password: "secret",
			},
			valid: true,
		},
		{
			name: "valid with API token",
			config: &AirflowConfig{
				Host:     "http://localhost:8080",
				APIToken: "my-token",
			},
			valid: true,
		},
		{
			name: "missing host",
			config: &AirflowConfig{
				Username: "admin",
			},
			valid: false,
		},
		{
			name: "invalid host",
			config: &AirflowConfig{
				Host:     "not-a-url",
				Username: "admin",
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
