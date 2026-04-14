package mysql

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
				"host":     "mysql.example.com",
				"port":     3306,
				"database": "mydb",
				"user":     "root",
				"password": "secret",
				"tls":      "true",
			},
			wantErr: false,
		},
		{
			name: "valid config with required fields",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "root",
				"password": "secret",
			},
			wantErr: false,
		},
		{
			name: "missing host",
			config: map[string]interface{}{
				"database": "mydb",
				"user":     "root",
				"password": "secret",
			},
			wantErr: true,
			errMsg:  "required",
		},
		{
			name: "missing database",
			config: map[string]interface{}{
				"host":     "localhost",
				"user":     "root",
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
				"user":     "root",
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
				"user":     "root",
				"password": "secret",
				"port":     65536,
			},
			wantErr: true,
			errMsg:  "must be at most",
		},
		{
			name: "invalid tls mode",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "root",
				"password": "secret",
				"tls":      "invalid",
			},
			wantErr: true,
			errMsg:  "must be one of",
		},
		{
			name: "valid tls mode - false",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "root",
				"password": "secret",
				"tls":      "false",
			},
			wantErr: false,
		},
		{
			name: "valid tls mode - true",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "root",
				"password": "secret",
				"tls":      "true",
			},
			wantErr: false,
		},
		{
			name: "valid tls mode - skip-verify",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "root",
				"password": "secret",
				"tls":      "skip-verify",
			},
			wantErr: false,
		},
		{
			name: "valid tls mode - preferred",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "root",
				"password": "secret",
				"tls":      "preferred",
			},
			wantErr: false,
		},
		{
			name: "port defaults to 3306",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "root",
				"password": "secret",
			},
			wantErr: false,
		},
		{
			name: "tls defaults to false",
			config: map[string]interface{}{
				"host":     "localhost",
				"database": "mydb",
				"user":     "root",
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

func TestMySQLConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *MySQLConfig
		valid  bool
	}{
		{
			name: "valid config",
			config: &MySQLConfig{
				Host:     "localhost",
				Port:     3306,
				Database: "mydb",
				User:     "root",
				Password: "secret",
				TLS:      "false",
			},
			valid: true,
		},
		{
			name: "missing host",
			config: &MySQLConfig{
				Database: "mydb",
				User:     "root",
				Password: "secret",
			},
			valid: false,
		},
		{
			name: "missing database",
			config: &MySQLConfig{
				Host:     "localhost",
				User:     "root",
				Password: "secret",
			},
			valid: false,
		},
		{
			name: "missing user",
			config: &MySQLConfig{
				Host:     "localhost",
				Database: "mydb",
				Password: "secret",
			},
			valid: false,
		},
		{
			name: "missing password",
			config: &MySQLConfig{
				Host:     "localhost",
				Database: "mydb",
				User:     "root",
			},
			valid: false,
		},
		{
			name: "invalid port",
			config: &MySQLConfig{
				Host:     "localhost",
				Port:     99999,
				Database: "mydb",
				User:     "root",
				Password: "secret",
			},
			valid: false,
		},
		{
			name: "invalid tls mode",
			config: &MySQLConfig{
				Host:     "localhost",
				Database: "mydb",
				User:     "root",
				Password: "secret",
				TLS:      "invalid",
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
