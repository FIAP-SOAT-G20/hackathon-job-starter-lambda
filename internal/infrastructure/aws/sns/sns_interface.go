package sns

import "context"

// SNSInterface defines the contract for SNS operations
type SNSInterface interface {
	Publish(ctx context.Context, message string) error
}
