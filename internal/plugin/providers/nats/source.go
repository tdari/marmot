// +marmot:name=NATS
// +marmot:description=Discovers JetStream streams from NATS servers.
// +marmot:status=experimental
// +marmot:features=Assets
package nats

//go:generate go run ../../../docgen/cmd/main.go

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/marmotdata/marmot/internal/core/asset"
	connectionnats "github.com/marmotdata/marmot/internal/core/connection/providers/nats"
	"github.com/marmotdata/marmot/internal/mrn"
	"github.com/marmotdata/marmot/internal/plugin"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
)

// Config for NATS plugin (discovery/pipeline fields only).
// Connection fields (host, port, token, username, password, credentials_file, tls, tls_insecure)
// are provided via the associated Connection and merged at runtime.
// +marmot:config
type Config struct {
	plugin.BaseConfig `json:",inline"`
}

// Example configuration for the plugin
// +marmot:example-config
var _ = `
filter:
  include:
    - "^ORDERS"
tags:
  - "nats"
  - "messaging"
`

type Source struct {
	config     *Config
	connConfig *connectionnats.NatsConfig
}

func (s *Source) Validate(rawConfig plugin.RawPluginConfig) (plugin.RawPluginConfig, error) {
	config, err := plugin.UnmarshalPluginConfig[Config](rawConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	if err := plugin.ValidateStruct(config); err != nil {
		return nil, err
	}

	s.config = config
	return rawConfig, nil
}

func (s *Source) Discover(ctx context.Context, pluginConfig plugin.RawPluginConfig) (*plugin.DiscoveryResult, error) {
	connConfig, err := plugin.UnmarshalPluginConfig[connectionnats.NatsConfig](pluginConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling connection config: %w", err)
	}

	s.connConfig = connConfig

	nc, err := s.connect()
	if err != nil {
		return nil, fmt.Errorf("connecting to NATS: %w", err)
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("creating JetStream context: %w", err)
	}

	var assets []asset.Asset

	streams := js.ListStreams(ctx)
	for info := range streams.Info() {
		a := s.createStreamAsset(info)
		assets = append(assets, a)
	}
	if err := streams.Err(); err != nil {
		return nil, fmt.Errorf("listing streams: %w", err)
	}

	return &plugin.DiscoveryResult{
		Assets: assets,
	}, nil
}

func (s *Source) connect() (*nats.Conn, error) {
	addr := fmt.Sprintf("nats://%s:%d", s.connConfig.Host, s.connConfig.Port)

	opts := []nats.Option{
		nats.Timeout(10 * time.Second),
		nats.Name("marmot-discovery"),
	}

	if s.connConfig.Token != "" {
		opts = append(opts, nats.Token(s.connConfig.Token))
	}

	if s.connConfig.Username != "" && s.connConfig.Password != "" {
		opts = append(opts, nats.UserInfo(s.connConfig.Username, s.connConfig.Password))
	}

	if s.connConfig.CredentialsFile != "" {
		opts = append(opts, nats.UserCredentials(s.connConfig.CredentialsFile))
	}

	if s.connConfig.TLS {
		opts = append(opts, nats.Secure(&tls.Config{
			InsecureSkipVerify: s.connConfig.TLSInsecure, //nolint:gosec // G402: user-controlled TLS config
		}))
	}

	return nats.Connect(addr, opts...)
}

func (s *Source) createStreamAsset(info *jetstream.StreamInfo) asset.Asset {
	metadata := make(map[string]interface{})

	metadata["stream_name"] = info.Config.Name
	metadata["subjects"] = strings.Join(info.Config.Subjects, ", ")
	metadata["retention_policy"] = info.Config.Retention.String()
	metadata["max_bytes"] = info.Config.MaxBytes
	metadata["max_msgs"] = info.Config.MaxMsgs
	metadata["max_age"] = info.Config.MaxAge.String()
	metadata["max_msg_size"] = int64(info.Config.MaxMsgSize)
	metadata["storage_type"] = info.Config.Storage.String()
	metadata["num_replicas"] = info.Config.Replicas
	metadata["duplicate_window"] = info.Config.Duplicates.String()
	metadata["discard_policy"] = info.Config.Discard.String()

	metadata["messages"] = info.State.Msgs
	metadata["bytes"] = info.State.Bytes
	metadata["consumer_count"] = info.State.Consumers
	metadata["first_seq"] = info.State.FirstSeq
	metadata["last_seq"] = info.State.LastSeq

	metadata["host"] = s.connConfig.Host
	metadata["port"] = s.connConfig.Port

	streamName := info.Config.Name
	mrnValue := mrn.New("Stream", "NATS", streamName)

	processedTags := plugin.InterpolateTags(s.config.Tags, metadata)

	return asset.Asset{
		Name:      &streamName,
		MRN:       &mrnValue,
		Type:      "Stream",
		Providers: []string{"NATS"},
		Metadata:  metadata,
		Tags:      processedTags,
		Sources: []asset.AssetSource{{
			Name:       "NATS",
			LastSyncAt: time.Now(),
			Properties: metadata,
			Priority:   1,
		}},
	}
}

func init() {
	meta := plugin.PluginMeta{
		ID:          "nats",
		Name:        "NATS",
		Description: "Discover JetStream streams from NATS servers",
		Icon:        "nats",
		Category:    "messaging",
		ConfigSpec:  plugin.GenerateConfigSpec(Config{}),
	}

	if err := plugin.GetRegistry().Register(meta, &Source{}); err != nil {
		log.Fatal().Err(err).Msg("Failed to register NATS plugin")
	}
}
