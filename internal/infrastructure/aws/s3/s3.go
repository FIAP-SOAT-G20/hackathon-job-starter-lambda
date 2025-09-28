package s3

import (
	"context"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	Client *s3.Client
}

func NewS3(cfg *config.Config) *S3 {
	return &S3{
		Client: s3.NewFromConfig(aws.Config{Region: cfg.AWS.Region, Credentials: credentials.NewStaticCredentialsProvider(cfg.AWS.AccessKey, cfg.AWS.SecretAccessKey, cfg.AWS.SessionToken)}),
	}
}

func (s *S3) GetObjectMetadata(ctx context.Context, bucket string, key string) (map[string]string, error) {
	object, err := s.Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return object.Metadata, nil
}
