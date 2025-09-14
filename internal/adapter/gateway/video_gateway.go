package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/aws/sns"
)

type VideoGateway struct {
	sns sns.SNS
}

func NewVideoGateway(sns sns.SNS) *VideoGateway {
	return &VideoGateway{sns: sns}
}

func (g *VideoGateway) UpdateVideoStatus(ctx context.Context, videoId int64, status dto.VideoProcessingStatus) error {

	json, err := json.Marshal(dto.VideoStatusDTO{
		VideoId:    videoId,
		Status:     string(status),
		OccurredAt: time.Now(),
	})
	if err != nil {
		return err
	}
	groupId := fmt.Sprintf("video-id-%d", videoId)
	dedupId := fmt.Sprintf("video-status-%s", string(status))
	filterKey := "video-status"
	filterValue := string(status)

	return g.sns.Publish(ctx, string(json), groupId, dedupId, filterKey, filterValue)
}
