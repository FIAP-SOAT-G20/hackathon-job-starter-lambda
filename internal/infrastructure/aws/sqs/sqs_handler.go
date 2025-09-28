package sqs

import (
	"context"
	"fmt"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/logger"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SqsHandler struct {
	sqsClient       *SqsClient
	queueURL        string
	maxMessages     int
	waitTimeSeconds int
	logger          *logger.Logger
}

// NewSqsHandler creates a new SQS handler with the provided SQS client
func NewSqsHandler(sqsClient *SqsClient, queueURL string, maxMessages int, waitTimeSeconds int, logger *logger.Logger) *SqsHandler {
	return &SqsHandler{
		sqsClient:       sqsClient,
		logger:          logger,
		queueURL:        queueURL,
		maxMessages:     maxMessages,
		waitTimeSeconds: waitTimeSeconds,
	}
}

// ReceiveMessages processes messages and optionally deletes them after processing
func (h *SqsHandler) ReceiveMessages(ctx context.Context, processor func(types.Message) (bool, error)) error {
	messages, err := h.sqsClient.ReceiveMessages(ctx, h.queueURL, h.maxMessages, h.waitTimeSeconds)
	if err != nil {
		return fmt.Errorf("failed to receive messages: %w", err)
	}

	for _, message := range messages {
		if reprocess, err := processor(message); err != nil {
			h.logger.Error("Failed to process message",
				"error", err.Error(),
				"messageId", *message.MessageId,
				"messageBody", *message.Body,
			)

			if reprocess {
				h.logger.Info("Reprocessing message", "messageId", *message.MessageId)
				continue // Reprocess the message
			}

			h.logger.Info("Sending message to dead-letter queue", "messageId", *message.MessageId)
			// TODO: send the message to a dead-letter queue for further investigation
		}

		if err := h.DeleteMessage(ctx, h.queueURL, *message.ReceiptHandle); err != nil {
			h.logger.Error("Failed to delete message",
				"error", err.Error(),
				"messageId", *message.MessageId,
				"messageBody", *message.Body,
			)
			continue
		}
	}

	return nil
}

// DeleteMessage deletes a specific message from the queue
func (h *SqsHandler) DeleteMessage(ctx context.Context, queueURL string, receiptHandle string) error {
	return h.sqsClient.DeleteMessage(ctx, queueURL, receiptHandle)
}
