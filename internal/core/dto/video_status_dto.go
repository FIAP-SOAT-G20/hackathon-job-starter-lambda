package dto

import "time"

type VideoProcessingStatus string

const (
	VideoStatusUploaded     VideoProcessingStatus = "uploaded"
	VideoStatusProcessing   VideoProcessingStatus = "processing"
	VideoStatusReprocessing VideoProcessingStatus = "reprocessing"
	VideoStatusFinished     VideoProcessingStatus = "finished"
	VideoStatusFailed       VideoProcessingStatus = "failed"
)

type VideoStatusDTO struct {
	VideoId    int64     `json:"video_id"`
	Status     string    `json:"status"`
	OccurredAt time.Time `json:"occurred_at"`
}
