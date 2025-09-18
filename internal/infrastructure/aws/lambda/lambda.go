package lambda

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/api"
	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-lambda-go/lambda"
)

// StartLambda is the function that tells lambda which function should be call to start lambda.
func StartLambda() {
	fmt.Println("🟢 Lambda is ready to receive requests!")
	lambda.Start(handleRequest)
}

// handleRequest responsible to handle lambda events
func handleRequest(ctx context.Context, req events.S3Event) error {
	infra := infrastructure.GetInfrastructure()
	l := infra.Logger
	cfg := infra.LambdaConfig
	k8sAPI := infra.K8sAPI

	l.InfoContext(ctx, "Starting lambda handler", "records length", len(req.Records))
	for _, record := range req.Records {
		l.InfoContext(ctx, "Processing record", "key", record.S3.Object.Key, "bucket", record.S3.Bucket.Name)

		metadata, err := infra.S3.GetObjectMetadata(ctx, record.S3.Bucket.Name, record.S3.Object.Key)
		if err != nil {
			l.ErrorContext(ctx, "Error getting object metadata", "error", err)
			return err
		}

		videoId, err := strconv.ParseInt(metadata["video-id"], 10, 64)
		if err != nil {
			l.ErrorContext(ctx, "Error parsing video id", "error", err)
			return err
		}

		splittedKey := strings.Split(record.S3.Object.Key, "/")
		fileName := splittedKey[len(splittedKey)-1]
		fileNameWithoutExtension := strings.Split(fileName, ".")[0]
		jobName := fmt.Sprintf("%s-%s", cfg.K8S.Job.Prefix, fileNameWithoutExtension)
		jobCheckerName := fmt.Sprintf("%s-%s-checker", cfg.K8S.Job.Prefix, fileNameWithoutExtension)
		l.InfoContext(ctx, "Creating job checker", "jobName", jobCheckerName)
		err = k8sAPI.CreateJob(ctx, &api.JobInput{
			Namespace:          cfg.K8S.Namespace,
			JobName:            jobCheckerName,
			Image:              cfg.K8S.Job.ImageChecker,
			Cmd:                cfg.K8S.Job.Command,
			ServiceAccountName: cfg.K8S.ServiceAccountName,
			Envs: map[string]string{
				"JOB_NAME":      jobName,
				"JOB_NAMESPACE": cfg.K8S.Namespace,
				"JOB_VIDEO_ID":  strconv.FormatInt(videoId, 10),
				"JOB_USER_ID":   strconv.FormatInt(infra.JobConfig.UserId, 10),
			},
			TtlSecondsAfterFinished: cfg.K8S.Job.TtlSecondsAfterFinished,
			ImageChecker:            cfg.K8S.Job.ImageChecker,
		})
		if err != nil {
			l.ErrorContext(ctx, "Error creating job checker", "error", err)
			return err
		}

		l.InfoContext(ctx, "Creating job", "jobName", jobName)
		err = k8sAPI.CreateJob(ctx, &api.JobInput{
			Namespace:               cfg.K8S.Namespace,
			JobName:                 jobName,
			Image:                   cfg.K8S.Job.Image,
			Cmd:                     cfg.K8S.Job.Command,
			Envs:                    cfg.K8S.Job.Envs,
			TtlSecondsAfterFinished: cfg.K8S.Job.TtlSecondsAfterFinished,
			ImageChecker:            cfg.K8S.Job.ImageChecker,
		})
		if err != nil {
			l.ErrorContext(ctx, "Error creating job", "error", err)
			return err
		}
	}

	return nil
}
