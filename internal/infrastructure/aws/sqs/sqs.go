package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/api"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// S3Event represents the S3 event structure from SQS message
type S3Event struct {
	Records []S3EventRecord `json:"Records"`
}

type S3EventRecord struct {
	S3 S3Record `json:"s3"`
}

type S3Record struct {
	Bucket S3Bucket `json:"bucket"`
	Object S3Object `json:"object"`
}

type S3Bucket struct {
	Name string `json:"name"`
}

type S3Object struct {
	Key string `json:"key"`
}

// SQSMessage represents the SQS message structure
type SQSMessage struct {
	MessageId     string `json:"MessageId"`
	ReceiptHandle string `json:"ReceiptHandle"`
	Body          string `json:"Body"`
}

// Consumer handles SQS message consumption
type Consumer struct {
	client           *sqs.Client
	queueURL         string
	cfg              *config.Config
	jobConfig        *config.JobConfig
	logger           *logger.Logger
	k8sAPI           *api.K8sAPI
	s3Client         S3Client
	workerPoolSize   int
	maxMessagesBatch int
}

// S3Client interface for S3 operations
type S3Client interface {
	GetObjectMetadata(ctx context.Context, bucket, key string) (map[string]string, error)
}

// NewConsumer creates a new SQS consumer
func NewConsumer(ctx context.Context, queueURL string, cfg *config.Config, jobConfig *config.JobConfig, logger *logger.Logger, k8sAPI *api.K8sAPI, s3Client S3Client) (*Consumer, error) {
	client := sqs.NewFromConfig(aws.Config{Region: cfg.AWS.Region, Credentials: credentials.NewStaticCredentialsProvider(cfg.AWS.AccessKey, cfg.AWS.SecretAccessKey, cfg.AWS.SessionToken)})

	return &Consumer{
		client:           client,
		queueURL:         queueURL,
		cfg:              cfg,
		jobConfig:        jobConfig,
		logger:           logger,
		k8sAPI:           k8sAPI,
		s3Client:         s3Client,
		workerPoolSize:   cfg.AWS.SQS.WorkerPoolSize,
		maxMessagesBatch: cfg.AWS.SQS.MaxMessagesBatch,
	}, nil
}

// Start begins consuming messages from the SQS queue
func (c *Consumer) Start(ctx context.Context) error {
	c.logger.InfoContext(ctx, "🟢 SQS Consumer is ready to receive messages!", "queue_url", c.queueURL)

	for {
		select {
		case <-ctx.Done():
			c.logger.InfoContext(ctx, "SQS Consumer stopped")
			return ctx.Err()
		default:
			if err := c.processMessages(ctx); err != nil {
				c.logger.ErrorContext(ctx, "Error processing messages", "error", err.Error())
				// Continue processing even if there's an error
				time.Sleep(5 * time.Second)
			}
		}
	}
}

// processMessages polls and processes messages from the queue
func (c *Consumer) processMessages(ctx context.Context) error {
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: int32(c.maxMessagesBatch),
		WaitTimeSeconds:     int32(c.cfg.AWS.SQS.WaitTimeSeconds), // Long polling
	}

	result, err := c.client.ReceiveMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to receive messages: %s", err.Error())
	}

	if len(result.Messages) == 0 {
		return nil // No messages to process
	}

	c.logger.InfoContext(ctx, "Processing messages in parallel", "count", len(result.Messages), "workers", c.workerPoolSize)

	// Create channels for worker pool
	messageChan := make(chan types.Message, len(result.Messages))
	errorChan := make(chan error)
	successChan := make(chan string) // receipt handles for successful messages

	// Start worker pool
	var wg sync.WaitGroup
	for i := 0; i < c.workerPoolSize; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			c.worker(ctx, workerID, messageChan, errorChan, successChan)
		}(i)
	}

	// Send messages to workers
	for _, message := range result.Messages {
		messageChan <- message
	}
	close(messageChan)

	// Wait for all workers to complete
	wg.Wait()
	close(errorChan)
	close(successChan)

	// Collect results
	var successfulReceiptHandles []string
	var errors []error

	// Collect successful receipt handles
	for receiptHandle := range successChan {
		successfulReceiptHandles = append(successfulReceiptHandles, receiptHandle)
	}

	// Collect errors
	for err := range errorChan {
		errors = append(errors, err)
	}

	// Log results
	c.logger.InfoContext(ctx, "Message processing completed",
		"successful", len(successfulReceiptHandles),
		"failed", len(errors))

	// Delete successful messages in batch
	if len(successfulReceiptHandles) > 0 {
		if err := c.deleteMessagesBatch(ctx, successfulReceiptHandles); err != nil {
			c.logger.ErrorContext(ctx, "Error deleting messages in batch", "error", err.Error())
		}
	}

	// Log errors
	for _, err := range errors {
		c.logger.ErrorContext(ctx, "Message processing error", "error", err.Error())
	}

	return nil
}

