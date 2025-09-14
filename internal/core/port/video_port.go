package port

import "github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/dto"

type VideoUsecase interface {
	UpdateVideoStatus(videoId int64, status dto.VideoProcessingStatus) error
}

type VideoGateway interface {
	UpdateVideoStatus(videoId int64, status dto.VideoProcessingStatus) error
}

type VideoController interface {
	UpdateVideoStatus(videoId int64, status dto.VideoProcessingStatus) error
}
