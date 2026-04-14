// +marmot:name=Trino
// +marmot:description=Trino distributed SQL query engine connection
// +marmot:status=stable
// +marmot:category=database
package trino

//go:generate go run ../../../../docgen/cmd/connection/main.go

import (
	"github.com/marmotdata/marmot/internal/core/connection"
)

// +marmot:config
type TrinoConfig struct {
	Host        string `json:"host" label:"Host" description:"Trino coordinator hostname" validate:"required" placeholder:"trino.company.com"`
	Port        int    `json:"port" label:"Port" description:"Trino coordinator port" validate:"omitempty,min=1,max=65535" default:"8080" placeholder:"8080"`
	User        string `json:"user" label:"User" description:"Username for authentication" validate:"required" placeholder:"marmot_reader"`
	Password    string `json:"password,omitempty" label:"Password" description:"Password (requires HTTPS)" sensitive:"true"`
	Secure      bool   `json:"secure,omitempty" label:"Secure" description:"Use HTTPS" default:"false"`
	SSLCertPath string `json:"ssl_cert_path,omitempty" label:"SSL Cert Path" description:"Path to TLS certificate file" placeholder:"/path/to/cert.pem"`
	AccessToken string `json:"access_token,omitempty" label:"Access Token" description:"JWT bearer token" sensitive:"true"`
}

// +marmot:example-config
var _ = `
host: trino.company.com
port: 8080
user: marmot_reader
password: your-password
secure: false
`

type Source struct{}

func (s *Source) Validate(rawConfig map[string]interface{}) error {
	config, err := connection.UnmarshalConfig[TrinoConfig](rawConfig)
	if err != nil {
		return err
	}

	if config.Port == 0 {
		config.Port = 8080
	}

	return connection.ValidateConfig(config)
}

func init() {
	connection.GetRegistry().Register(connection.ConnectionTypeMeta{
		ID:          "trino",
		Name:        "Trino",
		Description: "Trino distributed SQL query engine connection",
		Icon:        "trino",
		Category:    "database",
		ConfigSpec:  connection.GenerateConfigSpec(TrinoConfig{}),
	}, &Source{})
}
