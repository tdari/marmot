// +marmot:name=ClickHouse
// +marmot:description=ClickHouse columnar database connection
// +marmot:status=stable
// +marmot:category=database
package clickhouse

//go:generate go run ../../../../docgen/cmd/connection/main.go

import (
	"github.com/marmotdata/marmot/internal/core/connection"
)

// +marmot:config
type ClickHouseConfig struct {
	Host     string `json:"host" label:"Host" description:"ClickHouse server hostname or IP address" validate:"required" placeholder:"localhost"`
	Port     int    `json:"port" label:"Port" description:"ClickHouse server port" default:"9000" validate:"min=1,max=65535"`
	User     string `json:"user" label:"User" description:"ClickHouse username" validate:"required" placeholder:"default"`
	Password string `json:"password" label:"Password" description:"ClickHouse password" sensitive:"true"`
	Database string `json:"database" label:"Database" description:"Database name" validate:"required" placeholder:"default"`
	Secure   bool   `json:"secure" label:"Secure Connection" description:"Use secure connection (TLS)" default:"false"`
}

// +marmot:example-config
var _ = `
host: localhost
port: 9000
user: default
password: your-password
database: default
secure: false
`

type Source struct{}

func (s *Source) Validate(rawConfig map[string]interface{}) error {
	config, err := connection.UnmarshalConfig[ClickHouseConfig](rawConfig)
	if err != nil {
		return err
	}

	if config.Port == 0 {
		config.Port = 9000
	}

	if config.Database == "" {
		config.Database = "default"
	}

	return connection.ValidateConfig(config)
}

func init() {
	connection.GetRegistry().Register(connection.ConnectionTypeMeta{
		ID:          "clickhouse",
		Name:        "ClickHouse",
		Description: "ClickHouse columnar database connection",
		Icon:        "clickhouse",
		Category:    "database",
		ConfigSpec:  connection.GenerateConfigSpec(ClickHouseConfig{}),
	}, &Source{})
}
