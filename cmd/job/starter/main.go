package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/api"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/aws/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func init() {
	// Pre SQS consumer initialization
}

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

	sqsClient, err := sqs.NewSqsClient(infra.AWSClientFactory)
	if err != nil {
		infra.Logger.Error("Failed to create SQS client", "error", err.Error())
		os.Exit(1)
	}

	sqsHandler := sqs.NewSqsHandler(
		sqsClient,
		infra.Config.AWS.SQS.QueueURL,
		infra.Config.AWS.SQS.MaxMessagesBatch,
		infra.Config.AWS.SQS.WaitTimeSeconds,
		infra.Logger,
	)

	infra.Logger.Info("Starting SQS consumer", "queueURL", infra.Config.AWS.SQS.QueueURL)

	// Receive messages from SQS
	for {
		err = sqsHandler.ReceiveMessages(ctx, func(message types.Message) (bool, error) {
			infra.Logger.Info("Processing message", "message", message)

			var s3Event S3Event
			if err := json.Unmarshal([]byte(*message.Body), &s3Event); err != nil {
				return false, fmt.Errorf("failed to unmarshal S3 event: %s", err.Error())
			}

			for _, record := range s3Event.Records {
				err := processS3Record(ctx, infra, record)
				if err != nil {
					infra.Logger.Error("Failed to process message", "error", err.Error(), "messageID", *message.MessageId)
					return false, err
				}
			}

			return true, nil
		})
		if err != nil {
			infra.Logger.Error("Failed to receive messages", "error", err.Error())
		}
	}
}

func processS3Record(ctx context.Context, infra *infrastructure.Infrastructure, record S3EventRecord) error {
	infra.Logger.InfoContext(ctx, "Processing S3 record", "key", record.S3.Object.Key, "bucket", record.S3.Bucket.Name)

	// Get object metadata
	metadata, err := infra.S3.GetObjectMetadata(ctx, record.S3.Bucket.Name, record.S3.Object.Key)
	if err != nil {
		return fmt.Errorf("error getting object metadata: %s", err.Error())
	}

	infra.Logger.InfoContext(ctx, "Object metadata", "metadata", metadata)

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
	jobName := fmt.Sprintf("%s-%s", infra.Config.K8S.Job.Prefix, fileNameWithoutExtension)
	jobCheckerName := fmt.Sprintf("%s-%s-checker", infra.Config.K8S.Job.Prefix, fileNameWithoutExtension)

	// Create job checker
	infra.Logger.InfoContext(ctx, "Creating job checker", "jobName", jobCheckerName)
	err = infra.K8sAPI.CreateJob(ctx, &api.JobInput{
		Namespace:          infra.Config.K8S.Namespace,
		JobName:            jobCheckerName,
		Image:              infra.Config.K8S.Job.ImageChecker,
		Cmd:                infra.Config.K8S.Job.Command,
		ServiceAccountName: infra.Config.K8S.ServiceAccountName,
		Envs: map[string]string{
			"JOB_NAME":                           jobName,
			"JOB_NAMESPACE":                      infra.Config.K8S.Namespace,
			"JOB_VIDEO_ID":                       strconv.FormatInt(videoId, 10),
			"JOB_USER_ID":                        strconv.FormatInt(userId, 10),
			"AWS_ACCESS_KEY_ID":                  infra.Config.AWS.AccessKey,
			"AWS_SECRET_ACCESS_KEY":              infra.Config.AWS.SecretAccessKey,
			"AWS_SESSION_TOKEN":                  infra.Config.AWS.SessionToken,
			"AWS_REGION":                         infra.Config.AWS.Region,
			"AWS_SNS_TOPIC_ARN":                  infra.Config.AWS.SNS.TopicArn,
			"AWS_SQS_QUEUE_URL":                  infra.Config.AWS.SQS.QueueURL,
			"K8S_NAMESPACE":                      infra.Config.K8S.Namespace,
			"K8S_JOB_NAME":                       jobName,
			"K8S_JOB_IMAGE":                      infra.Config.K8S.Job.Image,
			"K8S_JOB_COMMAND":                    infra.Config.K8S.Job.Command,
			"K8S_JOB_PREFIX":                     infra.Config.K8S.Job.Prefix,
			"K8S_JOB_BACK_OFF_LIMIT":             strconv.FormatInt(int64(infra.Config.K8S.Job.BackOffLimit), 10),
			"K8S_JOB_IMAGE_CHECKER":              infra.Config.K8S.Job.ImageChecker,
			"K8S_JOB_TTL_SECONDS_AFTER_FINISHED": strconv.FormatInt(int64(infra.Config.K8S.Job.TtlSecondsAfterFinished.Seconds()), 10),
		},
		TtlSecondsAfterFinished: infra.Config.K8S.Job.TtlSecondsAfterFinished,
	})
	if err != nil {
		return fmt.Errorf("error creating job checker: %s", err.Error())
	}

	// Create main job
	infra.Logger.InfoContext(ctx, "Creating job", "jobName", jobName)
	err = infra.K8sAPI.CreateJob(ctx, &api.JobInput{
		Namespace: infra.Config.K8S.Namespace,
		JobName:   jobName,
		Image:     infra.Config.K8S.Job.Image,
		Cmd:       infra.Config.K8S.Job.Command,
		Envs: map[string]string{
			"VIDEO_KEY":             record.S3.Object.Key,
			"VIDEO_BUCKET":          record.S3.Bucket.Name,
			"PROCESSED_BUCKET":      record.S3.Bucket.Name,
			"AWS_ACCESS_KEY_ID":     infra.Config.AWS.AccessKey,
			"AWS_SECRET_ACCESS_KEY": infra.Config.AWS.SecretAccessKey,
			"AWS_SESSION_TOKEN":     infra.Config.AWS.SessionToken,
			"AWS_REGION":            infra.Config.AWS.Region,
		},
		TtlSecondsAfterFinished: infra.Config.K8S.Job.TtlSecondsAfterFinished,
	})
	if err != nil {
		return fmt.Errorf("error creating job: %s", err.Error())
	}

	return nil
}