// worker processes messages from the message channel
func (c *Consumer) worker(ctx context.Context, workerID int, messageChan <-chan types.Message, errorChan chan<- error, successChan chan<- string) {
	for message := range messageChan {
		c.logger.DebugContext(ctx, "Worker processing message", "worker_id", workerID, "message_id", *message.MessageId)

		// Try to process the message with retry logic
		if err := c.processMessageWithRetry(ctx, message, 3); err != nil {
			c.logger.ErrorContext(ctx, "Worker error processing message after retries",
				"worker_id", workerID,
				"message_id", *message.MessageId,
				"error", err.Error())
			errorChan <- err
		} else {
			// Send receipt handle for successful processing
			successChan <- *message.ReceiptHandle
		}
	}
}

// processMessageWithRetry processes a message with retry logic
func (c *Consumer) processMessageWithRetry(ctx context.Context, message types.Message, maxRetries int) error {
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt > 1 {
			c.logger.InfoContext(ctx, "Retrying message processing",
				"message_id", *message.MessageId,
				"attempt", attempt,
				"max_retries", maxRetries)

			// Exponential backoff: wait 2^attempt seconds
			backoffDuration := time.Duration(1<<uint(attempt-1)) * time.Second
			select {
			case <-time.After(backoffDuration):
				// continue after backoff
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if err := c.processMessage(ctx, message); err != nil {
			lastErr = err
			c.logger.WarnContext(ctx, "Message processing attempt failed",
				"message_id", *message.MessageId,
				"attempt", attempt,
				"error", err.Error())
			continue
		}

		// Success
		if attempt > 1 {
			c.logger.InfoContext(ctx, "Message processing succeeded after retry",
				"message_id", *message.MessageId,
				"attempt", attempt)
		}
		return nil
	}

	return fmt.Errorf("message processing failed after %d attempts: %s", maxRetries, lastErr.Error())
}

// processMessage processes a single SQS message
func (c *Consumer) processMessage(ctx context.Context, message types.Message) error {
	c.logger.InfoContext(ctx, "Processing message", "message_id", *message.MessageId)

	// Parse the S3 event from the message body
	var s3Event S3Event
	if err := json.Unmarshal([]byte(*message.Body), &s3Event); err != nil {
		return fmt.Errorf("failed to unmarshal S3 event: %s", err.Error())
	}

	c.logger.InfoContext(ctx, "Processing S3 event", "records_count", len(s3Event.Records))

	// Process each S3 record
	for _, record := range s3Event.Records {
		if err := c.processS3Record(ctx, record); err != nil {
			return fmt.Errorf("failed to process S3 record: %s", err.Error())
		}
	}

	return nil
}

// processS3Record processes a single S3 record (similar to the original Lambda handler)
func (c *Consumer) processS3Record(ctx context.Context, record S3EventRecord) error {
	c.logger.InfoContext(ctx, "Processing S3 record", "key", record.S3.Object.Key, "bucket", record.S3.Bucket.Name)

	// Get object metadata
	metadata, err := c.s3Client.GetObjectMetadata(ctx, record.S3.Bucket.Name, record.S3.Object.Key)
	if err != nil {
		return fmt.Errorf("error getting object metadata: %s", err.Error())
	}

	c.logger.InfoContext(ctx, "Object metadata", "metadata", metadata)

	// Parse video ID
	videoId, err := strconv.ParseInt(metadata["video-id"], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing video id: %s", err.Error())
	}

	// Parse video ID
	userId, err := strconv.ParseInt(metadata["user-id"], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing user id: %s", err.Error())
	}

	// Generate job names
	splittedKey := strings.Split(record.S3.Object.Key, "/")
	fileName := splittedKey[len(splittedKey)-1]
	fileNameWithoutExtension := strings.Split(fileName, ".")[0]
	jobName := fmt.Sprintf("%s-%s", c.cfg.K8S.Job.Prefix, fileNameWithoutExtension)
	jobCheckerName := fmt.Sprintf("%s-%s-checker", c.cfg.K8S.Job.Prefix, fileNameWithoutExtension)

	// Create job checker
	c.logger.InfoContext(ctx, "Creating job checker", "jobName", jobCheckerName)
	err = c.k8sAPI.CreateJob(ctx, &api.JobInput{
		Namespace:          c.cfg.K8S.Namespace,
		JobName:            jobCheckerName,
		Image:              c.cfg.K8S.Job.ImageChecker,
		Cmd:                c.cfg.K8S.Job.Command,
		ServiceAccountName: c.cfg.K8S.ServiceAccountName,
		Envs: map[string]string{
			"JOB_NAME":                           jobName,
			"JOB_NAMESPACE":                      c.cfg.K8S.Namespace,
			"JOB_VIDEO_ID":                       strconv.FormatInt(videoId, 10),
			"JOB_USER_ID":                        strconv.FormatInt(userId, 10),
			"AWS_ACCESS_KEY_ID":                  c.cfg.AWS.AccessKey,
			"AWS_SECRET_ACCESS_KEY":              c.cfg.AWS.SecretAccessKey,
			"AWS_SESSION_TOKEN":                  c.cfg.AWS.SessionToken,
			"AWS_REGION":                         c.cfg.AWS.Region,
			"AWS_SNS_TOPIC_ARN":                  c.cfg.AWS.SNS.TopicArn,
			"AWS_SQS_QUEUE_URL":                  c.cfg.AWS.SQS.QueueURL,
			"K8S_NAMESPACE":                      c.cfg.K8S.Namespace,
			"K8S_JOB_NAME":                       jobName,
			"K8S_JOB_IMAGE":                      c.cfg.K8S.Job.Image,
			"K8S_JOB_COMMAND":                    c.cfg.K8S.Job.Command,
			"K8S_JOB_PREFIX":                     c.cfg.K8S.Job.Prefix,
			"K8S_JOB_BACK_OFF_LIMIT":             strconv.FormatInt(int64(c.cfg.K8S.Job.BackOffLimit), 10),
			"K8S_JOB_IMAGE_CHECKER":              c.cfg.K8S.Job.ImageChecker,
			"K8S_JOB_TTL_SECONDS_AFTER_FINISHED": strconv.FormatInt(int64(c.cfg.K8S.Job.TtlSecondsAfterFinished.Seconds()), 10),
		},
		TtlSecondsAfterFinished: c.cfg.K8S.Job.TtlSecondsAfterFinished,
	})
	if err != nil {
		return fmt.Errorf("error creating job checker: %s", err.Error())
	}

	// Create main job
	c.logger.InfoContext(ctx, "Creating job", "jobName", jobName)
	err = c.k8sAPI.CreateJob(ctx, &api.JobInput{
		Namespace: c.cfg.K8S.Namespace,
		JobName:   jobName,
		Image:     c.cfg.K8S.Job.Image,
		Cmd:       c.cfg.K8S.Job.Command,
		Envs: map[string]string{
			"VIDEO_KEY":             record.S3.Object.Key,
			"VIDEO_BUCKET":          record.S3.Bucket.Name,
			"PROCEDSSED_BUCKET":     record.S3.Bucket.Name,
			"AWS_ACCESS_KEY_ID":     c.cfg.AWS.AccessKey,
			"AWS_SECRET_ACCESS_KEY": c.cfg.AWS.SecretAccessKey,
			"AWS_SESSION_TOKEN":     c.cfg.AWS.SessionToken,
			"AWS_REGION":            c.cfg.AWS.Region,
		},
		TtlSecondsAfterFinished: c.cfg.K8S.Job.TtlSecondsAfterFinished,
	})
	if err != nil {
		return fmt.Errorf("error creating job: %s", err.Error())
	}

	return nil
}

// deleteMessage deletes a processed message from the queue
// func (c *Consumer) deleteMessage(ctx context.Context, receiptHandle string) error {
// 	input := &sqs.DeleteMessageInput{
// 		QueueUrl:      aws.String(c.queueURL),
// 		ReceiptHandle: aws.String(receiptHandle),
// 	}

// 	_, err := c.client.DeleteMessage(ctx, input)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete message: %s", err.Error())
// 	}

// 	return nil
// }

// deleteMessagesBatch deletes multiple messages from the queue in batch
func (c *Consumer) deleteMessagesBatch(ctx context.Context, receiptHandles []string) error {
	if len(receiptHandles) == 0 {
		return nil
	}

	// SQS batch delete can handle up to 10 messages at a time
	const maxBatchSize = 10

	for i := 0; i < len(receiptHandles); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(receiptHandles) {
			end = len(receiptHandles)
		}

		batch := receiptHandles[i:end]
		if err := c.deleteBatch(ctx, batch); err != nil {
			return fmt.Errorf("failed to delete batch %d-%d: %s", i, end-1, err.Error())
		}
	}

	return nil
}

// deleteBatch deletes a batch of messages (up to 10)
func (c *Consumer) deleteBatch(ctx context.Context, receiptHandles []string) error {
	if len(receiptHandles) == 0 {
		return nil
	}

	entries := make([]types.DeleteMessageBatchRequestEntry, len(receiptHandles))
	for i, receiptHandle := range receiptHandles {
		entries[i] = types.DeleteMessageBatchRequestEntry{
			Id:            aws.String(fmt.Sprintf("msg-%d", i)),
			ReceiptHandle: aws.String(receiptHandle),
		}
	}

	input := &sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(c.queueURL),
		Entries:  entries,
	}

	result, err := c.client.DeleteMessageBatch(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete message batch: %s", err.Error())
	}

	// Log any failed deletions
	if len(result.Failed) > 0 {
		for _, failed := range result.Failed {
			c.logger.ErrorContext(ctx, "Failed to delete message in batch",
				"id", *failed.Id,
				"code", *failed.Code,
				"message", *failed.Message)
		}
	}

	return nil
}
