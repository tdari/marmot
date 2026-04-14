package trino

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
			name: "valid config with basic auth",
			config: map[string]interface{}{
				"host":     "trino.example.com",
				"user":     "marmot_reader",
				"password": "secret",
			},
			wantErr: false,
		},
		{
			name: "valid config with access token",
			config: map[string]interface{}{
				"host":         "trino.example.com",
				"user":         "marmot_reader",
				"access_token": "jwt-token",
			},
			wantErr: false,
		},
		{
			name: "missing host",
			config: map[string]interface{}{
				"user": "marmot_reader",
			},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "missing user",
			config: map[string]interface{}{
				"host": "localhost",
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
			name: "valid with port",
			config: map[string]interface{}{
				"host":     "trino.example.com",
				"port":     8080,
				"user":     "marmot_reader",
			},
			wantErr: false,
		},
		{
			name: "invalid port - too high",
			config: map[string]interface{}{
				"host":     "trino.example.com",
				"port":     65536,
				"user":     "marmot_reader",
			},
			wantErr: true,
			errMsg:  "must be at most",
		},
		{
			name: "valid with secure connection",
			config: map[string]interface{}{
				"host":     "trino.example.com",
				"user":     "marmot_reader",
				"secure":   true,
				"password": "secret",
			},
			wantErr: false,
		},
		{
			name: "valid with SSL cert path",
			config: map[string]interface{}{
				"host":           "trino.example.com",
				"user":           "marmot_reader",
				"secure":         true,
				"ssl_cert_path":  "/path/to/cert.pem",
				"access_token":   "jwt-token",
			},
			wantErr: false,
		},
		{
			name: "port defaults to 8080",
			config: map[string]interface{}{
				"host": "trino.example.com",
				"user": "marmot_reader",
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

func TestTrinoConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *TrinoConfig
		valid  bool
	}{
		{
			name: "valid minimal config",
			config: &TrinoConfig{
				Host: "localhost",
				User: "marmot_reader",
			},
			valid: true,
		},
		{
			name: "valid with password",
			config: &TrinoConfig{
				Host:     "localhost",
				User:     "marmot_reader",
				Password: "secret",
			},
			valid: true,
		},
		{
			name: "valid with access token",
			config: &TrinoConfig{
				Host:        "localhost",
				User:        "marmot_reader",
				AccessToken: "jwt-token",
			},
			valid: true,
		},
		{
			name: "missing host",
			config: &TrinoConfig{
				User: "marmot_reader",
			},
			valid: false,
		},
		{
			name: "missing user",
			config: &TrinoConfig{
				Host: "localhost",
			},
			valid: false,
		},
		{
			name: "invalid port",
			config: &TrinoConfig{
				Host: "localhost",
				Port: 99999,
				User: "marmot_reader",
			},
			valid: false,
		},
		{
			name: "valid with secure",
			config: &TrinoConfig{
				Host:   "localhost",
				Port:   8443,
				User:   "marmot_reader",
				Secure: true,
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
