package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/dto"
	mocks "github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/infrastructure/aws/sns/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestVideoGateway_UpdateVideoStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSNS := mocks.NewMockSNSInterface(ctrl)
	gateway := NewVideoGateway(mockSNS)

	t.Run("should successfully update video status", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		input := dto.UpdateVideoStatusInput{
			VideoId: 123,
			UserId:  456,
			Status:  dto.VideoStatusProcessing,
		}

		expectedGroupId := "video-id-123"
		expectedDedupId := "video-status-processing"
		expectedFilterKey := "status"
		expectedFilterValue := "processing"

		mockSNS.EXPECT().
			Publish(ctx, gomock.Any(), expectedGroupId, expectedDedupId, expectedFilterKey, expectedFilterValue).
			Return(nil).
			Times(1)

		// Act
		err := gateway.UpdateVideoStatus(ctx, input)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should return error when SNS publish fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		input := dto.UpdateVideoStatusInput{
			VideoId: 123,
			UserId:  456,
			Status:  dto.VideoStatusFailed,
		}
		expectedError := errors.New("SNS publish failed")

		mockSNS.EXPECT().
			Publish(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(expectedError).
			Times(1)

		// Act
		err := gateway.UpdateVideoStatus(ctx, input)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should handle different video statuses", func(t *testing.T) {
		testCases := []struct {
			name           string
			status         dto.VideoProcessingStatus
			expectedStatus string
		}{
			{"uploaded", dto.VideoStatusUploaded, "uploaded"},
			{"processing", dto.VideoStatusProcessing, "processing"},
			{"reprocessing", dto.VideoStatusReprocessing, "reprocessing"},
			{"finished", dto.VideoStatusFinished, "finished"},
			{"failed", dto.VideoStatusFailed, "failed"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Arrange
				ctx := context.Background()
				input := dto.UpdateVideoStatusInput{
					VideoId: 123,
					UserId:  456,
					Status:  tc.status,
				}

				expectedGroupId := "video-id-123"
				expectedDedupId := "video-status-" + tc.expectedStatus
				expectedFilterKey := "status"
				expectedFilterValue := tc.expectedStatus

				mockSNS.EXPECT().
					Publish(ctx, gomock.Any(), expectedGroupId, expectedDedupId, expectedFilterKey, expectedFilterValue).
					Return(nil).
					Times(1)

				// Act
				err := gateway.UpdateVideoStatus(ctx, input)

				// Assert
				assert.NoError(t, err)
			})
		}
	})

	t.Run("should handle zero values in input", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		input := dto.UpdateVideoStatusInput{
			VideoId: 0,
			UserId:  0,
			Status:  "",
		}

		expectedGroupId := "video-id-0"
		expectedDedupId := "video-status-"
		expectedFilterKey := "status"
		expectedFilterValue := ""

		mockSNS.EXPECT().
			Publish(ctx, gomock.Any(), expectedGroupId, expectedDedupId, expectedFilterKey, expectedFilterValue).
			Return(nil).
			Times(1)

		// Act
		err := gateway.UpdateVideoStatus(ctx, input)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should create correct JSON payload", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		input := dto.UpdateVideoStatusInput{
			VideoId: 789,
			UserId:  101112,
			Status:  dto.VideoStatusFinished,
		}

		// Capture the JSON payload passed to SNS
		var capturedMessage string
		mockSNS.EXPECT().
			Publish(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(ctx context.Context, message string, groupId string, dedupId string, filterKey string, filterValue string) {
				capturedMessage = message
			}).
			Return(nil).
			Times(1)

		// Act
		err := gateway.UpdateVideoStatus(ctx, input)

		// Assert
		assert.NoError(t, err)

		// Verify the JSON payload structure
		var payload dto.VideoStatusPayload
		err = json.Unmarshal([]byte(capturedMessage), &payload)
		assert.NoError(t, err)
		assert.Equal(t, input.VideoId, payload.VideoId)
		assert.Equal(t, input.UserId, payload.UserId)
		assert.Equal(t, string(input.Status), payload.Status)
		assert.NotZero(t, payload.OccurredAt)
	})
}

func TestNewVideoGateway(t *testing.T) {
	t.Run("should create VideoGateway with SNS", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSNS := mocks.NewMockSNSInterface(ctrl)

		// Act
		gateway := NewVideoGateway(mockSNS)

		// Assert
		assert.NotNil(t, gateway)
		// Note: We can't access the private field directly, but we can verify the gateway was created
	})

	t.Run("should create VideoGateway with nil SNS", func(t *testing.T) {
		// Act
		gateway := NewVideoGateway(nil)

		// Assert
		assert.NotNil(t, gateway)
		// Note: We can't access the private field directly, but we can verify the gateway was created
	})
}
