// +marmot:name=Apache Kafka
// +marmot:description=Kafka streaming platform connection
// +marmot:status=stable
// +marmot:category=messaging
package kafka

//go:generate go run ../../../../docgen/cmd/connection/main.go

import (
	"fmt"

	"github.com/marmotdata/marmot/internal/core/connection"
)

type KafkaAuthConfig struct {
	Type      string `json:"type" label:"Auth Type" description:"Authentication type" validate:"omitempty,oneof=none sasl_plaintext sasl_ssl ssl" default:"none"`
	Username  string `json:"username,omitempty" label:"Username" description:"SASL username" placeholder:"your-username"`
	Password  string `json:"password,omitempty" label:"Password" description:"SASL password" sensitive:"true"`
	Mechanism string `json:"mechanism,omitempty" label:"SASL Mechanism" description:"SASL mechanism" validate:"omitempty,oneof=PLAIN SCRAM-SHA-256 SCRAM-SHA-512" default:"PLAIN"`
}

type KafkaTLSConfig struct {
	Enabled    bool   `json:"enabled" label:"Enable TLS" description:"Whether to enable TLS" default:"false"`
	CertPath   string `json:"cert_path,omitempty" label:"Certificate Path" description:"Path to TLS certificate file" placeholder:"/path/to/cert.pem"`
	KeyPath    string `json:"key_path,omitempty" label:"Key Path" description:"Path to TLS key file" placeholder:"/path/to/key.pem"`
	CACertPath string `json:"ca_cert_path,omitempty" label:"CA Certificate Path" description:"Path to TLS CA certificate file" placeholder:"/path/to/ca.pem"`
	SkipVerify bool   `json:"skip_verify,omitempty" label:"Skip Verification" description:"Skip TLS certificate verification" default:"false"`
}

type KafkaSchemaRegistryConfig struct {
	URL        string            `json:"url" label:"URL" description:"Schema Registry URL" validate:"omitempty,url" placeholder:"https://schema-registry:8081"`
	Config     map[string]string `json:"config,omitempty" label:"Config" description:"Additional Schema Registry configuration"`
	Enabled    bool              `json:"enabled" label:"Enabled" description:"Whether to use Schema Registry" default:"false"`
	SkipVerify bool              `json:"skip_verify,omitempty" label:"Skip Verification" description:"Skip TLS certificate verification" default:"false"`
}

// +marmot:config
type KafkaConfig struct {
	BootstrapServers string                     `json:"bootstrap_servers" label:"Bootstrap Servers" description:"Comma-separated list of bootstrap servers" validate:"required" placeholder:"kafka-1:9092,kafka-2:9092"`
	ClientID         string                     `json:"client_id" label:"Client ID" description:"Client ID for the consumer" placeholder:"marmot-discovery"`
	Authentication   *KafkaAuthConfig           `json:"authentication,omitempty" label:"Authentication" description:"Authentication configuration"`
	ConsumerConfig   map[string]string          `json:"consumer_config,omitempty" label:"Consumer Config" description:"Additional consumer configuration"`
	ClientTimeout    int                        `json:"client_timeout_seconds" label:"Client Timeout" description:"Request timeout in seconds" validate:"omitempty,min=1,max=300" default:"30"`
	TLS              *KafkaTLSConfig            `json:"tls,omitempty" label:"TLS" description:"TLS configuration"`
	SchemaRegistry   *KafkaSchemaRegistryConfig `json:"schema_registry,omitempty" label:"Schema Registry" description:"Schema Registry configuration"`
}

// +marmot:example-config
var _ = `
bootstrap_servers: kafka-1:9092,kafka-2:9092
client_id: marmot-discovery
client_timeout_seconds: 30
authentication:
  type: sasl_ssl
  username: your-username
  password: your-password
  mechanism: PLAIN
tls:
  enabled: true
`

type Source struct{}

func (s *Source) Validate(rawConfig map[string]interface{}) error {
	config, err := connection.UnmarshalConfig[KafkaConfig](rawConfig)
	if err != nil {
		return err
	}

	if config.TLS == nil {
		config.TLS = &KafkaTLSConfig{
			Enabled: true,
		}
	}

	if config.Authentication != nil {
		authType := config.Authentication.Type
		if authType == "sasl_plaintext" || authType == "sasl_ssl" {
			if config.Authentication.Username == "" {
				return fmt.Errorf("username is required for %s authentication", authType)
			}
			if config.Authentication.Password == "" {
				return fmt.Errorf("password is required for %s authentication", authType)
			}
			if config.Authentication.Mechanism == "" {
				return fmt.Errorf("mechanism is required for %s authentication", authType)
			}
		}
	}

	return connection.ValidateConfig(config)
}

func init() {
	connection.GetRegistry().Register(connection.ConnectionTypeMeta{
		ID:          "kafka",
		Name:        "Apache Kafka",
		Description: "Kafka streaming platform connection",
		Icon:        "kafka",
		Category:    "messaging",
		ConfigSpec:  connection.GenerateConfigSpec(KafkaConfig{}),
	}, &Source{})
}
