package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type JobInput struct {
	Namespace               string
	JobName                 string
	Image                   string
	Cmd                     string
	TtlSecondsAfterFinished time.Duration
	Envs                    map[string]string
	BackOffLimit            int32
	ImageChecker            string
}

type K8sAPI struct {
	Client *kubernetes.Clientset
}

func NewK8sAPI(client *kubernetes.Clientset) *K8sAPI {
	return &K8sAPI{Client: client}
}

func (k *K8sAPI) CreateJob(ctx context.Context, jobInput *JobInput) error {
	err := validateParams(jobInput.Namespace, jobInput.JobName, jobInput.Image, jobInput.Cmd)
	if err != nil {
		return err
	}

	finalJobName := jobInput.JobName
	jobs := k.Client.BatchV1().Jobs(jobInput.Namespace)
	var backOffLimit int32 = jobInput.BackOffLimit

	envVars := make([]v1.EnvVar, 0)
	if jobInput.Envs != nil {
		log.Info().Msg(fmt.Sprintf("Job %s envs:", jobInput.JobName))
		for key, value := range jobInput.Envs {
			log.Info().Msg(fmt.Sprintf("%s: %s", key, value))
			envVars = append(envVars, v1.EnvVar{Name: key, Value: value})
		}
	} else {
		log.Info().Msg(fmt.Sprintf("No environment variables was set for job %s ", jobInput.JobName))
	}

	imagePullSecrets := make([]v1.LocalObjectReference, 0)
	var ttlSecondsAfterFinished int32 = int32(jobInput.TtlSecondsAfterFinished.Seconds())
	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      finalJobName,
			Namespace: jobInput.Namespace,
		},
		Spec: batchv1.JobSpec{
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            "job-checker-" + jobInput.JobName,
							Image:           jobInput.ImageChecker,
							ImagePullPolicy: v1.PullIfNotPresent,
							Env: []v1.EnvVar{
								{
									Name:  "JOB_NAME",
									Value: jobInput.JobName,
								},
								{
									Name:  "NAMESPACE",
									Value: jobInput.Namespace,
								},
							},
						},
						{
							Name:            jobInput.JobName,
							Image:           jobInput.Image,
							ImagePullPolicy: v1.PullIfNotPresent,
							Env:             envVars,
						},
					},
					RestartPolicy:    v1.RestartPolicyOnFailure,
					ImagePullSecrets: imagePullSecrets,
				},
			},
			BackoffLimit: &backOffLimit,
		},
	}

	_, err = jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		log.Error().Err(err).Any("job", finalJobName).Any("namespace", jobInput.Namespace).Msg("Error creating job")
		return err
	}

	log.Info().Any("job", finalJobName).Any("namespace", jobInput.Namespace).Msg("Job created successfully")
	return nil
}

func validateParams(namespace, jobName, image, cmd string) error {
	if namespace == "" || jobName == "" || image == "" || cmd == "" {
		return errors.New("the following envs are mandatory: K8S_NAMESPACE, K8S_JOB_NAME, K8S_JOB_IMAGE, K8S_JOB_COMMAND")
	}
	return nil
}

func (k *K8sAPI) GetJobStatus(ctx context.Context, jobName, namespace string) (string, error) {
	jobs := k.Client.BatchV1().Jobs(namespace)
	job, err := jobs.Get(ctx, jobName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(job.Status.Conditions[0].Type), nil
}
