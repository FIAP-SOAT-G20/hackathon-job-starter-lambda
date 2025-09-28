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
	ServiceAccountName      string
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
	var backOffLimit = jobInput.BackOffLimit

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
	var ttlSecondsAfterFinished = int32(jobInput.TtlSecondsAfterFinished.Seconds())
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
							Name:            jobInput.JobName,
							Image:           jobInput.Image,
							ImagePullPolicy: v1.PullAlways,
							Env:             envVars,
						},
					},
					RestartPolicy:      v1.RestartPolicyNever,
					ImagePullSecrets:   imagePullSecrets,
					ServiceAccountName: jobInput.ServiceAccountName,
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

	watch, err := k.Client.BatchV1().
		Jobs(jobInput.Namespace).
		Watch(ctx, metav1.ListOptions{
			FieldSelector: "metadata.name=" + finalJobName,
		})
	if err != nil {
		log.Error().Err(err).Any("job", finalJobName).Any("namespace", jobInput.Namespace).Msg("Error watching job")
		return err
	}
	for event := range watch.ResultChan() {
		job := event.Object.(*batchv1.Job)
		if job.Status.Active > 0 {
			log.Info().Any("job", finalJobName).Any("namespace", jobInput.Namespace).Msg("Job started successfully")
			break
		}
		if job.Status.Failed > 0 {
			log.Info().
				Any("job", finalJobName).
				Any("namespace", jobInput.Namespace).
				Err(fmt.Errorf("job failed")).
				Int32("failedPods", job.Status.Failed).
				Msg("Job failed")

			pods, _ := k.Client.CoreV1().Pods(jobInput.Namespace).List(ctx, metav1.ListOptions{
				LabelSelector: "job-name=" + finalJobName,
			})
			for _, pod := range pods.Items {
				for _, cs := range pod.Status.ContainerStatuses {
					if cs.State.Terminated != nil {
						log.Error().
							Str("job", finalJobName).
							Str("namespace", jobInput.Namespace).
							Str("pod", pod.Name).
							Int32("exitCode", cs.State.Terminated.ExitCode).
							Str("reason", cs.State.Terminated.Reason).
							Str("message", cs.State.Terminated.Message).
							Msg("Job failed with container error")
					}
				}
			}
			return fmt.Errorf("job failed")
		}
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

func (k *K8sAPI) GetLastJobStatus(ctx context.Context, jobName, namespace string) (string, error) {
	jobs := k.Client.BatchV1().Jobs(namespace)
	job, err := jobs.Get(ctx, jobName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// Check if job has any conditions
	if len(job.Status.Conditions) == 0 {
		// If no conditions are set yet, check the job status directly
		if job.Status.Active > 0 {
			return "Running", nil
		}
		if job.Status.Succeeded > 0 {
			return "Complete", nil
		}
		if job.Status.Failed > 0 {
			return "Failed", nil
		}
		// Job is still pending
		return "Pending", nil
	}

	return string(job.Status.Conditions[len(job.Status.Conditions)-1].Type), nil
}
