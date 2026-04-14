package azureblob

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
			name: "valid config with connection string",
			config: map[string]interface{}{
				"connection_string": "DefaultEndpointsProtocol=https;AccountName=myaccount;AccountKey=mykey==;EndpointSuffix=core.windows.net",
				"container_name":    "mycontainer",
			},
			wantErr: false,
		},
		{
			name: "valid config with account name and key",
			config: map[string]interface{}{
				"account_name":   "mystorageaccount",
				"account_key":    "myaccountkey==",
				"container_name": "mycontainer",
			},
			wantErr: false,
		},
		{
			name: "missing both connection string and account name",
			config: map[string]interface{}{
				"container_name": "mycontainer",
			},
			wantErr: true,
			errMsg:  "connection_string or account_name",
		},
		{
			name: "account name without account key",
			config: map[string]interface{}{
				"account_name":   "mystorageaccount",
				"container_name": "mycontainer",
			},
			wantErr: true,
			errMsg:  "account_key is required",
		},
		{
			name: "empty config",
			config: map[string]interface{}{},
			wantErr: true,
			errMsg:  "connection_string or account_name",
		},
		{
			name: "valid with custom endpoint",
			config: map[string]interface{}{
				"connection_string": "DefaultEndpointsProtocol=https;AccountName=myaccount;AccountKey=mykey==;EndpointSuffix=core.windows.net",
				"endpoint":          "https://mystorageaccount.blob.core.windows.net",
			},
			wantErr: false,
		},
		{
			name: "valid with account name, key and endpoint",
			config: map[string]interface{}{
				"account_name":   "mystorageaccount",
				"account_key":    "myaccountkey==",
				"endpoint":       "https://mystorageaccount.blob.core.windows.net",
				"container_name": "mycontainer",
			},
			wantErr: false,
		},
		{
			name: "connection string takes precedence over account name",
			config: map[string]interface{}{
				"connection_string": "DefaultEndpointsProtocol=https;AccountName=myaccount;AccountKey=mykey==;EndpointSuffix=core.windows.net",
				"account_name":      "ignored",
				"account_key":       "ignored",
				"container_name":    "mycontainer",
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

func TestAzureBlobConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config *AzureBlobConfig
		valid  bool
	}{
		{
			name: "valid with connection string",
			config: &AzureBlobConfig{
				ConnectionString: "DefaultEndpointsProtocol=https;AccountName=myaccount;AccountKey=mykey==;EndpointSuffix=core.windows.net",
				ContainerName:    "mycontainer",
			},
			valid: true,
		},
		{
			name: "valid with account name and key",
			config: &AzureBlobConfig{
				AccountName:   "mystorageaccount",
				AccountKey:    "myaccountkey==",
				ContainerName: "mycontainer",
			},
			valid: true,
		},
		{
			name: "valid with custom endpoint",
			config: &AzureBlobConfig{
				ConnectionString: "DefaultEndpointsProtocol=https;AccountName=myaccount;AccountKey=mykey==;EndpointSuffix=core.windows.net",
				Endpoint:         "https://mystorageaccount.blob.core.windows.net",
				ContainerName:    "mycontainer",
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
