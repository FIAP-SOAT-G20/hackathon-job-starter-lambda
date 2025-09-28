package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVideoProcessingStatus(t *testing.T) {
	t.Run("should have correct status constants", func(t *testing.T) {
		assert.Equal(t, VideoProcessingStatus("UPLOADED"), VideoStatusUploaded)
		assert.Equal(t, VideoProcessingStatus("PROCESSING"), VideoStatusProcessing)
		assert.Equal(t, VideoProcessingStatus("REPROCESSING"), VideoStatusReprocessing)
		assert.Equal(t, VideoProcessingStatus("FINISHED"), VideoStatusFinished)
		assert.Equal(t, VideoProcessingStatus("FAILED"), VideoStatusFailed)
	})

	t.Run("should convert status to string correctly", func(t *testing.T) {
		assert.Equal(t, "UPLOADED", string(VideoStatusUploaded))
		assert.Equal(t, "PROCESSING", string(VideoStatusProcessing))
		assert.Equal(t, "REPROCESSING", string(VideoStatusReprocessing))
		assert.Equal(t, "FINISHED", string(VideoStatusFinished))
		assert.Equal(t, "FAILED", string(VideoStatusFailed))
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
		expectedStatus := "PROCESSING"
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
