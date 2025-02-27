package initialization

import (
	"kava/internal/configuration"
	"testing"

	"github.com/stretchr/testify/assert"

	"os"
)

func TestCreateLogger(t *testing.T) {
	// Очистка после тестов
	defer func() {
		os.Remove("kava.log")
		os.Remove("test.log")
	}()

	t.Run("Create logger with nil config", func(t *testing.T) {
		logger, err := CreateLogger(nil)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})

	t.Run("Create logger with empty config", func(t *testing.T) {
		cfg := &configuration.LoggingConfig{}
		logger, err := CreateLogger(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})

	t.Run("Create logger with valid debug level", func(t *testing.T) {
		cfg := &configuration.LoggingConfig{
			Level: "debug",
		}
		logger, err := CreateLogger(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})

	t.Run("Create logger with valid info level", func(t *testing.T) {
		cfg := &configuration.LoggingConfig{
			Level: "info",
		}
		logger, err := CreateLogger(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})

	t.Run("Create logger with valid warn level", func(t *testing.T) {
		cfg := &configuration.LoggingConfig{
			Level: "warn",
		}
		logger, err := CreateLogger(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})

	t.Run("Create logger with valid error level", func(t *testing.T) {
		cfg := &configuration.LoggingConfig{
			Level: "error",
		}
		logger, err := CreateLogger(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})

	t.Run("Create logger with invalid level", func(t *testing.T) {
		cfg := &configuration.LoggingConfig{
			Level: "invalid",
		}
		logger, err := CreateLogger(cfg)
		assert.Error(t, err)
		assert.Nil(t, logger)
		assert.Equal(t, "logging level is incorrect", err.Error())
	})

	t.Run("Create logger with custom output path", func(t *testing.T) {
		cfg := &configuration.LoggingConfig{
			Level:  "info",
			Output: "test.log",
		}
		logger, err := CreateLogger(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, logger)

		// Проверяем, что файл лога создан
		_, err = os.Stat("test.log")
		assert.NoError(t, err)
	})

	t.Run("Check default values", func(t *testing.T) {
		logger, err := CreateLogger(nil)
		assert.NoError(t, err)
		assert.NotNil(t, logger)

		// Проверяем, что файл лога создан с дефолтным именем
		_, err = os.Stat("kava.log")
		assert.NoError(t, err)
	})
}
