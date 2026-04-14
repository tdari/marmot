// +marmot:name=Azure Blob Storage
// +marmot:description=Azure Blob Storage connection
// +marmot:status=stable
// +marmot:category=storage
package azureblob

//go:generate go run ../../../../docgen/cmd/connection/main.go

import (
	"fmt"

	"github.com/marmotdata/marmot/internal/core/connection"
)

// +marmot:config
type AzureBlobConfig struct {
	ConnectionString string `json:"connection_string" label:"Connection String" description:"Azure Storage connection string" sensitive:"true" placeholder:"DefaultEndpointsProtocol=https;AccountName=..."`
	AccountName      string `json:"account_name" label:"Account Name" description:"Storage account name" placeholder:"mystorageaccount" show_when:"connection_string:"`
	AccountKey       string `json:"account_key" label:"Account Key" description:"Storage account key" sensitive:"true" show_when:"connection_string:"`
	ContainerName    string `json:"container_name" label:"Container Name" description:"Blob container name" placeholder:"mycontainer"`
	Endpoint         string `json:"endpoint" label:"Custom Endpoint" description:"Custom blob endpoint URL" placeholder:"https://mystorageaccount.blob.core.windows.net"`
}

// +marmot:example-config
var _ = `
account_name: mystorageaccount
account_key: your-api-secret
container_name: mycontainer
`

type Source struct{}

func (s *Source) Validate(rawConfig map[string]interface{}) error {
	config, err := connection.UnmarshalConfig[AzureBlobConfig](rawConfig)
	if err != nil {
		return err
	}

	if config.ConnectionString == "" && config.AccountName == "" {
		return fmt.Errorf("either connection_string or account_name must be provided")
	}

	if config.AccountName != "" && config.AccountKey == "" && config.ConnectionString == "" {
		return fmt.Errorf("account_key is required when using account_name")
	}

	return connection.ValidateConfig(config)
}

func init() {
	connection.GetRegistry().Register(connection.ConnectionTypeMeta{
		ID:          "azureblob",
		Name:        "Azure Blob Storage",
		Description: "Azure Blob Storage connection",
		Icon:        "azure-blob",
		Category:    "storage",
		ConfigSpec:  connection.GenerateConfigSpec(AzureBlobConfig{}),
	}, &Source{})
}
