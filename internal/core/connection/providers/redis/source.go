// +marmot:name=Redis
// +marmot:description=Redis in-memory data store connection
// +marmot:status=stable
// +marmot:category=database
package redis

//go:generate go run ../../../../docgen/cmd/connection/main.go

import (
	"github.com/marmotdata/marmot/internal/core/connection"
)

// +marmot:config
type RedisConfig struct {
	Host        string `json:"host" label:"Host" description:"Redis server hostname or IP address" validate:"required" placeholder:"localhost"`
	Port        int    `json:"port,omitempty" label:"Port" description:"Redis server port" default:"6379" validate:"omitempty,min=1,max=65535"`
	Password    string `json:"password,omitempty" label:"Password" description:"Password for authentication" sensitive:"true"`
	Username    string `json:"username,omitempty" label:"Username" description:"Username for ACL authentication"`
	DB          int    `json:"db,omitempty" label:"Database" description:"Default database number" default:"0" validate:"omitempty,min=0,max=15"`
	TLS         bool   `json:"tls,omitempty" label:"Enable TLS" description:"Enable TLS connection" default:"false"`
	TLSInsecure bool   `json:"tls_insecure,omitempty" label:"TLS Insecure" description:"Skip TLS certificate verification" default:"false"`
}

// +marmot:example-config
var _ = `
host: localhost
port: 6379
password: your-password
db: 0
tls: false
`

type Source struct{}

func (s *Source) Validate(rawConfig map[string]interface{}) error {
	config, err := connection.UnmarshalConfig[RedisConfig](rawConfig)
	if err != nil {
		return err
	}
	if config.Port == 0 {
		config.Port = 6379
	}
	if config.DB < 0 {
		config.DB = 0
	}
	return connection.ValidateConfig(config)
}

func init() {
	connection.GetRegistry().Register(connection.ConnectionTypeMeta{
		ID:          "redis",
		Name:        "Redis",
		Description: "Redis in-memory data store connection",
		Icon:        "redis",
		Category:    "database",
		ConfigSpec:  connection.GenerateConfigSpec(RedisConfig{}),
	}, &Source{})
}
