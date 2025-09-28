package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {

	// Environment
	Environment string

	// K8S Settings
	K8S struct {
		Namespace          string
		ContextName        string
		MasterUrl          string
		ServiceAccountName string
		Job                struct {
			Prefix                  string
			Image                   string
			Command                 string
			Envs                    map[string]string
			TtlSecondsAfterFinished time.Duration
			BackOffLimit            int32
			JobName                 string
			ImageChecker            string
		}
	}

	AWS struct {
		Region          string
		AccessKey       string
		SecretAccessKey string
		SessionToken    string
		SNS             struct {
			TopicArn string
		}
		SQS struct {
			QueueURL         string
			WorkerPoolSize   int
			MaxMessagesBatch int
			WaitTimeSeconds  int
		}
	}
}

type JobConfig struct {
	JobName   string
	Namespace string
	VideoId   int64
	UserId    int64
}

func LoadLambdaConfig() *Config {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: .env file not found or error loading it: %v", err)
	}

	// Environment
	environment := getEnv("ENVIRONMENT", "development")

	// K8S Settings
	k8sNamespace := getEnv("K8S_NAMESPACE", "default")
	k8sImage := getEnv("K8S_JOB_IMAGE", "ghcr.io/fiap-soat-g20/hackathon-job-starter-lambda:latest")
	k8sCommand := getEnv("K8S_JOB_COMMAND", "echo \"Hello, World\"")
	k8sJobPrefix := getEnv("K8S_JOB_PREFIX", "video-processor")
	k8sServiceAccountName := getEnv("K8S_SERVICE_ACCOUNT_NAME", "job-checker-sa")
	k8sJobEnvs := getEnvsWithPrefix("K8S_JOB_ENV_")
	k8sJobTtlSecondsAfterFinished, err := time.ParseDuration(getEnv("K8S_JOB_TTL_SECONDS_AFTER_FINISHED", "10s"))
	if err != nil {
		log.Printf("Warning: K8S_JOB_TTL_SECONDS_AFTER_FINISHED is not a valid duration: %v. Setting to 10s", err)
		k8sJobTtlSecondsAfterFinished = 10 * time.Second
	}
	k8sJobBackOffLimit, err := strconv.Atoi(getEnv("K8S_JOB_BACK_OFF_LIMIT", "3"))
	if err != nil {
		log.Printf("Warning: K8S_JOB_BACK_OFF_LIMIT is not a valid integer: %v. Setting to 3", err)
		k8sJobBackOffLimit = 3
	}
	k8sJobImageChecker := getEnv("K8S_JOB_IMAGE_CHECKER", "docker.io/library/job-checker:latest")

	awsRegion := getEnv("AWS_REGION", "us-east-1")
	awsAccessKey := getEnv("AWS_ACCESS_KEY_ID", "")
	awsSecretAccessKey := getEnv("AWS_SECRET_ACCESS_KEY", "")
	awsSnsTopicArn := getEnv("AWS_SNS_TOPIC_ARN", "")
	awsSqsQueueURL := getEnv("AWS_SQS_QUEUE_URL", "")
	awsSessionToken := getEnv("AWS_SESSION_TOKEN", "")

	// SQS Consumer settings
	sqsWorkerPoolSize := getIntEnv("SQS_WORKER_POOL_SIZE", 5)

	sqsMaxMessagesBatch := getIntEnv("SQS_MAX_MESSAGES_BATCH", 10)
	sqsWaitTimeSeconds := getIntEnv("SQS_WAIT_TIME_SECONDS", 20)
	config := &Config{}

	config.Environment = environment
	config.K8S.Namespace = k8sNamespace
	config.K8S.ServiceAccountName = k8sServiceAccountName
	config.K8S.Job.Image = k8sImage
	config.K8S.Job.Command = k8sCommand
	config.K8S.Job.Prefix = k8sJobPrefix
	config.K8S.Job.Envs = k8sJobEnvs
	config.K8S.Job.TtlSecondsAfterFinished = k8sJobTtlSecondsAfterFinished
	config.K8S.Job.BackOffLimit = int32(k8sJobBackOffLimit)
	config.K8S.Job.ImageChecker = k8sJobImageChecker
	config.AWS.Region = awsRegion
	config.AWS.AccessKey = awsAccessKey
	config.AWS.SecretAccessKey = awsSecretAccessKey
	config.AWS.SNS.TopicArn = awsSnsTopicArn
	config.AWS.SQS.QueueURL = awsSqsQueueURL
	config.AWS.SQS.WorkerPoolSize = sqsWorkerPoolSize
	config.AWS.SQS.MaxMessagesBatch = sqsMaxMessagesBatch
	config.AWS.SQS.WaitTimeSeconds = sqsWaitTimeSeconds
	config.AWS.SessionToken = awsSessionToken
	return config
}

func LoadJobConfig() *JobConfig {
	jobName := getEnv("JOB_NAME", "video-processor")
	namespace := getEnv("JOB_NAMESPACE", "default")
	videoId, err := strconv.ParseInt(getEnv("JOB_VIDEO_ID", "0"), 10, 64)
	if err != nil {
		log.Printf("Warning: JOB_VIDEO_ID is not a valid integer: %v. Setting to 0", err)
		videoId = 0
	}
	userId, err := strconv.ParseInt(getEnv("JOB_USER_ID", "0"), 10, 64)
	if err != nil {
		log.Printf("Warning: JOB_USER_ID is not a valid integer: %v. Setting to 0", err)
		userId = 0
	}
	return &JobConfig{
		JobName:   jobName,
		Namespace: namespace,
		VideoId:   videoId,
		UserId:    userId,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvsWithPrefix(prefix string) map[string]string {
	envs := make(map[string]string)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, prefix) {
			fullEnvName := strings.Split(env, "=")[0]
			envName := strings.TrimPrefix(fullEnvName, prefix)
			envs[envName] = getEnv(fullEnvName, "")
		}
	}
	return envs
}

func getIntEnv(key string, defaultValue int) int {
	value, err := strconv.Atoi(getEnv(key, strconv.Itoa(defaultValue)))
	if err != nil || value < 0 {
		log.Printf("Warning: %s is not a valid integer: %v. Setting to %d", key, err, defaultValue)
		return defaultValue
	}
	return value
}
