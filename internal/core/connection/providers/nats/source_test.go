package nats

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
			name: "valid config with all fields",
			config: map[string]interface{}{
				"host":             "nats.example.com",
				"port":             4222,
				"token":            "secret-token",
				"username":         "user",
				"password":         "pass",
				"credentials_file": "/path/to/creds.creds",
				"tls":              true,
				"tls_insecure":     true,
			},
			wantErr: false,
		},
		{
			name: "valid config with required field only",
			config: map[string]interface{}{
				"host": "localhost",
			},
			wantErr: false,
		},
		{
			name: "empty config",
			config: map[string]interface{}{},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "missing host",
			config: map[string]interface{}{
				"port": 4222,
			},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "invalid port - too high",
			config: map[string]interface{}{
				"host": "localhost",
				"port": 65536,
			},
			wantErr: true,
			errMsg:  "must be at most",
		},
		{
			name: "valid config with token auth",
			config: map[string]interface{}{
				"host":  "localhost",
				"token": "my-token",
			},
			wantErr: false,
		},
		{
			name: "valid config with user/password auth",
			config: map[string]interface{}{
				"host":     "localhost",
				"username": "user",
				"password": "pass",
			},
			wantErr: false,
		},
		{
			name: "valid config with credentials file",
			config: map[string]interface{}{
				"host":             "localhost",
				"credentials_file": "/etc/nats/creds.creds",
			},
			wantErr: false,
		},
		{
			name: "valid config with TLS",
			config: map[string]interface{}{
				"host":         "localhost",
				"tls":          true,
				"tls_insecure": false,
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

func TestNatsConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *NatsConfig
		valid  bool
	}{
		{
			name: "valid minimal config",
			config: &NatsConfig{
				Host: "localhost",
			},
			valid: true,
		},
		{
			name: "valid with port",
			config: &NatsConfig{
				Host: "localhost",
				Port: 4222,
			},
			valid: true,
		},
		{
			name: "missing host",
			config: &NatsConfig{
				Port: 4222,
			},
			valid: false,
		},
		{
			name: "invalid port - too high",
			config: &NatsConfig{
				Host: "localhost",
				Port: 99999,
			},
			valid: false,
		},
		{
			name: "valid with all auth fields",
			config: &NatsConfig{
				Host:            "localhost",
				Token:           "token",
				Username:        "user",
				Password:        "pass",
				CredentialsFile: "/path/to/creds",
			},
			valid: true,
		},
		{
			name: "valid with TLS",
			config: &NatsConfig{
				Host: "localhost",
				TLS:  true,
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
