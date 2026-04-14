// +marmot:name=Apache Iceberg (REST Catalog)
// +marmot:description=Apache Iceberg REST catalog connection
// +marmot:status=stable
// +marmot:category=data-lake
package iceberg

//go:generate go run ../../../../docgen/cmd/connection/main.go

import (
	"fmt"
	"net/url"

	"github.com/marmotdata/marmot/internal/core/connection"
)

// +marmot:config
type IcebergRESTConfig struct {
	URI        string            `json:"uri" label:"URI" description:"REST catalog URI" validate:"required,url" placeholder:"http://localhost:8181"`
	Warehouse  string            `json:"warehouse" label:"Warehouse" description:"Warehouse identifier" validate:"required"`
	Credential string            `json:"credential,omitempty" label:"Credential" description:"Credential for OAuth2 client credentials authentication" sensitive:"true"`
	Token      string            `json:"token,omitempty" label:"Token" description:"Bearer token for authentication" sensitive:"true"`
	Properties map[string]string `json:"properties,omitempty" label:"Properties" description:"Additional catalog properties"`
	Prefix     string            `json:"prefix,omitempty" label:"Prefix" description:"Optional prefix for the REST catalog"`
}

// +marmot:example-config
var _ = `
uri: http://localhost:8181
warehouse: my-warehouse
token: your-api-secret
`

type Source struct{}

func (s *Source) Validate(rawConfig map[string]interface{}) error {
	config, err := connection.UnmarshalConfig[IcebergRESTConfig](rawConfig)
	if err != nil {
		return err
	}

	u, err := url.ParseRequestURI(config.URI)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("uri must be a valid URL")
	}

	if config.Credential == "" && config.Token == "" {
		return fmt.Errorf("authentication required: provide either credential or token")
	}

	return connection.ValidateConfig(config)
}

func init() {
	connection.GetRegistry().Register(connection.ConnectionTypeMeta{
		ID:          "iceberg-rest",
		Name:        "Apache Iceberg (REST Catalog)",
		Description: "Apache Iceberg REST catalog connection",
		Icon:        "iceberg",
		Category:    "data-lake",
		ConfigSpec:  connection.GenerateConfigSpec(IcebergRESTConfig{}),
	}, &Source{})
}
