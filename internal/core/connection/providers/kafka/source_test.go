package kafka

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
			name: "valid config with bootstrap servers only",
			config: map[string]interface{}{
				"bootstrap_servers": "kafka-1:9092,kafka-2:9092",
			},
			wantErr: false,
		},
		{
			name: "valid config with client ID",
			config: map[string]interface{}{
				"bootstrap_servers": "kafka-1:9092",
				"client_id":         "marmot-consumer",
			},
			wantErr: false,
		},
		{
			name: "missing bootstrap servers",
			config: map[string]interface{}{
				"client_id": "marmot-consumer",
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
			name: "valid SASL plaintext auth",
			config: map[string]interface{}{
				"bootstrap_servers": "kafka-1:9092",
				"authentication": map[string]interface{}{
					"type":      "sasl_plaintext",
					"username":  "user",
					"password":  "pass",
					"mechanism": "PLAIN",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid SASL auth - missing username",
			config: map[string]interface{}{
				"bootstrap_servers": "kafka-1:9092",
				"authentication": map[string]interface{}{
					"type":      "sasl_plaintext",
					"password":  "pass",
					"mechanism": "PLAIN",
				},
			},
			wantErr: true,
			errMsg:  "username is required",
		},
		{
			name: "invalid SASL auth - missing password",
			config: map[string]interface{}{
				"bootstrap_servers": "kafka-1:9092",
				"authentication": map[string]interface{}{
					"type":      "sasl_plaintext",
					"username":  "user",
					"mechanism": "PLAIN",
				},
			},
			wantErr: true,
			errMsg:  "password is required",
		},
		{
			name: "invalid SASL auth - missing mechanism",
			config: map[string]interface{}{
				"bootstrap_servers": "kafka-1:9092",
				"authentication": map[string]interface{}{
					"type":     "sasl_plaintext",
					"username": "user",
					"password": "pass",
				},
			},
			wantErr: true,
			errMsg:  "mechanism is required",
		},
		{
			name: "valid SASL SSL auth",
			config: map[string]interface{}{
				"bootstrap_servers": "kafka-1:9092",
				"authentication": map[string]interface{}{
					"type":      "sasl_ssl",
					"username":  "user",
					"password":  "pass",
					"mechanism": "SCRAM-SHA-256",
				},
			},
			wantErr: false,
		},
		{
			name: "valid none auth",
			config: map[string]interface{}{
				"bootstrap_servers": "kafka-1:9092",
				"authentication": map[string]interface{}{
					"type": "none",
				},
			},
			wantErr: false,
		},
		{
			name: "valid SSL auth",
			config: map[string]interface{}{
				"bootstrap_servers": "kafka-1:9092",
				"authentication": map[string]interface{}{
					"type": "ssl",
				},
			},
			wantErr: false,
		},
		{
			name: "valid with TLS config",
			config: map[string]interface{}{
				"bootstrap_servers": "kafka-1:9092",
				"tls": map[string]interface{}{
					"enabled": true,
					"cert_path": "/path/to/cert.pem",
					"key_path": "/path/to/key.pem",
				},
			},
			wantErr: false,
		},
		{
			name: "valid with schema registry",
			config: map[string]interface{}{
				"bootstrap_servers": "kafka-1:9092",
				"schema_registry": map[string]interface{}{
					"enabled": true,
					"url": "https://schema-registry:8081",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid auth type",
			config: map[string]interface{}{
				"bootstrap_servers": "kafka-1:9092",
				"authentication": map[string]interface{}{
					"type": "invalid_type",
				},
			},
			wantErr: true,
			errMsg:  "must be one of",
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

func TestKafkaConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *KafkaConfig
		valid  bool
	}{
		{
			name: "valid minimal config",
			config: &KafkaConfig{
				BootstrapServers: "kafka-1:9092",
			},
			valid: true,
		},
		{
			name: "missing bootstrap servers",
			config: &KafkaConfig{},
			valid: false,
		},
		{
			name: "valid with auth",
			config: &KafkaConfig{
				BootstrapServers: "kafka-1:9092",
				Authentication: &KafkaAuthConfig{
					Type:      "sasl_plaintext",
					Username:  "user",
					Password:  "pass",
					Mechanism: "PLAIN",
				},
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
