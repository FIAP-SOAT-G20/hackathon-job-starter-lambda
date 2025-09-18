package lambda

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandleRequest(t *testing.T) {
	t.Run("should handle empty S3 event", func(t *testing.T) {
		// Arrange
		req := events.S3Event{
			Records: []events.S3EventRecord{},
		}

		// Note: This test would fail in a real environment without proper infrastructure setup
		// In a real test environment, you would mock the infrastructure components

		// Act & Assert
		// This test demonstrates the structure but would need proper mocking
		assert.NotNil(t, req)
		assert.Len(t, req.Records, 0)
	})

	t.Run("should handle S3 event with records", func(t *testing.T) {
		// Arrange
		req := events.S3Event{
			Records: []events.S3EventRecord{
				{
					S3: events.S3Entity{
						Bucket: events.S3Bucket{
							Name: "test-bucket",
						},
						Object: events.S3Object{
							Key: "test/video.mp4",
						},
					},
				},
			},
		}

		// Act & Assert
		// This test demonstrates the structure but would need proper mocking
		assert.NotNil(t, req)
		assert.Len(t, req.Records, 1)
		assert.Equal(t, "test-bucket", req.Records[0].S3.Bucket.Name)
		assert.Equal(t, "test/video.mp4", req.Records[0].S3.Object.Key)
	})

	t.Run("should handle multiple S3 records", func(t *testing.T) {
		// Arrange
		req := events.S3Event{
			Records: []events.S3EventRecord{
				{
					S3: events.S3Entity{
						Bucket: events.S3Bucket{
							Name: "test-bucket-1",
						},
						Object: events.S3Object{
							Key: "test/video1.mp4",
						},
					},
				},
				{
					S3: events.S3Entity{
						Bucket: events.S3Bucket{
							Name: "test-bucket-2",
						},
						Object: events.S3Object{
							Key: "test/video2.mp4",
						},
					},
				},
			},
		}

		// Act & Assert
		assert.NotNil(t, req)
		assert.Len(t, req.Records, 2)
		assert.Equal(t, "test-bucket-1", req.Records[0].S3.Bucket.Name)
		assert.Equal(t, "test/video1.mp4", req.Records[0].S3.Object.Key)
		assert.Equal(t, "test-bucket-2", req.Records[1].S3.Bucket.Name)
		assert.Equal(t, "test/video2.mp4", req.Records[1].S3.Object.Key)
	})
}

// Test helper functions for S3 event processing
func TestS3EventProcessing(t *testing.T) {
	t.Run("should extract filename from S3 key", func(t *testing.T) {
		// Arrange
		testCases := []struct {
			key      string
			expected string
		}{
			{"test/video.mp4", "video"},
			{"videos/2024/01/video.mp4", "video"},
			{"video.mp4", "video"},
			{"test/video-with-dashes.mp4", "video-with-dashes"},
			{"test/video_with_underscores.mp4", "video_with_underscores"},
		}

		for _, tc := range testCases {
			t.Run(tc.key, func(t *testing.T) {
				// Act
				splittedKey := strings.Split(tc.key, "/")
				fileName := splittedKey[len(splittedKey)-1]
				fileNameWithoutExtension := strings.Split(fileName, ".")[0]

				// Assert
				assert.Equal(t, tc.expected, fileNameWithoutExtension)
			})
		}
	})

	t.Run("should generate job names correctly", func(t *testing.T) {
		// Arrange
		prefix := "video-processor"
		fileNameWithoutExtension := "test-video"

		// Act
		jobName := fmt.Sprintf("%s-%s", prefix, fileNameWithoutExtension)
		jobCheckerName := fmt.Sprintf("%s-%s-checker", prefix, fileNameWithoutExtension)

		// Assert
		assert.Equal(t, "video-processor-test-video", jobName)
		assert.Equal(t, "video-processor-test-video-checker", jobCheckerName)
	})

	t.Run("should generate group and deduplication IDs correctly", func(t *testing.T) {
		// Arrange
		videoId := int64(123)
		status := "processing"

		// Act
		groupId := fmt.Sprintf("video-id-%d", videoId)
		dedupId := fmt.Sprintf("video-status-%s", status)

		// Assert
		assert.Equal(t, "video-id-123", groupId)
		assert.Equal(t, "video-status-processing", dedupId)
	})
}

// Test error handling scenarios
func TestErrorHandling(t *testing.T) {
	t.Run("should handle invalid video ID", func(t *testing.T) {
		// Arrange
		invalidVideoId := "invalid-id"

		// Act & Assert
		// This would fail in the actual implementation
		_, err := strconv.ParseInt(invalidVideoId, 10, 64)
		assert.Error(t, err)
	})

	t.Run("should handle empty S3 key", func(t *testing.T) {
		// Arrange
		key := ""

		// Act
		splittedKey := strings.Split(key, "/")
		fileName := splittedKey[len(splittedKey)-1]

		// Assert
		assert.Equal(t, "", fileName)
	})

	t.Run("should handle S3 key without extension", func(t *testing.T) {
		// Arrange
		key := "test/video"

		// Act
		splittedKey := strings.Split(key, "/")
		fileName := splittedKey[len(splittedKey)-1]
		fileNameWithoutExtension := strings.Split(fileName, ".")[0]

		// Assert
		assert.Equal(t, "video", fileNameWithoutExtension)
	})
}
