package usecase

import (
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/port"
)

type VideoUsecase struct {
	videoGateway port.VideoGateway
}

func NewVideoUsecase(videoGateway port.VideoGateway) *VideoUsecase {
	return &VideoUsecase{videoGateway: videoGateway}
}

func (u *VideoUsecase) UpdateVideoStatus(videoId int64, status dto.VideoProcessingStatus) error {
	return u.videoGateway.UpdateVideoStatus(videoId, status)
}
