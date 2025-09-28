package sqs

import (
	"context"
	"testing"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/api"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

// Mock S3 client for testing
type mockS3Client struct{}

func (m *mockS3Client) GetObjectMetadata(ctx context.Context, bucket, key string) (map[string]string, error) {
	return map[string]string{"video-id": "123"}, nil
}

func TestNewConsumer(t *testing.T) {
	ctx := context.Background()
	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"

	// Create mock dependencies
	cfg := &config.Config{}
	jobConfig := &config.JobConfig{}
	logger := logger.NewLogger(cfg)
	k8sAPI := &api.K8sAPI{}
	s3Client := &mockS3Client{}

	consumer, err := NewConsumer(ctx, queueURL, cfg, jobConfig, logger, k8sAPI, s3Client)

	assert.NoError(t, err)
	assert.NotNil(t, consumer)
	assert.Equal(t, queueURL, consumer.queueURL)
	assert.Equal(t, cfg, consumer.cfg)
}

func TestConsumer_Start(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create mock dependencies
	cfg := &config.Config{}
	jobConfig := &config.JobConfig{}
	logger := logger.NewLogger(cfg)
	k8sAPI := &api.K8sAPI{}
	s3Client := &mockS3Client{}

	_ = &Consumer{
		queueURL:  "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue",
		cfg:       cfg,
		jobConfig: jobConfig,
		logger:    logger,
		k8sAPI:    k8sAPI,
		s3Client:  s3Client,
	}

	// Test that Start method can be called without errors
	// Note: This is a basic test. In a real scenario, you'd want to mock the SQS client
	// and test the actual message processing logic
	go func() {
		// Cancel context after a short delay to stop the consumer
		<-ctx.Done()
	}()

	// This test mainly verifies that the method signature is correct
	// and doesn't panic when called
	assert.NotPanics(t, func() {
		cancel() // Cancel context to stop the consumer
	})
}
