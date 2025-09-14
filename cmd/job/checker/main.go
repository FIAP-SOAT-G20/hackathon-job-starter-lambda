package main

import (
	"context"
	"os"
	"time"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure"
)

func main() {
	infra := infrastructure.GetInfrastructure()
	l := infra.Logger
	jobConfig := infra.JobConfig
	k8sAPI := infra.K8sAPI

	// Create a logger with MDC-like attributes
	mdcLogger := l.With(
		"jobName", jobConfig.JobName,
		"namespace", jobConfig.Namespace,
		"component", "job-checker",
		"version", "1.0.0",
	)

	var backoffLimit int = 0

	for {
		mdcLogger.Info("Checking jobs")
		time.Sleep(5 * time.Second)
		jobStatus, err := k8sAPI.GetJobStatus(context.Background(), jobConfig.JobName, jobConfig.Namespace)
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
		switch jobStatus {
		case "Complete":
			mdcLogger.Info("Job completed")
			os.Exit(0)
		case "Failed":
			mdcLogger.Info("Job failed")
			os.Exit(1)
		case "Pending":
			mdcLogger.Info("Job pending")
		}
	}
}
