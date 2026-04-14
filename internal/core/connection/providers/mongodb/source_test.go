package mongodb

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
			name: "valid config with connection URI and port",
			config: map[string]interface{}{
				"connection_uri": "mongodb://localhost:27017",
				"port": 27017,
			},
			wantErr: false,
		},
		{
			name: "valid config with mongodb+srv URI and port",
			config: map[string]interface{}{
				"connection_uri": "mongodb+srv://user:pass@cluster.mongodb.net/dbname",
				"port": 27017,
			},
			wantErr: false,
		},
		{
			name: "valid config with host/port",
			config: map[string]interface{}{
				"host": "localhost",
				"port": 27017,
				"user": "admin",
				"password": "secret",
			},
			wantErr: false,
		},
		{
			name: "valid config with auth_source",
			config: map[string]interface{}{
				"connection_uri": "mongodb://localhost:27017",
				"auth_source": "admin",
				"port": 27017,
			},
			wantErr: false,
		},
		{
			name: "valid config with TLS",
			config: map[string]interface{}{
				"host": "mongodb.example.com",
				"port": 27017,
				"user": "admin",
				"password": "secret",
				"tls": true,
				"tls_insecure": false,
			},
			wantErr: false,
		},
		{
			name: "valid with minimal host/port",
			config: map[string]interface{}{
				"host": "localhost",
				"port": 27017,
			},
			wantErr: false,
		},
		{
			name: "invalid port - too high",
			config: map[string]interface{}{
				"host": "localhost",
				"port": 99999,
			},
			wantErr: true,
			errMsg:  "must be at most",
		},
		{
			name: "valid with URI and TLS insecure",
			config: map[string]interface{}{
				"connection_uri": "mongodb+srv://localhost:27017",
				"port": 27017,
				"tls": true,
				"tls_insecure": true,
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

func TestMongoDBConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *MongoDBConfig
		valid  bool
	}{
		{
			name: "valid with connection URI",
			config: &MongoDBConfig{
				ConnectionURI: "mongodb://localhost:27017",
				Port: 27017,
			},
			valid: true,
		},
		{
			name: "valid with host port user password",
			config: &MongoDBConfig{
				Host:     "localhost",
				Port:     27017,
				User:     "admin",
				Password: "secret",
			},
			valid: true,
		},
		{
			name: "valid with host and port",
			config: &MongoDBConfig{
				Host: "localhost",
				Port: 27017,
			},
			valid: true,
		},
		{
			name: "valid with connection URI",
			config: &MongoDBConfig{
				ConnectionURI: "mongodb://localhost:27017",
				Port: 27017,
			},
			valid: true,
		},
		{
			name: "invalid port",
			config: &MongoDBConfig{
				Host: "localhost",
				Port: 99999,
			},
			valid: false,
		},
		{
			name: "valid with auth source",
			config: &MongoDBConfig{
				ConnectionURI: "mongodb://localhost:27017",
				Port: 27017,
				AuthSource: "admin",
			},
			valid: true,
		},
		{
			name: "valid with TLS",
			config: &MongoDBConfig{
				Host: "localhost",
				Port: 27017,
				TLS:  true,
				TLSInsecure: false,
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
