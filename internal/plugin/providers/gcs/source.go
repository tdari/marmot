// +marmot:name=Google Cloud Storage
// +marmot:description=Discovers buckets from Google Cloud Storage.
// +marmot:status=experimental
// +marmot:features=Assets
package gcs

//go:generate go run ../../../docgen/cmd/main.go

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/marmotdata/marmot/internal/core/asset"
	"github.com/marmotdata/marmot/internal/core/connection/providers/gcs"
	"github.com/marmotdata/marmot/internal/core/lineage"
	"github.com/marmotdata/marmot/internal/mrn"
	"github.com/marmotdata/marmot/internal/plugin"
	"github.com/rs/zerolog/log"
)

// Config for Google Cloud Storage plugin (discovery/pipeline fields only).
// Connection fields (project_id, credentials_file, credentials_json, endpoint, disable_auth)
// are provided via the associated Connection and merged at runtime.
// +marmot:config
type Config struct {
	plugin.BaseConfig `json:",inline"`

	// Discovery options
	IncludeMetadata    bool `json:"include_metadata" description:"Include bucket metadata like labels" default:"true"`
	IncludeObjectCount bool `json:"include_object_count" description:"Count objects in each bucket (can be slow for large buckets)" default:"false"`
}

// Example configuration for the plugin
// +marmot:example-config
var _ = `
include_metadata: true
include_object_count: false
filter:
  include:
    - "^data-.*"
  exclude:
    - ".*-temp$"
tags:
  - "gcs"
  - "storage"
`

type Source struct {
	config     *Config
	connConfig *gcs.GCSConfig
	client     *storage.Client
}

func (s *Source) Validate(rawConfig plugin.RawPluginConfig) (plugin.RawPluginConfig, error) {
	config, err := plugin.UnmarshalPluginConfig[Config](rawConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	if err := plugin.ValidateStruct(config); err != nil {
		return nil, err
	}

	connConfig, err := plugin.UnmarshalPluginConfig[gcs.GCSConfig](rawConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling connection config: %w", err)
	}

	s.config = config
	s.connConfig = connConfig
	return rawConfig, nil
}

func (s *Source) Discover(ctx context.Context, pluginConfig plugin.RawPluginConfig) (*plugin.DiscoveryResult, error) {
	config, err := plugin.UnmarshalPluginConfig[Config](pluginConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}
	s.config = config

	connConfig, err := plugin.UnmarshalPluginConfig[gcs.GCSConfig](pluginConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling connection config: %w", err)
	}
	s.connConfig = connConfig

	client, err := s.createClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating GCS client: %w", err)
	}
	defer client.Close()
	s.client = client

	buckets, err := s.discoverBuckets(ctx)
	if err != nil {
		return nil, fmt.Errorf("discovering buckets: %w", err)
	}

	var assets []asset.Asset
	var lineages []lineage.LineageEdge

	for _, bucket := range buckets {
		asset, err := s.createBucketAsset(ctx, bucket)
		if err != nil {
			log.Warn().Err(err).Str("bucket", bucket.Name).Msg("Failed to create asset for bucket")
			continue
		}
		assets = append(assets, asset)
	}

	return &plugin.DiscoveryResult{
		Assets:  assets,
		Lineage: lineages,
	}, nil
}

func (s *Source) createClient(ctx context.Context) (*storage.Client, error) {
	var opts []option.ClientOption

	if s.connConfig.Endpoint != "" {
		opts = append(opts, option.WithEndpoint(s.connConfig.Endpoint))
	}

	switch {
	case s.connConfig.DisableAuth:
		opts = append(opts, option.WithoutAuthentication())
	case s.connConfig.CredentialsJSON != "":
		opts = append(opts, option.WithCredentialsJSON([]byte(s.connConfig.CredentialsJSON)))
	case s.connConfig.CredentialsFile != "":
		opts = append(opts, option.WithCredentialsFile(s.connConfig.CredentialsFile))
	}

	return storage.NewClient(ctx, opts...)
}

func (s *Source) discoverBuckets(ctx context.Context) ([]*storage.BucketAttrs, error) {
	var buckets []*storage.BucketAttrs

	it := s.client.Buckets(ctx, s.connConfig.ProjectID)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("iterating buckets: %w", err)
		}
		buckets = append(buckets, attrs)
	}

	return buckets, nil
}

func (s *Source) createBucketAsset(ctx context.Context, bucket *storage.BucketAttrs) (asset.Asset, error) {
	bucketName := bucket.Name

	metadata := make(map[string]interface{})
	metadata["bucket_name"] = bucketName
	metadata["location"] = bucket.Location
	metadata["location_type"] = bucket.LocationType
	metadata["storage_class"] = bucket.StorageClass
	metadata["created"] = bucket.Created.Format(time.RFC3339)

	if bucket.VersioningEnabled {
		metadata["versioning"] = "enabled"
	} else {
		metadata["versioning"] = "disabled"
	}

	if bucket.RequesterPays {
		metadata["requester_pays"] = true
	}

	if bucket.DefaultEventBasedHold {
		metadata["default_event_based_hold"] = true
	}

	if bucket.RetentionPolicy != nil {
		metadata["retention_period_seconds"] = bucket.RetentionPolicy.RetentionPeriod.Seconds()
	}

	if bucket.Encryption != nil && bucket.Encryption.DefaultKMSKeyName != "" {
		metadata["encryption"] = "customer-managed"
		metadata["kms_key"] = bucket.Encryption.DefaultKMSKeyName
	} else {
		metadata["encryption"] = "google-managed"
	}

	if bucket.Logging != nil && bucket.Logging.LogBucket != "" {
		metadata["logging_enabled"] = true
		metadata["log_bucket"] = bucket.Logging.LogBucket
	}

	if s.config.IncludeMetadata && len(bucket.Labels) > 0 {
		for key, value := range bucket.Labels {
			metadata["label_"+key] = value
		}
	}

	if bucket.Lifecycle.Rules != nil && len(bucket.Lifecycle.Rules) > 0 {
		metadata["lifecycle_rules_count"] = len(bucket.Lifecycle.Rules)
	}

	if s.config.IncludeObjectCount {
		count, err := s.countObjects(ctx, bucketName)
		if err != nil {
			log.Warn().Err(err).Str("bucket", bucketName).Msg("Failed to count objects")
		} else {
			metadata["object_count"] = count
		}
	}

	mrnValue := mrn.New("Bucket", "GCS", bucketName)

	processedTags := plugin.InterpolateTags(s.config.Tags, metadata)

	return asset.Asset{
		Name:      &bucketName,
		MRN:       &mrnValue,
		Type:      "Bucket",
		Providers: []string{"GCS"},
		Metadata:  metadata,
		Tags:      processedTags,
		Sources: []asset.AssetSource{{
			Name:       "GCS",
			LastSyncAt: time.Now(),
			Properties: metadata,
			Priority:   1,
		}},
	}, nil
}

func (s *Source) countObjects(ctx context.Context, bucketName string) (int64, error) {
	var count int64

	it := s.client.Bucket(bucketName).Objects(ctx, nil)
	for {
		_, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}
		count++
	}

	return count, nil
}

func init() {
	meta := plugin.PluginMeta{
		ID:          "gcs",
		Name:        "Google Cloud Storage",
		Description: "Discover buckets from Google Cloud Storage",
		Icon:        "gcs",
		Category:    "storage",
		ConfigSpec:  plugin.GenerateConfigSpec(Config{}),
	}

	if err := plugin.GetRegistry().Register(meta, &Source{}); err != nil {
		log.Fatal().Err(err).Msg("Failed to register GCS plugin")
	}
}
