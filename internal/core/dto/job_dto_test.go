package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartJobInput(t *testing.T) {
	t.Run("should create StartJobInput with valid data", func(t *testing.T) {
		// Arrange
		expectedName := "test-job"
		expectedVariables := map[string]string{
			"ENV_VAR_1": "value1",
			"ENV_VAR_2": "value2",
		}

		// Act
		input := StartJobInput{
			Name:      expectedName,
			Variables: expectedVariables,
		}

		// Assert
		assert.Equal(t, expectedName, input.Name)
		assert.Equal(t, expectedVariables, input.Variables)
		assert.Len(t, input.Variables, 2)
	})

	t.Run("should create StartJobInput with empty variables", func(t *testing.T) {
		// Arrange
		expectedName := "test-job"

		// Act
		input := StartJobInput{
			Name:      expectedName,
			Variables: make(map[string]string),
		}

		// Assert
		assert.Equal(t, expectedName, input.Name)
		assert.NotNil(t, input.Variables)
		assert.Len(t, input.Variables, 0)
	})

	t.Run("should create StartJobInput with nil variables", func(t *testing.T) {
		// Arrange
		expectedName := "test-job"

		// Act
		input := StartJobInput{
			Name:      expectedName,
			Variables: nil,
		}

		// Assert
		assert.Equal(t, expectedName, input.Name)
		assert.Nil(t, input.Variables)
	})
}
