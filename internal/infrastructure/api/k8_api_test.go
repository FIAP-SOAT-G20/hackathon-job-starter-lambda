package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes"
)

func TestValidateParams(t *testing.T) {
	t.Run("should return nil for valid parameters", func(t *testing.T) {
		// Arrange
		namespace := "test-namespace"
		jobName := "test-job"
		image := "test-image:latest"
		cmd := "test-command"

		// Act
		err := validateParams(namespace, jobName, image, cmd)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("should return error for empty namespace", func(t *testing.T) {
		// Arrange
		namespace := ""
		jobName := "test-job"
		image := "test-image:latest"
		cmd := "test-command"

		// Act
		err := validateParams(namespace, jobName, image, cmd)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mandatory")
	})

	t.Run("should return error for empty job name", func(t *testing.T) {
		// Arrange
		namespace := "test-namespace"
		jobName := ""
		image := "test-image:latest"
		cmd := "test-command"

		// Act
		err := validateParams(namespace, jobName, image, cmd)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mandatory")
	})

	t.Run("should return error for empty image", func(t *testing.T) {
		// Arrange
		namespace := "test-namespace"
		jobName := "test-job"
		image := ""
		cmd := "test-command"

		// Act
		err := validateParams(namespace, jobName, image, cmd)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mandatory")
	})

	t.Run("should return error for empty command", func(t *testing.T) {
		// Arrange
		namespace := "test-namespace"
		jobName := "test-job"
		image := "test-image:latest"
		cmd := ""

		// Act
		err := validateParams(namespace, jobName, image, cmd)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mandatory")
	})

	t.Run("should return error for all empty parameters", func(t *testing.T) {
		// Arrange
		namespace := ""
		jobName := ""
		image := ""
		cmd := ""

		// Act
		err := validateParams(namespace, jobName, image, cmd)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mandatory")
	})
}

func TestJobInput(t *testing.T) {
	t.Run("should create JobInput with valid data", func(t *testing.T) {
		// Arrange
		expectedNamespace := "test-namespace"
		expectedJobName := "test-job"
		expectedImage := "test-image:latest"
		expectedCmd := "test-command"
		expectedTtl := 30 * time.Second
		expectedEnvs := map[string]string{
			"ENV_VAR_1": "value1",
			"ENV_VAR_2": "value2",
		}
		expectedBackOffLimit := int32(3)
		expectedImageChecker := "checker-image:latest"
		expectedServiceAccountName := "test-service-account"

		// Act
		jobInput := &JobInput{
			Namespace:               expectedNamespace,
			JobName:                 expectedJobName,
			Image:                   expectedImage,
			Cmd:                     expectedCmd,
			TtlSecondsAfterFinished: expectedTtl,
			Envs:                    expectedEnvs,
			BackOffLimit:            expectedBackOffLimit,
			ImageChecker:            expectedImageChecker,
			ServiceAccountName:      expectedServiceAccountName,
		}

		// Assert
		assert.Equal(t, expectedNamespace, jobInput.Namespace)
		assert.Equal(t, expectedJobName, jobInput.JobName)
		assert.Equal(t, expectedImage, jobInput.Image)
		assert.Equal(t, expectedCmd, jobInput.Cmd)
		assert.Equal(t, expectedTtl, jobInput.TtlSecondsAfterFinished)
		assert.Equal(t, expectedEnvs, jobInput.Envs)
		assert.Equal(t, expectedBackOffLimit, jobInput.BackOffLimit)
		assert.Equal(t, expectedImageChecker, jobInput.ImageChecker)
		assert.Equal(t, expectedServiceAccountName, jobInput.ServiceAccountName)
	})

	t.Run("should create JobInput with zero values", func(t *testing.T) {
		// Act
		jobInput := &JobInput{}

		// Assert
		assert.Equal(t, "", jobInput.Namespace)
		assert.Equal(t, "", jobInput.JobName)
		assert.Equal(t, "", jobInput.Image)
		assert.Equal(t, "", jobInput.Cmd)
		assert.Equal(t, time.Duration(0), jobInput.TtlSecondsAfterFinished)
		assert.Nil(t, jobInput.Envs)
		assert.Equal(t, int32(0), jobInput.BackOffLimit)
		assert.Equal(t, "", jobInput.ImageChecker)
		assert.Equal(t, "", jobInput.ServiceAccountName)
	})

	t.Run("should create JobInput with nil environment variables", func(t *testing.T) {
		// Act
		jobInput := &JobInput{
			Namespace: "test-namespace",
			JobName:   "test-job",
			Image:     "test-image:latest",
			Cmd:       "test-command",
			Envs:      nil,
		}

		// Assert
		assert.Nil(t, jobInput.Envs)
	})

	t.Run("should create JobInput with empty environment variables", func(t *testing.T) {
		// Act
		jobInput := &JobInput{
			Namespace: "test-namespace",
			JobName:   "test-job",
			Image:     "test-image:latest",
			Cmd:       "test-command",
			Envs:      make(map[string]string),
		}

		// Assert
		assert.NotNil(t, jobInput.Envs)
		assert.Len(t, jobInput.Envs, 0)
	})
}

func TestNewK8sAPI(t *testing.T) {
	t.Run("should create K8sAPI with client", func(t *testing.T) {
		// Arrange
		client := &kubernetes.Clientset{}

		// Act
		k8sAPI := NewK8sAPI(client)

		// Assert
		assert.NotNil(t, k8sAPI)
		assert.Equal(t, client, k8sAPI.Client)
	})

	t.Run("should create K8sAPI with nil client", func(t *testing.T) {
		// Act
		k8sAPI := NewK8sAPI(nil)

		// Assert
		assert.NotNil(t, k8sAPI)
		assert.Nil(t, k8sAPI.Client)
	})
}
