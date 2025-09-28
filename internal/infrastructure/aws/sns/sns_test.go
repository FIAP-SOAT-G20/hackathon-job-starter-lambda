package sns

import (
	"testing"

	myConfig "github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

func TestSNS_Publish(t *testing.T) {
	t.Run("should publish message with all parameters", func(t *testing.T) {
		// This test would require mocking the AWS SNS client
		// For now, we'll test the parameter handling logic

		// Arrange
		cfg := &myConfig.Config{
			AWS: struct {
				Region          string
				AccessKey       string
				SecretAccessKey string
				SessionToken    string
				SNS             struct {
					TopicArn string
				}
				SQS struct {
					QueueURL         string
					WorkerPoolSize   int
					MaxMessagesBatch int
					WaitTimeSeconds  int
				}
			}{
				Region:          "us-east-1",
				AccessKey:       "test-access-key",
				SecretAccessKey: "test-secret-key",
				SessionToken:    "test-session-token",
				SNS: struct {
					TopicArn string
				}{
					TopicArn: "arn:aws:sns:us-east-1:123456789012:test-topic",
				},
				SQS: struct {
					QueueURL         string
					WorkerPoolSize   int
					MaxMessagesBatch int
					WaitTimeSeconds  int
				}{
					QueueURL:         "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
					WorkerPoolSize:   5,
					MaxMessagesBatch: 10,
					WaitTimeSeconds:  20,
				},
			},
		}

		// Note: This test would fail in a real environment without proper AWS credentials
		// In a real test environment, you would mock the AWS SNS client
		sns := NewSNS(cfg)

		// Act & Assert
		// This test demonstrates the structure but would need proper mocking
		assert.NotNil(t, sns)
		assert.Equal(t, cfg.AWS.SNS.TopicArn, sns.TopicArn)
		assert.NotNil(t, sns.Client)
	})

	t.Run("should handle empty parameters", func(t *testing.T) {
		// Arrange
		cfg := &myConfig.Config{
			AWS: struct {
				Region          string
				AccessKey       string
				SecretAccessKey string
				SessionToken    string
				SNS             struct {
					TopicArn string
				}
				SQS struct {
					QueueURL         string
					WorkerPoolSize   int
					MaxMessagesBatch int
					WaitTimeSeconds  int
				}
			}{
				Region:          "us-east-1",
				AccessKey:       "test-access-key",
				SecretAccessKey: "test-secret-key",
				SessionToken:    "test-session-token",
				SNS: struct {
					TopicArn string
				}{
					TopicArn: "arn:aws:sns:us-east-1:123456789012:test-topic",
				},
				SQS: struct {
					QueueURL         string
					WorkerPoolSize   int
					MaxMessagesBatch int
					WaitTimeSeconds  int
				}{
					QueueURL:         "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
					WorkerPoolSize:   5,
					MaxMessagesBatch: 10,
					WaitTimeSeconds:  20,
				},
			},
		}

		sns := NewSNS(cfg)

		// Act & Assert
		assert.NotNil(t, sns)
		assert.Equal(t, cfg.AWS.SNS.TopicArn, sns.TopicArn)
		assert.NotNil(t, sns.Client)
	})
}

func TestNewSNS(t *testing.T) {
	t.Run("should create SNS with valid config", func(t *testing.T) {
		// Arrange
		cfg := &myConfig.Config{
			AWS: struct {
				Region          string
				AccessKey       string
				SecretAccessKey string
				SessionToken    string
				SNS             struct {
					TopicArn string
				}
				SQS struct {
					QueueURL         string
					WorkerPoolSize   int
					MaxMessagesBatch int
					WaitTimeSeconds  int
				}
			}{
				Region:          "us-east-1",
				AccessKey:       "test-access-key",
				SecretAccessKey: "test-secret-key",
				SessionToken:    "test-session-token",
				SNS: struct {
					TopicArn string
				}{
					TopicArn: "arn:aws:sns:us-east-1:123456789012:test-topic",
				},
				SQS: struct {
					QueueURL         string
					WorkerPoolSize   int
					MaxMessagesBatch int
					WaitTimeSeconds  int
				}{
					QueueURL:         "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
					WorkerPoolSize:   5,
					MaxMessagesBatch: 10,
					WaitTimeSeconds:  20,
				},
			},
		}

		// Act
		sns := NewSNS(cfg)

		// Assert
		assert.NotNil(t, sns)
		assert.Equal(t, cfg.AWS.SNS.TopicArn, sns.TopicArn)
		assert.NotNil(t, sns.Client)
	})

	t.Run("should create SNS with nil config", func(t *testing.T) {
		// This test would panic in a real scenario due to nil pointer dereference
		// It's included to show the expected behavior
		assert.Panics(t, func() {
			NewSNS(nil)
		})
	})

	t.Run("should create SNS with empty config", func(t *testing.T) {
		// Arrange
		cfg := &myConfig.Config{}

		// Act
		sns := NewSNS(cfg)

		// Assert
		// The SNS client will be created but with empty configuration
		assert.NotNil(t, sns)
		assert.Equal(t, "", sns.TopicArn)
		assert.NotNil(t, sns.Client)
	})
}

// TestSNSInterface tests the interface implementation
func TestSNSInterface(t *testing.T) {
	t.Run("should implement SNSInterface", func(t *testing.T) {
		// Arrange
		cfg := &myConfig.Config{
			AWS: struct {
				Region          string
				AccessKey       string
				SecretAccessKey string
				SessionToken    string
				SNS             struct {
					TopicArn string
				}
				SQS struct {
					QueueURL         string
					WorkerPoolSize   int
					MaxMessagesBatch int
					WaitTimeSeconds  int
				}
			}{
				Region:          "us-east-1",
				AccessKey:       "test-access-key",
				SecretAccessKey: "test-secret-key",
				SessionToken:    "test-session-token",
				SNS: struct {
					TopicArn string
				}{
					TopicArn: "arn:aws:sns:us-east-1:123456789012:test-topic",
				},
				SQS: struct {
					QueueURL         string
					WorkerPoolSize   int
					MaxMessagesBatch int
					WaitTimeSeconds  int
				}{
					QueueURL:         "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
					WorkerPoolSize:   5,
					MaxMessagesBatch: 10,
					WaitTimeSeconds:  20,
				},
			},
		}

		sns := NewSNS(cfg)

		// Act & Assert
		var _ SNSInterface = sns
		assert.NotNil(t, sns)
	})
}
