// +marmot:name=MongoDB
// +marmot:description=MongoDB NoSQL database connection
// +marmot:status=stable
// +marmot:category=database
package mongodb

//go:generate go run ../../../../docgen/cmd/connection/main.go

import (
	"github.com/marmotdata/marmot/internal/core/connection"
)

// +marmot:config
type MongoDBConfig struct {
	ConnectionURI string `json:"connection_uri" label:"Connection URI" description:"MongoDB connection URI (mongodb:// or mongodb+srv://)" placeholder:"mongodb://localhost:27017"`
	Host          string `json:"host" label:"Host" description:"MongoDB server hostname or IP" placeholder:"localhost" show_when:"connection_uri:"`
	Port          int    `json:"port" label:"Port" description:"MongoDB server port" default:"27017" validate:"min=1,max=65535" show_when:"connection_uri:"`
	User          string `json:"user" label:"User" description:"MongoDB username" placeholder:"admin" show_when:"connection_uri:"`
	Password      string `json:"password" label:"Password" description:"MongoDB password" sensitive:"true" show_when:"connection_uri:"`
	AuthSource    string `json:"auth_source" label:"Auth Source" description:"Authentication database" default:"admin" placeholder:"admin"`
	TLS           bool   `json:"tls" label:"Enable TLS" description:"Enable TLS/SSL connection" default:"false"`
	TLSInsecure   bool   `json:"tls_insecure" label:"TLS Insecure" description:"Skip TLS certificate verification" default:"false"`
}

// +marmot:example-config
var _ = `
host: localhost
port: 27017
user: admin
password: your-password
auth_source: admin
tls: false
`

type Source struct{}

func (s *Source) Validate(rawConfig map[string]interface{}) error {
	config, err := connection.UnmarshalConfig[MongoDBConfig](rawConfig)
	if err != nil {
		return err
	}
	return connection.ValidateConfig(config)
}

func init() {
	connection.GetRegistry().Register(connection.ConnectionTypeMeta{
		ID:          "mongodb",
		Name:        "MongoDB",
		Description: "MongoDB NoSQL database connection",
		Icon:        "mongodb",
		Category:    "database",
		ConfigSpec:  connection.GenerateConfigSpec(MongoDBConfig{}),
	}, &Source{})
}
