package gateway

import (
	"context"
	"encoding/json"
	"time"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/port"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/aws/sns"
)

type VideoGateway struct {
	sns sns.SNSInterface
}

func NewVideoGateway(sns sns.SNSInterface) port.VideoGateway {
	return &VideoGateway{sns: sns}
}

func (g *VideoGateway) UpdateVideoStatus(ctx context.Context, input dto.UpdateVideoStatusInput) error {

	json, err := json.Marshal(dto.VideoStatusPayload{
		VideoId:    input.VideoId,
		UserId:     input.UserId,
		Status:     string(input.Status),
		OccurredAt: time.Now(),
	})
	if err != nil {
		return err
	}

	return g.sns.Publish(ctx, string(json))
}
