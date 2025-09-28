package aws

import (
	"context"
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// ClientFactory provides AWS clients with shared configuration
type ClientFactory struct {
	config awssdk.Config
}

// NewClientFactory creates a new AWS client factory with shared configuration
func NewClientFactory(ctx context.Context, region string) (*ClientFactory, error) {
	var awsConfig awssdk.Config
	var err error

	// Check if AWS environment variables are set
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if accessKeyID != "" && secretAccessKey != "" {
		// Use explicit AWS credentials from environment variables
		awsConfig, err = awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(region),
			awsconfig.WithCredentialsProvider(
				credentials.StaticCredentialsProvider{
					Value: awssdk.Credentials{
						AccessKeyID:     accessKeyID,
						SecretAccessKey: secretAccessKey,
						SessionToken:    os.Getenv("AWS_SESSION_TOKEN"),
					},
				},
			),
		)
	} else {
		// Fall back to default AWS configuration (e.g., IAM roles, default profile)
		awsConfig, err = awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(region),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &ClientFactory{
		config: awsConfig,
	}, nil
}

// GetConfig returns the AWS configuration for creating service clients
func (f *ClientFactory) GetConfig() awssdk.Config {
	return f.config
}

// GetRegion returns the configured AWS region
func (f *ClientFactory) GetRegion() string {
	return f.config.Region
}
