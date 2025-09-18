package controller

import (
	"errors"
	"testing"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/dto"
	mocks "github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/port/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestVideoController_UpdateVideoStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVideoUsecase := mocks.NewMockVideoUsecase(ctrl)
	controller := NewVideoController(mockVideoUsecase)

	t.Run("should successfully update video status", func(t *testing.T) {
		// Arrange
		videoId := int64(123)
		status := dto.VideoStatusProcessing

		expectedInput := dto.UpdateVideoStatusInput{
			VideoId: videoId,
			Status:  status,
		}

		mockVideoUsecase.EXPECT().
			UpdateVideoStatus(gomock.Any(), expectedInput).
			Return(nil).
			Times(1)

		// Act
		err := controller.UpdateVideoStatus(videoId, status)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should return error when usecase fails", func(t *testing.T) {
		// Arrange
		videoId := int64(123)
		status := dto.VideoStatusFailed
		expectedError := errors.New("usecase error")

		expectedInput := dto.UpdateVideoStatusInput{
			VideoId: videoId,
			Status:  status,
		}

		mockVideoUsecase.EXPECT().
			UpdateVideoStatus(gomock.Any(), expectedInput).
			Return(expectedError).
			Times(1)

		// Act
		err := controller.UpdateVideoStatus(videoId, status)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should handle different video statuses", func(t *testing.T) {
		testCases := []struct {
			name   string
			status dto.VideoProcessingStatus
		}{
			{"uploaded", dto.VideoStatusUploaded},
			{"processing", dto.VideoStatusProcessing},
			{"reprocessing", dto.VideoStatusReprocessing},
			{"finished", dto.VideoStatusFinished},
			{"failed", dto.VideoStatusFailed},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Arrange
				videoId := int64(123)

				expectedInput := dto.UpdateVideoStatusInput{
					VideoId: videoId,
					Status:  tc.status,
				}

				mockVideoUsecase.EXPECT().
					UpdateVideoStatus(gomock.Any(), expectedInput).
					Return(nil).
					Times(1)

				// Act
				err := controller.UpdateVideoStatus(videoId, tc.status)

				// Assert
				assert.NoError(t, err)
			})
		}
	})

	t.Run("should handle zero video id", func(t *testing.T) {
		// Arrange
		videoId := int64(0)
		status := dto.VideoStatusProcessing

		expectedInput := dto.UpdateVideoStatusInput{
			VideoId: videoId,
			Status:  status,
		}

		mockVideoUsecase.EXPECT().
			UpdateVideoStatus(gomock.Any(), expectedInput).
			Return(nil).
			Times(1)

		// Act
		err := controller.UpdateVideoStatus(videoId, status)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should handle negative video id", func(t *testing.T) {
		// Arrange
		videoId := int64(-1)
		status := dto.VideoStatusProcessing

		expectedInput := dto.UpdateVideoStatusInput{
			VideoId: videoId,
			Status:  status,
		}

		mockVideoUsecase.EXPECT().
			UpdateVideoStatus(gomock.Any(), expectedInput).
			Return(nil).
			Times(1)

		// Act
		err := controller.UpdateVideoStatus(videoId, status)

		// Assert
		assert.NoError(t, err)
	})
}

func TestNewVideoController(t *testing.T) {
	t.Run("should create VideoController with usecase", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockVideoUsecase := mocks.NewMockVideoUsecase(ctrl)

		// Act
		controller := NewVideoController(mockVideoUsecase)

		// Assert
		assert.NotNil(t, controller)
		assert.Equal(t, mockVideoUsecase, controller.videoUsecase)
	})

	t.Run("should create VideoController with nil usecase", func(t *testing.T) {
		// Act
		controller := NewVideoController(nil)

		// Assert
		assert.NotNil(t, controller)
		assert.Nil(t, controller.videoUsecase)
	})
}
