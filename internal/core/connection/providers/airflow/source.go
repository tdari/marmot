// +marmot:name=Apache Airflow
// +marmot:description=Apache Airflow workflow orchestration connection
// +marmot:status=stable
// +marmot:category=orchestration
package airflow

//go:generate go run ../../../../docgen/cmd/connection/main.go

import (
	"fmt"
	"strings"

	"github.com/marmotdata/marmot/internal/core/connection"
)

// +marmot:config
type AirflowConfig struct {
	Host     string `json:"host" label:"Host" description:"Airflow webserver URL" validate:"required,url" placeholder:"http://localhost:8080"`
	Username string `json:"username,omitempty" label:"Username" description:"Username for basic authentication"`
	Password string `json:"password,omitempty" label:"Password" description:"Password for basic authentication" sensitive:"true"`
	APIToken string `json:"api_token,omitempty" label:"API Token" description:"API token for authentication (alternative to basic auth)" sensitive:"true"`
}

// +marmot:example-config
var _ = `
host: http://localhost:8080
username: admin
password: your-password
`

type Source struct{}

func (s *Source) Validate(rawConfig map[string]interface{}) error {
	config, err := connection.UnmarshalConfig[AirflowConfig](rawConfig)
	if err != nil {
		return err
	}

	config.Host = strings.TrimSuffix(config.Host, "/")

	if config.Username == "" && config.APIToken == "" {
		return fmt.Errorf("authentication required: provide either username/password or api_token")
	}

	return connection.ValidateConfig(config)
}

func init() {
	connection.GetRegistry().Register(connection.ConnectionTypeMeta{
		ID:          "airflow",
		Name:        "Apache Airflow",
		Description: "Apache Airflow workflow orchestration connection",
		Icon:        "airflow",
		Category:    "orchestration",
		ConfigSpec:  connection.GenerateConfigSpec(AirflowConfig{}),
	}, &Source{})
}
