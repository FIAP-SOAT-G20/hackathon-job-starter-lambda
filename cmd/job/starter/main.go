package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/aws/sqs"
)

func init() {
	// Pre SQS consumer initialization
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	// Get infrastructure
	infra := infrastructure.GetInfrastructure()

	// Create SQS consumer
	consumer, err := sqs.NewConsumer(ctx, infra.Config.AWS.SQS.QueueURL, infra.Config, infra.JobConfig, infra.Logger, infra.K8sAPI, infra.S3)
	if err != nil {
		infra.Logger.ErrorContext(ctx, "Failed to create SQS consumer", "error", err)
		os.Exit(1)
	}

	// Start consuming messages
	if err := consumer.Start(ctx); err != nil {
		infra.Logger.ErrorContext(ctx, "SQS consumer failed", "error", err)
		os.Exit(1)
	}
}
