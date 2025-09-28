package infrastructure

import (
	"context"
	"fmt"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/api"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/aws"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/aws/s3"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/aws/sns"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/k8s"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/logger"
	"k8s.io/client-go/kubernetes"
)

var l *logger.Logger

var k8sClient *kubernetes.Clientset
var k8sAPI *api.K8sAPI
var cfg *config.Config

type Infrastructure struct {
	Context          context.Context
	K8sAPI           *api.K8sAPI
	Config           *config.Config
	JobConfig        *config.JobConfig
	Logger           *logger.Logger
	SNS              *sns.SNS
	S3               *s3.S3
	AWSClientFactory *aws.ClientFactory
}

var infrastructure *Infrastructure

// init function is called during application startup. So, at this moment is initialized
// all structures and also the database connection
func init() {
	fmt.Println("🟠 Initing SQS consumer application")
	cfg = config.LoadLambdaConfig()
	jobConfig := config.LoadJobConfig()
	l = logger.NewLogger(cfg)
	ctx := context.Background()

	// init k8s client
	var err error
	k8sClient, err = k8s.ConnectToK8s(ctx, l, cfg)
	if err != nil {
		panic(err)
	}
	awsClientFactory, err := aws.NewClientFactory(ctx, cfg.AWS.Region)
	if err != nil {
		panic(err)
	}
	k8sAPI = api.NewK8sAPI(k8sClient)

	infrastructure = &Infrastructure{
		Context:          ctx,
		K8sAPI:           k8sAPI,
		Config:           cfg,
		JobConfig:        jobConfig,
		Logger:           l,
		SNS:              sns.NewSNS(cfg),
		S3:               s3.NewS3(cfg),
		AWSClientFactory: awsClientFactory,
	}

}

func GetInfrastructure() *Infrastructure {
	return infrastructure
}
