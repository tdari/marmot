// +marmot:name=PostgreSQL
// +marmot:description=PostgreSQL database connection
// +marmot:status=stable
// +marmot:category=database
package postgresql

//go:generate go run ../../../../docgen/cmd/connection/main.go

import (
	"github.com/marmotdata/marmot/internal/core/connection"
)

// +marmot:config
type PostgreSQLConfig struct {
	Host     string `json:"host" label:"Host" description:"PostgreSQL server hostname or IP address" validate:"required" placeholder:"localhost"`
	Port     int    `json:"port" label:"Port" description:"PostgreSQL server port" default:"5432" validate:"min=1,max=65535"`
	Database string `json:"database" label:"Database" description:"Database name to connect to" validate:"required" placeholder:"mydb"`
	User     string `json:"user" label:"User" description:"Database username for authentication" validate:"required"`
	Password string `json:"password" label:"Password" description:"Database password for authentication" sensitive:"true" validate:"required"`
	SSLMode  string `json:"ssl_mode" label:"SSL Mode" description:"SSL connection mode" default:"disable" validate:"oneof=disable allow prefer require verify-ca verify-full"`
}

// +marmot:example-config
var _ = `
host: localhost
port: 5432
database: mydb
user: admin
password: your-password
ssl_mode: disable
`

type Source struct{}

func (s *Source) Validate(rawConfig map[string]interface{}) error {
	config, err := connection.UnmarshalConfig[PostgreSQLConfig](rawConfig)
	if err != nil {
		return err
	}
	if config.Port == 0 {
		config.Port = 5432
	}
	if config.SSLMode == "" {
		config.SSLMode = "disable"
	}
	return connection.ValidateConfig(config)
}

func init() {
	connection.GetRegistry().Register(connection.ConnectionTypeMeta{
		ID:          "postgresql",
		Name:        "PostgreSQL",
		Description: "PostgreSQL database connection",
		Icon:        "postgresql",
		Category:    "database",
		ConfigSpec:  connection.GenerateConfigSpec(PostgreSQLConfig{}),
	}, &Source{})
}
