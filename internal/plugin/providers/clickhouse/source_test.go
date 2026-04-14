package clickhouse

import (
	"testing"

	"github.com/marmotdata/marmot/internal/core/connection/providers/clickhouse"
	"github.com/marmotdata/marmot/internal/plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSource_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      plugin.RawPluginConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "valid config",
			config: plugin.RawPluginConfig{
				"host":     "localhost",
				"port":     9000,
				"user":     "default",
				"password": "password",
				"database": "default",
			},
			wantErr: false,
		},
		{
			name: "empty config",
			config: plugin.RawPluginConfig{},
			wantErr: false,
		},
		{
			name: "config with connection fields",
			config: plugin.RawPluginConfig{
				"host": "localhost",
				"user": "default",
			},
			wantErr: false,
		},
		{
			name: "config with secure connection",
			config: plugin.RawPluginConfig{
				"host":   "clickhouse.example.com",
				"user":   "admin",
				"secure": true,
			},
			wantErr: false,
		},
		{
			name: "config with filters",
			config: plugin.RawPluginConfig{
				"filter": map[string]interface{}{
					"include": []interface{}{"^analytics.*"},
					"exclude": []interface{}{".*_temp$"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Source{}
			_, err := s.Validate(tt.config)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSource_ValidateDefaults(t *testing.T) {
	// Validate no longer sets connConfig - it's set by Discover
	s := &Source{}
	_, err := s.Validate(plugin.RawPluginConfig{})

	require.NoError(t, err)
	require.NotNil(t, s.config)
}

func TestConfig_Defaults(t *testing.T) {
	config := &clickhouse.ClickHouseConfig{}

	assert.Equal(t, "", config.Password)
	assert.Equal(t, 0, config.Port)
	assert.Equal(t, "", config.Database)
	assert.False(t, config.Secure)
}
