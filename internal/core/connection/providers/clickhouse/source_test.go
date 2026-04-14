package clickhouse

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
				"host":     "clickhouse.example.com",
				"port":     9000,
				"user":     "default",
				"password": "password123",
				"database": "mydb",
				"secure":   true,
			},
			wantErr: false,
		},
		{
			name: "valid config with required fields only",
			config: map[string]interface{}{
				"host":     "localhost",
				"user":     "default",
				"database": "default",
			},
			wantErr: false,
		},
		{
			name: "missing host",
			config: map[string]interface{}{
				"user":     "default",
				"database": "default",
			},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "missing user",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "default",
			},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "missing database - defaults to 'default'",
			config: map[string]interface{}{
				"host": "localhost",
				"user": "default",
			},
			wantErr: false, // Database defaults to "default"
		},
		{
			name: "empty config",
			config: map[string]interface{}{},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "invalid port - too high",
			config: map[string]interface{}{
				"host":     "localhost",
				"user":     "default",
				"database": "default",
				"port":     65536,
			},
			wantErr: true,
			errMsg:  "must be at most",
		},
		{
			name: "valid with secure connection",
			config: map[string]interface{}{
				"host":     "secure.clickhouse.com",
				"user":     "admin",
				"database": "metrics",
				"secure":   true,
			},
			wantErr: false,
		},
		{
			name: "port defaults to 9000",
			config: map[string]interface{}{
				"host":     "localhost",
				"user":     "default",
				"database": "default",
			},
			wantErr: false,
		},
		{
			name: "database defaults to default",
			config: map[string]interface{}{
				"host": "localhost",
				"user": "default",
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

func TestClickHouseConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *ClickHouseConfig
		valid  bool
	}{
		{
			name: "valid config",
			config: &ClickHouseConfig{
				Host:     "localhost",
				Port:     9000,
				User:     "default",
				Database: "default",
			},
			valid: true,
		},
		{
			name: "missing host",
			config: &ClickHouseConfig{
				Port:     9000,
				User:     "default",
				Database: "default",
			},
			valid: false,
		},
		{
			name: "missing user",
			config: &ClickHouseConfig{
				Host:     "localhost",
				Database: "default",
			},
			valid: false,
		},
		{
			name: "missing database",
			config: &ClickHouseConfig{
				Host: "localhost",
				User: "default",
			},
			valid: false,
		},
		{
			name: "invalid port - too high",
			config: &ClickHouseConfig{
				Host:     "localhost",
				User:     "default",
				Database: "default",
				Port:     99999,
			},
			valid: false,
		},
		{
			name: "valid with password",
			config: &ClickHouseConfig{
				Host:     "localhost",
				Port:     9000,
				User:     "default",
				Database: "default",
				Password: "secret",
			},
			valid: true,
		},
		{
			name: "valid with secure",
			config: &ClickHouseConfig{
				Host:     "localhost",
				Port:     9000,
				User:     "default",
				Database: "default",
				Secure:   true,
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
