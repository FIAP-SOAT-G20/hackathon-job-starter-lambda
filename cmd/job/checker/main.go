package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/adapter/gateway"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/usecase"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/config"
)

func main() {
	infra := infrastructure.GetInfrastructure()
	l := infra.Logger
	jobConfig := infra.JobConfig
	k8sAPI := infra.K8sAPI
	videoUsecase := usecase.NewVideoUsecase(gateway.NewVideoGateway(*infra.SNS))
	ctx := context.Background()

	mdcLogger := l.With(
		"jobName", jobConfig.JobName,
		"namespace", jobConfig.Namespace,
		"component", "job-checker",
		"version", "1.0.0",
	)

	updateVideoStatus(ctx, mdcLogger, videoUsecase, jobConfig, dto.VideoStatusUploaded)

	var backoffLimit int = 0
	var jobPending = false

	for {
		time.Sleep(1 * time.Second)
		jobStatus, err := k8sAPI.GetLastJobStatus(context.Background(), jobConfig.JobName, jobConfig.Namespace)
		mdcLogger.Info("Job status", "jobStatus", jobStatus)
		if err != nil {
			mdcLogger.Error("Error getting job status",
				"error", err,
				"backoffLimit", backoffLimit,
			)
			backoffLimit++
			if backoffLimit > 3 {
				mdcLogger.Error("Job failed", "backoffLimit", backoffLimit)
				os.Exit(1)
			}
		}
		mdcLogger.Info(fmt.Sprintf("Job %s", strings.ToLower(jobStatus)))
		switch jobStatus {
		case "Complete":
			updateVideoStatus(ctx, mdcLogger, videoUsecase, jobConfig, dto.VideoStatusFinished)
			os.Exit(0)
		case "Failed":
			updateVideoStatus(ctx, mdcLogger, videoUsecase, jobConfig, dto.VideoStatusFailed)
			os.Exit(1)
		case "Pending":
		case "Running":
			if !jobPending {
				updateVideoStatus(ctx, mdcLogger, videoUsecase, jobConfig, dto.VideoStatusProcessing)
				jobPending = true
			}
		}
	}
}

func updateVideoStatus(ctx context.Context, mdcLogger *slog.Logger, videoUsecase *usecase.VideoUsecase, jobConfig *config.JobConfig, status dto.VideoProcessingStatus) {
	err := videoUsecase.UpdateVideoStatus(ctx, dto.UpdateVideoStatusInput{
		VideoId: jobConfig.VideoId,
		UserId:  jobConfig.UserId,
		Status:  status,
	})
	if err != nil {
		mdcLogger.Error("Error updating video status", "error", err)
		os.Exit(1)
	}
}
