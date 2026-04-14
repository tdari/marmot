// +marmot:name=Kafka
// +marmot:description=This plugin discovers Kafka topics from Kafka clusters.
// +marmot:status=experimental
// +marmot:features=Assets
package kafka

//go:generate go run ../../../docgen/cmd/main.go

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/marmotdata/marmot/internal/core/asset"
	"github.com/marmotdata/marmot/internal/core/connection/providers/kafka"
	"github.com/marmotdata/marmot/internal/plugin"
	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

// Type aliases so the rest of the package (auth.go, client.go) continues
// to use unqualified names while the authoritative definitions live in
// the shared connection/kafka package.
type (
	AuthConfig           = kafka.KafkaAuthConfig
	TLSConfig            = kafka.KafkaTLSConfig
	SchemaRegistryConfig = kafka.KafkaSchemaRegistryConfig
)

// Config for Kafka plugin (discovery/pipeline fields only).
// Connection fields (bootstrap_servers, auth, TLS, etc.) are provided
// via the associated Connection and merged at runtime.
// +marmot:config
type Config struct {
	plugin.BaseConfig `json:",inline"`

	IncludePartitionInfo bool `json:"include_partition_info" description:"Whether to include partition information in metadata" default:"true"`
	IncludeTopicConfig   bool `json:"include_topic_config" description:"Whether to include topic configuration in metadata" default:"true"`
}

// Example configuration for the plugin
// +marmot:example-config
var _ = `
include_partition_info: true
include_topic_config: true
tags:
  - "kafka"
  - "streaming"
`

type Source struct {
	config         *Config
	connConfig     *kafka.KafkaConfig
	client         *kgo.Client
	admin          *kadm.Client
	schemaRegistry schemaregistry.Client
}

func (c *Config) ApplyDefaults() {
	c.IncludePartitionInfo = true
	c.IncludeTopicConfig = true
}

func (s *Source) Validate(rawConfig plugin.RawPluginConfig) (plugin.RawPluginConfig, error) {
	config, err := plugin.UnmarshalPluginConfig[Config](rawConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	config.ApplyDefaults()

	if err := plugin.ValidateStruct(config); err != nil {
		return nil, err
	}

	connConfig, err := plugin.UnmarshalPluginConfig[kafka.KafkaConfig](rawConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling connection config: %w", err)
	}

	s.config = config
	s.connConfig = connConfig
	return rawConfig, nil
}

func (s *Source) Discover(ctx context.Context, pluginConfig plugin.RawPluginConfig) (*plugin.DiscoveryResult, error) {
	if err := s.initClient(ctx); err != nil {
		return nil, fmt.Errorf("initializing Kafka client: %w", err)
	}
	defer s.closeClient()

	if s.connConfig.SchemaRegistry != nil && s.connConfig.SchemaRegistry.Enabled {
		if err := s.initSchemaRegistry(); err != nil {
			log.Warn().Err(err).Msg("Failed to initialize Schema Registry client")
		}
	}

	topics, err := s.discoverTopics(ctx)
	if err != nil {
		return nil, fmt.Errorf("discovering topics: %w", err)
	}

	var assets []asset.Asset
	for _, topic := range topics {
		asset, err := s.createTopicAsset(ctx, topic)
		if err != nil {
			log.Warn().Err(err).Str("topic", topic).Msg("Failed to create asset for topic")
			continue
		}
		assets = append(assets, asset)
	}

	return &plugin.DiscoveryResult{
		Assets: assets,
	}, nil
}

func init() {
	meta := plugin.PluginMeta{
		ID:          "kafka",
		Name:        "Kafka",
		Description: "Discover Kafka topics from Kafka clusters",
		Icon:        "kafka",
		Category:    "streaming",
		ConfigSpec:  plugin.GenerateConfigSpec(Config{}),
	}

	if err := plugin.GetRegistry().Register(meta, &Source{}); err != nil {
		log.Fatal().Err(err).Msg("Failed to register Kafka plugin")
	}
}
