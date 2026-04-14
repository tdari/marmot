package aws

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
			name: "valid config with region only",
			config: map[string]interface{}{
				"region": "us-east-1",
			},
			wantErr: false,
		},
		{
			name: "valid config with access key credentials",
			config: map[string]interface{}{
				"region": "us-west-2",
				"id":     "AKIAIOSFODNN7EXAMPLE",
				"secret": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			},
			wantErr: false,
		},
		{
			name: "valid config with session token",
			config: map[string]interface{}{
				"region": "us-east-1",
				"id":     "AKIAIOSFODNN7EXAMPLE",
				"secret": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				"token":  "session-token",
			},
			wantErr: false,
		},
		{
			name: "valid config with profile",
			config: map[string]interface{}{
				"region":      "us-east-1",
				"use_default": false,
				"profile":     "myprofile",
			},
			wantErr: false,
		},
		{
			name: "valid config with use_default",
			config: map[string]interface{}{
				"region":      "us-east-1",
				"use_default": true,
			},
			wantErr: false,
		},
		{
			name: "valid config with role ARN",
			config: map[string]interface{}{
				"region": "us-east-1",
				"id":     "AKIAIOSFODNN7EXAMPLE",
				"secret": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				"role":   "arn:aws:iam::123456789012:role/MyRole",
			},
			wantErr: false,
		},
		{
			name: "valid config with role and external ID",
			config: map[string]interface{}{
				"region":             "us-east-1",
				"id":                 "AKIAIOSFODNN7EXAMPLE",
				"secret":             "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				"role":               "arn:aws:iam::123456789012:role/MyRole",
				"role_external_id":   "external-id-123",
			},
			wantErr: false,
		},
		{
			name: "valid config with custom endpoint",
			config: map[string]interface{}{
				"region":   "us-east-1",
				"endpoint": "https://s3-compatible.example.com",
			},
			wantErr: false,
		},
		{
			name: "missing region",
			config: map[string]interface{}{
				"id":     "AKIAIOSFODNN7EXAMPLE",
				"secret": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
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

func TestAWSConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *AWSConfig
		valid  bool
	}{
		{
			name: "valid minimal config",
			config: &AWSConfig{
				Region: "us-east-1",
			},
			valid: true,
		},
		{
			name: "valid with credentials",
			config: &AWSConfig{
				Region: "us-east-1",
				ID:     "AKIAIOSFODNN7EXAMPLE",
				Secret: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			},
			valid: true,
		},
		{
			name: "missing region",
			config: &AWSConfig{
				ID:     "AKIAIOSFODNN7EXAMPLE",
				Secret: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			},
			valid: false,
		},
		{
			name: "valid with profile",
			config: &AWSConfig{
				Region:     "us-east-1",
				UseDefault: false,
				Profile:    "myprofile",
			},
			valid: true,
		},
		{
			name: "valid with use_default",
			config: &AWSConfig{
				Region:     "us-east-1",
				UseDefault: true,
			},
			valid: true,
		},
		{
			name: "valid with role",
			config: &AWSConfig{
				Region: "us-east-1",
				Role:   "arn:aws:iam::123456789012:role/MyRole",
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
