package postgresql

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
				"host":     "postgres.example.com",
				"port":     5432,
				"database": "mydb",
				"user":     "postgres",
				"password": "secret",
				"ssl_mode": "require",
			},
			wantErr: false,
		},
		{
			name: "valid config with required fields",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "postgres",
				"password": "secret",
			},
			wantErr: false,
		},
		{
			name: "missing host",
			config: map[string]interface{}{
				"database": "mydb",
				"user":     "postgres",
				"password": "secret",
			},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "missing database",
			config: map[string]interface{}{
				"host":     "localhost",
				"user":     "postgres",
				"password": "secret",
			},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "missing user",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"password": "secret",
			},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "missing password",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "postgres",
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
			name: "invalid port - too high",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "postgres",
				"password": "secret",
				"port":     65536,
			},
			wantErr: true,
			errMsg:  "must be at most",
		},
		{
			name: "invalid ssl_mode",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "postgres",
				"password": "secret",
				"ssl_mode": "invalid",
			},
			wantErr: true,
			errMsg:  "must be one of",
		},
		{
			name: "valid ssl_mode - disable",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "postgres",
				"password": "secret",
				"ssl_mode": "disable",
			},
			wantErr: false,
		},
		{
			name: "valid ssl_mode - require",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "postgres",
				"password": "secret",
				"ssl_mode": "require",
			},
			wantErr: false,
		},
		{
			name: "port defaults to 5432",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "postgres",
				"password": "secret",
			},
			wantErr: false,
		},
		{
			name: "ssl_mode defaults to disable",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "postgres",
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

func TestPostgreSQLConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *PostgreSQLConfig
		valid  bool
	}{
		{
			name: "valid config",
			config: &PostgreSQLConfig{
				Host:     "localhost",
				Port:     5432,
				Database: "mydb",
				User:     "postgres",
				Password: "secret",
				SSLMode:  "disable",
			},
			valid: true,
		},
		{
			name: "missing host",
			config: &PostgreSQLConfig{
				Database: "mydb",
				User:     "postgres",
				Password: "secret",
			},
			valid: false,
		},
		{
			name: "missing database",
			config: &PostgreSQLConfig{
				Host:     "localhost",
				User:     "postgres",
				Password: "secret",
			},
			valid: false,
		},
		{
			name: "missing user",
			config: &PostgreSQLConfig{
				Host:     "localhost",
				Database: "mydb",
				Password: "secret",
			},
			valid: false,
		},
		{
			name: "missing password",
			config: &PostgreSQLConfig{
				Host:     "localhost",
				Database: "mydb",
				User:     "postgres",
			},
			valid: false,
		},
		{
			name: "invalid port",
			config: &PostgreSQLConfig{
				Host:     "localhost",
				Port:     99999,
				Database: "mydb",
				User:     "postgres",
				Password: "secret",
			},
			valid: false,
		},
		{
			name: "invalid ssl_mode",
			config: &PostgreSQLConfig{
				Host:     "localhost",
				Database: "mydb",
				User:     "postgres",
				Password: "secret",
				SSLMode:  "invalid",
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
