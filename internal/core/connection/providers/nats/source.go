// +marmot:name=NATS
// +marmot:description=NATS messaging system connection
// +marmot:status=stable
// +marmot:category=messaging
package nats

//go:generate go run ../../../../docgen/cmd/connection/main.go

import (
	"github.com/marmotdata/marmot/internal/core/connection"
)

// +marmot:config
type NatsConfig struct {
	Host            string `json:"host" label:"Host" description:"NATS server hostname or IP address" validate:"required" placeholder:"localhost"`
	Port            int    `json:"port,omitempty" label:"Port" description:"NATS server port" default:"4222" validate:"omitempty,min=1,max=65535"`
	Token           string `json:"token,omitempty" label:"Token" description:"Authentication token" sensitive:"true"`
	Username        string `json:"username,omitempty" label:"Username" description:"Username for authentication"`
	Password        string `json:"password,omitempty" label:"Password" description:"Password for authentication" sensitive:"true"`
	CredentialsFile string `json:"credentials_file,omitempty" label:"Credentials File" description:"Path to NATS credentials file (.creds)"`
	TLS             bool   `json:"tls,omitempty" label:"Enable TLS" description:"Enable TLS connection" default:"false"`
	TLSInsecure     bool   `json:"tls_insecure,omitempty" label:"TLS Insecure" description:"Skip TLS certificate verification" default:"false"`
}

// +marmot:example-config
var _ = `
host: localhost
port: 4222
username: your-username
password: your-password
tls: false
`

type Source struct{}

func (s *Source) Validate(rawConfig map[string]interface{}) error {
	config, err := connection.UnmarshalConfig[NatsConfig](rawConfig)
	if err != nil {
		return err
	}
	return connection.ValidateConfig(config)
}

func init() {
	connection.GetRegistry().Register(connection.ConnectionTypeMeta{
		ID:          "nats",
		Name:        "NATS",
		Description: "NATS messaging system connection",
		Icon:        "nats",
		Category:    "messaging",
		ConfigSpec:  connection.GenerateConfigSpec(NatsConfig{}),
	}, &Source{})
}
