package sns

import (
	"context"
	"log"

	myConfig "github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type SNS struct {
	Client   *sns.Client
	TopicArn string
}

func NewSNS(cfg *myConfig.Config) *SNS {
	snsClient := sns.NewFromConfig(aws.Config{Region: cfg.AWS.Region, Credentials: credentials.NewStaticCredentialsProvider(cfg.AWS.AccessKey, cfg.AWS.SecretAccessKey, cfg.AWS.SessionToken)})
	return &SNS{
		Client:   snsClient,
		TopicArn: cfg.AWS.SNS.TopicArn,
	}
}

func (s *SNS) Publish(ctx context.Context, message string, groupId string, dedupId string, filterKey string, filterValue string) error {
	publishInput := sns.PublishInput{TopicArn: aws.String(s.TopicArn), Message: aws.String(message)}
	if groupId != "" {
		publishInput.MessageGroupId = aws.String(groupId)
	}
	if dedupId != "" {
		publishInput.MessageDeduplicationId = aws.String(dedupId)
	}
	if filterKey != "" && filterValue != "" {
		publishInput.MessageAttributes = map[string]types.MessageAttributeValue{
			filterKey: {DataType: aws.String("String"), StringValue: aws.String(filterValue)},
		}
	}
	_, err := s.Client.Publish(ctx, &publishInput)
	if err != nil {
		log.Printf("Couldn't publish message to topic %v. Here's why: %v", s.TopicArn, err)
	}
	return err
}
