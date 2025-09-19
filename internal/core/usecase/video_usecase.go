package usecase

import (
	"context"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/port"
)

type VideoUsecase struct {
	videoGateway port.VideoGateway
}

func NewVideoUsecase(videoGateway port.VideoGateway) *VideoUsecase {
	return &VideoUsecase{videoGateway: videoGateway}
}

func (u *VideoUsecase) UpdateVideoStatus(ctx context.Context, input dto.UpdateVideoStatusInput) error {
	return u.videoGateway.UpdateVideoStatus(ctx, input)
}
