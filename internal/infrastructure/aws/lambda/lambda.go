package lambda

import (
	"context"
	"fmt"
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
		splittedKey := strings.Split(record.S3.Object.Key, "/")
		fileName := splittedKey[len(splittedKey)-1]
		fileNameWithoutExtension := strings.Split(fileName, ".")[0]
		jobName := fmt.Sprintf("%s-%s", cfg.K8S.Job.Prefix, fileNameWithoutExtension)
		l.InfoContext(ctx, "Creating job", "jobName", jobName)
		k8sAPI.CreateJob(ctx, &api.JobInput{
			Namespace:               cfg.K8S.Namespace,
			JobName:                 jobName,
			Image:                   cfg.K8S.Job.Image,
			Cmd:                     cfg.K8S.Job.Command,
			Envs:                    cfg.K8S.Job.Envs,
			TtlSecondsAfterFinished: cfg.K8S.Job.TtlSecondsAfterFinished,
			ImageChecker:            cfg.K8S.Job.ImageChecker,
		})
	}

	return nil
}
