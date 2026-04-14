// +marmot:name=Google Cloud Storage
// +marmot:description=GCS object storage connection
// +marmot:status=stable
// +marmot:category=storage
package gcs

//go:generate go run ../../../../docgen/cmd/connection/main.go

import (
	"github.com/marmotdata/marmot/internal/core/connection"
)

// +marmot:config
type GCSConfig struct {
	ProjectID       string `json:"project_id" label:"Project ID" description:"Google Cloud project ID" validate:"required" placeholder:"my-gcp-project"`
	BucketName      string `json:"bucket_name" label:"Bucket Name" description:"GCS bucket name" placeholder:"my-bucket"`
	CredentialsFile string `json:"credentials_file" label:"Credentials File" description:"Path to service account JSON key file" placeholder:"/path/to/key.json"`
	CredentialsJSON string `json:"credentials_json" label:"Credentials JSON" description:"Service account JSON key content" sensitive:"true"`
	Endpoint        string `json:"endpoint" label:"Custom Endpoint" description:"Custom GCS endpoint URL" placeholder:"https://storage.googleapis.com"`
	DisableAuth     bool   `json:"disable_auth" label:"Disable Authentication" description:"Disable authentication (for public buckets)" default:"false"`
}

// +marmot:example-config
var _ = `
project_id: my-gcp-project
bucket_name: my-bucket
credentials_file: /path/to/key.json
`

type Source struct{}

func (s *Source) Validate(rawConfig map[string]interface{}) error {
	config, err := connection.UnmarshalConfig[GCSConfig](rawConfig)
	if err != nil {
		return err
	}
	return connection.ValidateConfig(config)
}

func init() {
	connection.GetRegistry().Register(connection.ConnectionTypeMeta{
		ID:          "gcs",
		Name:        "Google Cloud Storage",
		Description: "GCS object storage connection",
		Icon:        "gcs",
		Category:    "storage",
		ConfigSpec:  connection.GenerateConfigSpec(GCSConfig{}),
	}, &Source{})
}
