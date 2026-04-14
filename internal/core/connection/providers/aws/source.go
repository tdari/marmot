// +marmot:name=AWS
// +marmot:description=AWS connection
// +marmot:status=stable
// +marmot:category=cloud
package aws

//go:generate go run ../../../../docgen/cmd/connection/main.go

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/marmotdata/marmot/internal/core/connection"
)

// +marmot:config
type AWSConfig struct {
	Region         string `json:"region" label:"Region" description:"AWS region" validate:"required" placeholder:"us-east-1"`
	ID             string `json:"id" label:"Access Key ID" description:"AWS access key ID" sensitive:"true" placeholder:"AKIAIOSFODNN7EXAMPLE"`
	Secret         string `json:"secret" label:"Secret Access Key" description:"AWS secret access key" sensitive:"true" placeholder:"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"`
	Token          string `json:"token,omitempty" label:"Session Token" description:"AWS session token for temporary credentials" sensitive:"true"`
	UseDefault     bool   `json:"use_default" label:"Use Default Credentials" description:"Use default AWS credentials from environment/config files" default:"false"`
	Profile        string `json:"profile" label:"Profile" description:"AWS profile name" placeholder:"default" show_when:"use_default:false"`
	Role           string `json:"role" label:"Role ARN" description:"AWS IAM role ARN to assume" placeholder:"arn:aws:iam::123456789012:role/MyRole" show_when:"use_default:false"`
	RoleExternalID string `json:"role_external_id" label:"Role External ID" description:"External ID for role assumption" placeholder:"external-id" show_when:"use_default:false"`
	Endpoint       string `json:"endpoint" label:"Custom Endpoint" description:"Custom S3-compatible endpoint URL" placeholder:"https://s3.amazonaws.com"`
}

// +marmot:example-config
var _ = `
region: us-east-1
id: your-access-key-id
secret: your-api-secret
use_default: false
`

func (a *AWSConfig) NewAWSConfig(ctx context.Context) (aws.Config, error) {
	var opts []func(*config.LoadOptions) error

	if a.Region != "" {
		opts = append(opts, config.WithRegion(a.Region))
	}

	// If UseDefault is true or no explicit credentials provided, use default credential chain
	if a.UseDefault || (a.ID == "" && a.Profile == "") {
		// Just load default config - will follow AWS credential chain
		if a.Profile != "" {
			opts = append(opts, config.WithSharedConfigProfile(a.Profile))
		}
	} else {
		// Use explicit credentials if provided
		if a.ID != "" && a.Secret != "" {
			provider := credentials.NewStaticCredentialsProvider(
				a.ID,
				a.Secret,
				a.Token,
			)
			opts = append(opts, config.WithCredentialsProvider(provider))
		}

		if a.Profile != "" {
			opts = append(opts, config.WithSharedConfigProfile(a.Profile))
		}
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("loading AWS config: %w", err)
	}

	if a.Role != "" {
		stsClient := sts.NewFromConfig(cfg)
		assumeRoleOpts := func(o *stscreds.AssumeRoleOptions) {
			if a.RoleExternalID != "" {
				o.ExternalID = aws.String(a.RoleExternalID)
			}
		}

		provider := stscreds.NewAssumeRoleProvider(stsClient, a.Role, assumeRoleOpts)
		cfg.Credentials = aws.NewCredentialsCache(provider)
	}

	if a.Endpoint != "" {
		cfg.BaseEndpoint = aws.String(a.Endpoint)
	}

	return cfg, nil
}

type Source struct{}

func (s *Source) Validate(rawConfig map[string]interface{}) error {
	config, err := connection.UnmarshalConfig[AWSConfig](rawConfig)
	if err != nil {
		return err
	}
	return connection.ValidateConfig(config)
}

func init() {
	connection.GetRegistry().Register(connection.ConnectionTypeMeta{
		ID:          "aws",
		Name:        "AWS",
		Description: "AWS connection",
		Icon:        "aws",
		Category:    "cloud",
		ConfigSpec:  connection.GenerateConfigSpec(AWSConfig{}),
	}, &Source{})
}
