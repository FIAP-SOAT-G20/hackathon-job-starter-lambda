package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/dto"
	mocks "github.com/FIAP-SOAT-G20/hackathon-job-starter-lambda/internal/core/port/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestVideoUsecase_UpdateVideoStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockVideoGateway := mocks.NewMockVideoGateway(ctrl)
	usecase := NewVideoUsecase(mockVideoGateway)

	t.Run("should successfully update video status", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		input := dto.UpdateVideoStatusInput{
			VideoId: 123,
			UserId:  456,
			Status:  dto.VideoStatusProcessing,
		}

		mockVideoGateway.EXPECT().
			UpdateVideoStatus(ctx, input).
			Return(nil).
			Times(1)

		// Act
		err := usecase.UpdateVideoStatus(ctx, input)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should return error when gateway fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		input := dto.UpdateVideoStatusInput{
			VideoId: 123,
			UserId:  456,
			Status:  dto.VideoStatusFailed,
		}
		expectedError := errors.New("gateway error")

		mockVideoGateway.EXPECT().
			UpdateVideoStatus(ctx, input).
			Return(expectedError).
			Times(1)

		// Act
		err := usecase.UpdateVideoStatus(ctx, input)

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
				ctx := context.Background()
				input := dto.UpdateVideoStatusInput{
					VideoId: 123,
					UserId:  456,
					Status:  tc.status,
				}

				mockVideoGateway.EXPECT().
					UpdateVideoStatus(ctx, input).
					Return(nil).
					Times(1)

				// Act
				err := usecase.UpdateVideoStatus(ctx, input)

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

		mockVideoGateway.EXPECT().
			UpdateVideoStatus(ctx, input).
			Return(nil).
			Times(1)

		// Act
		err := usecase.UpdateVideoStatus(ctx, input)

		// Assert
		assert.NoError(t, err)
	})
}

func TestNewVideoUsecase(t *testing.T) {
	t.Run("should create VideoUsecase with gateway", func(t *testing.T) {
		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockVideoGateway := mocks.NewMockVideoGateway(ctrl)

		// Act
		usecase := NewVideoUsecase(mockVideoGateway)

		// Assert
		assert.NotNil(t, usecase)
		assert.Equal(t, mockVideoGateway, usecase.videoGateway)
	})

	t.Run("should create VideoUsecase with nil gateway", func(t *testing.T) {
		// Act
		usecase := NewVideoUsecase(nil)

		// Assert
		assert.NotNil(t, usecase)
		assert.Nil(t, usecase.videoGateway)
	})
}
