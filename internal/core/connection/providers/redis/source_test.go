package redis

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
				"host":     "redis.example.com",
				"port":     6379,
				"password": "secret",
				"username": "default",
				"db":       0,
				"tls":      true,
				"tls_insecure": true,
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
				"port": 6379,
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
			name: "invalid db - too high",
			config: map[string]interface{}{
				"host": "localhost",
				"db":   16,
			},
			wantErr: true,
			errMsg:  "must be at most",
		},
		{
			name: "valid config with TLS",
			config: map[string]interface{}{
				"host":         "redis.example.com",
				"tls":          true,
				"tls_insecure": false,
			},
			wantErr: false,
		},
		{
			name: "port defaults to 6379",
			config: map[string]interface{}{
				"host": "localhost",
			},
			wantErr: false,
		},
		{
			name: "db defaults to 0",
			config: map[string]interface{}{
				"host": "localhost",
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

func TestRedisConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *RedisConfig
		valid  bool
	}{
		{
			name: "valid config",
			config: &RedisConfig{
				Host: "localhost",
				Port: 6379,
				DB:   0,
			},
			valid: true,
		},
		{
			name: "missing host",
			config: &RedisConfig{
				Port: 6379,
			},
			valid: false,
		},
		{
			name: "invalid port",
			config: &RedisConfig{
				Host: "localhost",
				Port: 99999,
			},
			valid: false,
		},
		{
			name: "negative db",
			config: &RedisConfig{
				Host: "localhost",
				DB:   -1,
			},
			valid: false,
		},
		{
			name: "db out of range",
			config: &RedisConfig{
				Host: "localhost",
				DB:   16,
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
