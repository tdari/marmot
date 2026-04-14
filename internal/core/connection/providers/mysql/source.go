// +marmot:name=MySQL
// +marmot:description=MySQL database connection
// +marmot:status=stable
// +marmot:category=database
package mysql

//go:generate go run ../../../../docgen/cmd/connection/main.go

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/marmotdata/marmot/internal/core/connection"
)

// +marmot:config
type MySQLConfig struct {
	Host     string `json:"host" label:"Host" description:"MySQL server hostname or IP address" validate:"required" placeholder:"localhost"`
	Port     int    `json:"port" label:"Port" description:"MySQL server port" default:"3306" validate:"min=1,max=65535"`
	Database string `json:"database" label:"Database" description:"Database name to connect to" validate:"required" placeholder:"mydb"`
	User     string `json:"user" label:"User" description:"Database username for authentication" validate:"required"`
	Password string `json:"password" label:"Password" description:"Database password for authentication" sensitive:"true" validate:"required"`
	TLS      string `json:"tls" label:"TLS Mode" description:"TLS/SSL mode" default:"false" validate:"oneof=false true skip-verify preferred" placeholder:"false"`
}

// +marmot:example-config
var _ = `
host: localhost
port: 3306
database: mydb
user: admin
password: your-password
tls: "false"
`

type Source struct{}

func (s *Source) Validate(rawConfig map[string]interface{}) error {

	config, err := connection.UnmarshalConfig[MySQLConfig](rawConfig)
	if err != nil {
		return err
	}
	// Apply defaults for optional fields
	if config.Port == 0 {
		config.Port = 3306
	}
	if config.TLS == "" {
		config.TLS = "false"
	}
	return connection.ValidateConfig(config)
}

func init() {
	connection.GetRegistry().Register(connection.ConnectionTypeMeta{
		ID:          "mysql",
		Name:        "MySQL",
		Description: "MySQL database connection",
		Icon:        "mysql",
		Category:    "database",
		ConfigSpec:  connection.GenerateConfigSpec(MySQLConfig{}),
	}, &Source{})
}
