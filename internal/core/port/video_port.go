package port

import (
	"context"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/dto"
)

type VideoUsecase interface {
	UpdateVideoStatus(ctx context.Context, input dto.UpdateVideoStatusInput) error
}

type VideoGateway interface {
	UpdateVideoStatus(ctx context.Context, input dto.UpdateVideoStatusInput) error
}

type VideoController interface {
	UpdateVideoStatus(ctx context.Context, input dto.UpdateVideoStatusInput) error
}
