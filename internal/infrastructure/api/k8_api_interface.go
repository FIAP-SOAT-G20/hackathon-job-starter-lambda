package api

import "context"

// K8sAPIInterface defines the contract for Kubernetes API operations
type K8sAPIInterface interface {
	CreateJob(ctx context.Context, jobInput *JobInput) error
	GetLastJobStatus(ctx context.Context, jobName, namespace string) (string, error)
}
