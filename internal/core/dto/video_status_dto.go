package dto

import "time"

type VideoProcessingStatus string

const (
	VideoStatusUploaded     VideoProcessingStatus = "UPLOADED"
	VideoStatusProcessing   VideoProcessingStatus = "PROCESSING"
	VideoStatusReprocessing VideoProcessingStatus = "REPROCESSING"
	VideoStatusFinished     VideoProcessingStatus = "FINISHED"
	VideoStatusFailed       VideoProcessingStatus = "FAILED"
)

type UpdateVideoStatusInput struct {
	VideoId int64                 `json:"video_id"`
	UserId  int64                 `json:"user_id"`
	Status  VideoProcessingStatus `json:"status"`
}

type VideoStatusPayload struct {
	VideoId    int64     `json:"video_id"`
	UserId     int64     `json:"user_id"`
	Status     string    `json:"status"`
	OccurredAt time.Time `json:"occurred_at"`
}
