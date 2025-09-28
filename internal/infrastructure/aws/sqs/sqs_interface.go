package sqs

import (
	"context"
)

// ConsumerInterface defines the interface for SQS consumer
type ConsumerInterface interface {
	Start(ctx context.Context) error
}
