// +marmot:name=Google BigQuery
// +marmot:description=Google BigQuery data warehouse connection
// +marmot:status=stable
// +marmot:category=database
package bigquery

//go:generate go run ../../../../docgen/cmd/connection/main.go

import (
	"fmt"

	"github.com/marmotdata/marmot/internal/core/connection"
)

// +marmot:config
type BigQueryConfig struct {
	ProjectID             string `json:"project_id" label:"Project ID" description:"Google Cloud project ID" validate:"required" placeholder:"my-gcp-project"`
	Dataset               string `json:"dataset" label:"Dataset" description:"BigQuery dataset name" placeholder:"my_dataset"`
	CredentialsJSON       string `json:"credentials_json" label:"Service Account Key" description:"JSON key file content for service account" sensitive:"true"`
	CredentialsPath       string `json:"credentials_path" label:"Credentials Path" description:"Path to service account JSON key file" placeholder:"/path/to/key.json"`
	UseDefaultCredentials bool   `json:"use_default_credentials" label:"Use Default Credentials" description:"Use application default credentials" default:"false"`
	Location              string `json:"location" label:"Location" description:"BigQuery dataset location" default:"US" placeholder:"US"`
}

// +marmot:example-config
var _ = `
project_id: my-gcp-project
dataset: my_dataset
credentials_path: /path/to/key.json
location: US
`

type Source struct{}

func (s *Source) Validate(rawConfig map[string]interface{}) error {
	config, err := connection.UnmarshalConfig[BigQueryConfig](rawConfig)
	if err != nil {
		return err
	}

	authMethods := 0
	if config.CredentialsPath != "" {
		authMethods++
	}
	if config.CredentialsJSON != "" {
		authMethods++
	}
	if config.UseDefaultCredentials {
		authMethods++
	}

	if authMethods == 0 {
		return fmt.Errorf("at least one authentication method must be provided: credentials_path, credentials_json, or use_default_credentials")
	}
	if authMethods > 1 {
		return fmt.Errorf("only one authentication method should be provided")
	}

	return connection.ValidateConfig(config)
}

func init() {
	connection.GetRegistry().Register(connection.ConnectionTypeMeta{
		ID:          "bigquery",
		Name:        "Google BigQuery",
		Description: "Google BigQuery data warehouse connection",
		Icon:        "bigquery",
		Category:    "database",
		ConfigSpec:  connection.GenerateConfigSpec(BigQueryConfig{}),
	}, &Source{})
}
