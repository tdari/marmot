package nats

import (
	"testing"

	"github.com/marmotdata/marmot/internal/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSource_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  plugin.RawPluginConfig
		wantErr string
	}{
		{
			name: "valid config with plugin fields",
			config: plugin.RawPluginConfig{
				"filter": map[string]interface{}{
					"include": []interface{}{"^ORDERS"},
				},
			},
		},
		{
			name: "empty config",
			config: plugin.RawPluginConfig{},
		},
		{
			name: "with host and port",
			config: plugin.RawPluginConfig{
				"host": "localhost",
				"port": 4222,
			},
		},
		{
			name: "with token",
			config: plugin.RawPluginConfig{
				"host":  "localhost",
				"token": "s3cr3t",
			},
		},
		{
			name: "with credentials file",
			config: plugin.RawPluginConfig{
				"host":             "localhost",
				"credentials_file": "/path/to/creds.creds",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Source{}
			_, err := s.Validate(tt.config)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSource_ValidateDefaults(t *testing.T) {
	// Validate no longer sets connConfig - it's set by Discover
	// Just verify that Validate succeeds
	s := &Source{}
	_, err := s.Validate(plugin.RawPluginConfig{
		"host": "localhost",
	})
	require.NoError(t, err)
}
