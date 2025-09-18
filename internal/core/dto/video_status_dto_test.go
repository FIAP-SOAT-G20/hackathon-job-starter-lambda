package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVideoProcessingStatus(t *testing.T) {
	t.Run("should have correct status constants", func(t *testing.T) {
		assert.Equal(t, VideoProcessingStatus("uploaded"), VideoStatusUploaded)
		assert.Equal(t, VideoProcessingStatus("processing"), VideoStatusProcessing)
		assert.Equal(t, VideoProcessingStatus("reprocessing"), VideoStatusReprocessing)
		assert.Equal(t, VideoProcessingStatus("finished"), VideoStatusFinished)
		assert.Equal(t, VideoProcessingStatus("failed"), VideoStatusFailed)
	})

	t.Run("should convert status to string correctly", func(t *testing.T) {
		assert.Equal(t, "uploaded", string(VideoStatusUploaded))
		assert.Equal(t, "processing", string(VideoStatusProcessing))
		assert.Equal(t, "reprocessing", string(VideoStatusReprocessing))
		assert.Equal(t, "finished", string(VideoStatusFinished))
		assert.Equal(t, "failed", string(VideoStatusFailed))
	})
}

func TestUpdateVideoStatusInput(t *testing.T) {
	t.Run("should create UpdateVideoStatusInput with valid data", func(t *testing.T) {
		// Arrange
		expectedVideoId := int64(123)
		expectedUserId := int64(456)
		expectedStatus := VideoStatusProcessing

		// Act
		input := UpdateVideoStatusInput{
			VideoId: expectedVideoId,
			UserId:  expectedUserId,
			Status:  expectedStatus,
		}

		// Assert
		assert.Equal(t, expectedVideoId, input.VideoId)
		assert.Equal(t, expectedUserId, input.UserId)
		assert.Equal(t, expectedStatus, input.Status)
	})

	t.Run("should create UpdateVideoStatusInput with zero values", func(t *testing.T) {
		// Act
		input := UpdateVideoStatusInput{}

		// Assert
		assert.Equal(t, int64(0), input.VideoId)
		assert.Equal(t, int64(0), input.UserId)
		assert.Equal(t, VideoProcessingStatus(""), input.Status)
	})
}

func TestVideoStatusPayload(t *testing.T) {
	t.Run("should create VideoStatusPayload with valid data", func(t *testing.T) {
		// Arrange
		expectedVideoId := int64(123)
		expectedUserId := int64(456)
		expectedStatus := "processing"
		expectedOccurredAt := time.Now()

		// Act
		payload := VideoStatusPayload{
			VideoId:    expectedVideoId,
			UserId:     expectedUserId,
			Status:     expectedStatus,
			OccurredAt: expectedOccurredAt,
		}

		// Assert
		assert.Equal(t, expectedVideoId, payload.VideoId)
		assert.Equal(t, expectedUserId, payload.UserId)
		assert.Equal(t, expectedStatus, payload.Status)
		assert.Equal(t, expectedOccurredAt, payload.OccurredAt)
	})

	t.Run("should create VideoStatusPayload with zero values", func(t *testing.T) {
		// Act
		payload := VideoStatusPayload{}

		// Assert
		assert.Equal(t, int64(0), payload.VideoId)
		assert.Equal(t, int64(0), payload.UserId)
		assert.Equal(t, "", payload.Status)
		assert.Equal(t, time.Time{}, payload.OccurredAt)
	})
}
