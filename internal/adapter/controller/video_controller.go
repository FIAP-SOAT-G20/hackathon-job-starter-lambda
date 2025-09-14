package controller

import (
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/dto"
	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/port"
)

type VideoController struct {
	videoUsecase port.VideoUsecase
}

func NewVideoController(videoUsecase port.VideoUsecase) *VideoController {
	return &VideoController{videoUsecase: videoUsecase}
}

func (c *VideoController) UpdateVideoStatus(videoId int64, status dto.VideoProcessingStatus) error {
	return c.videoUsecase.UpdateVideoStatus(videoId, status)
}
