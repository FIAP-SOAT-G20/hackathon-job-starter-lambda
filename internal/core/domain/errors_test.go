package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError(t *testing.T) {
	t.Run("should create ValidationError with message and error", func(t *testing.T) {
		// Arrange
		expectedMessage := "validation failed"
		expectedErr := errors.New("field is required")

		// Act
		validationErr := &ValidationError{
			Message: expectedMessage,
			Err:     expectedErr,
		}

		// Assert
		assert.Equal(t, expectedMessage, validationErr.Message)
		assert.Equal(t, expectedErr, validationErr.Err)
		assert.Equal(t, expectedErr.Error(), validationErr.Error())
	})

	t.Run("should create ValidationError with only message", func(t *testing.T) {
		// Arrange
		expectedMessage := "validation failed"

		// Act
		validationErr := &ValidationError{
			Message: expectedMessage,
			Err:     nil,
		}

		// Assert
		assert.Equal(t, expectedMessage, validationErr.Message)
		assert.Nil(t, validationErr.Err)
		assert.Equal(t, expectedMessage, validationErr.Error())
	})

	t.Run("should create ValidationError using NewValidationError", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("field is required")

		// Act
		validationErr := NewValidationError(expectedErr)

		// Assert
		assert.Equal(t, ErrValidationError, validationErr.Message)
		assert.Equal(t, expectedErr, validationErr.Err)
		assert.Equal(t, expectedErr.Error(), validationErr.Error())
	})
}

func TestNotFoundError(t *testing.T) {
	t.Run("should create NotFoundError with message", func(t *testing.T) {
		// Arrange
		expectedMessage := "resource not found"

		// Act
		notFoundErr := &NotFoundError{
			Message: expectedMessage,
		}

		// Assert
		assert.Equal(t, expectedMessage, notFoundErr.Message)
		assert.Equal(t, expectedMessage, notFoundErr.Error())
	})

	t.Run("should create NotFoundError using NewNotFoundError", func(t *testing.T) {
		// Arrange
		expectedMessage := "resource not found"

		// Act
		notFoundErr := NewNotFoundError(expectedMessage)

		// Assert
		assert.Equal(t, expectedMessage, notFoundErr.Message)
		assert.Equal(t, expectedMessage, notFoundErr.Error())
	})
}

func TestInternalError(t *testing.T) {
	t.Run("should create InternalError with message and error", func(t *testing.T) {
		// Arrange
		expectedMessage := "internal server error"
		expectedErr := errors.New("database connection failed")

		// Act
		internalErr := &InternalError{
			Message: expectedMessage,
			Err:     expectedErr,
		}

		// Assert
		assert.Equal(t, expectedMessage, internalErr.Message)
		assert.Equal(t, expectedErr, internalErr.Err)
		assert.Equal(t, expectedErr.Error(), internalErr.Error())
	})

	t.Run("should create InternalError with only message", func(t *testing.T) {
		// Arrange
		expectedMessage := "internal server error"

		// Act
		internalErr := &InternalError{
			Message: expectedMessage,
			Err:     nil,
		}

		// Assert
		assert.Equal(t, expectedMessage, internalErr.Message)
		assert.Nil(t, internalErr.Err)
		assert.Equal(t, expectedMessage, internalErr.Error())
	})

	t.Run("should create InternalError using NewInternalError", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("database connection failed")

		// Act
		internalErr := NewInternalError(expectedErr)

		// Assert
		assert.Equal(t, ErrInternalError, internalErr.Message)
		assert.Equal(t, expectedErr, internalErr.Err)
		assert.Equal(t, expectedErr.Error(), internalErr.Error())
	})
}

func TestInvalidInputError(t *testing.T) {
	t.Run("should create InvalidInputError with message", func(t *testing.T) {
		// Arrange
		expectedMessage := "invalid input provided"

		// Act
		invalidInputErr := &InvalidInputError{
			Message: expectedMessage,
		}

		// Assert
		assert.Equal(t, expectedMessage, invalidInputErr.Message)
		assert.Equal(t, expectedMessage, invalidInputErr.Error())
	})

	t.Run("should create InvalidInputError using NewInvalidInputError", func(t *testing.T) {
		// Arrange
		expectedMessage := "invalid input provided"

		// Act
		invalidInputErr := NewInvalidInputError(expectedMessage)

		// Assert
		assert.Equal(t, expectedMessage, invalidInputErr.Message)
		assert.Equal(t, expectedMessage, invalidInputErr.Error())
	})
}

func TestErrorConstants(t *testing.T) {
	t.Run("should have correct error constants", func(t *testing.T) {
		assert.Equal(t, "data conflicts with existing data", ErrConflict)
		assert.Equal(t, "data not found", ErrNotFound)
		assert.Equal(t, "invalid parameter", ErrInvalidParam)
		assert.Equal(t, "invalid query parameters", ErrInvalidQueryParams)
		assert.Equal(t, "invalid body", ErrInvalidBody)
		assert.Equal(t, "invalid token duration format", ErrTokenDuration)
		assert.Equal(t, "error creating token", ErrTokenCreation)
		assert.Equal(t, "access token has expired", ErrExpiredToken)
		assert.Equal(t, "access token is invalid", ErrInvalidToken)
		assert.Equal(t, "invalid status transition", ErrOrderInvalidStatusTransition)
		assert.Equal(t, "order without products", ErrOrderWithoutProducts)
		assert.Equal(t, "product is mandatory", ErrProductIsMandatory)
		assert.Equal(t, "staff is mandatory", ErrStaffIdIsMandatory)
		assert.Equal(t, "order is mandatory", ErrOrderIsMandatory)
		assert.Equal(t, "order is not on status open", ErrOrderIsNotOpen)
		assert.Equal(t, "invalid role", ErrRoleInvalid)
		assert.Equal(t, "page must be greater than zero", ErrPageMustBeGreaterThanZero)
		assert.Equal(t, "limit must be between 1 and 100", ErrLimitMustBeBetween1And100)
		assert.Equal(t, "internal server error", ErrInternalError)
		assert.Equal(t, "unknown error", ErrUnknownError)
		assert.Equal(t, "validation error", ErrValidationError)
		assert.Equal(t, "invalid input", ErrInvalidInput)
		assert.Equal(t, "precondition failed", ErrPreconditionFailed)
		assert.Equal(t, "failed to create payment external", ErrFailedToCreatePaymentExternal)
		assert.Equal(t, "failed to fetch customer", ErrFetchingCustomer)
	})
}
