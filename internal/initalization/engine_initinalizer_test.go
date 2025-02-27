package initialization

import (
    "kava/internal/configuration"
    "testing"
    "go.uber.org/zap"
    "github.com/stretchr/testify/assert"
)

func TestCreateEngine(t *testing.T) {
    // Create a test logger
    logger, _ := zap.NewDevelopment()

    t.Run("Create engine with nil config", func(t *testing.T) {
        engine, err := CreateEngine(nil, logger)
        assert.NoError(t, err)
        assert.NotNil(t, engine)
    })

    t.Run("Create engine with empty config type", func(t *testing.T) {
        cfg := &configuration.EngineConfig{
            Type: "",
        }
        engine, err := CreateEngine(cfg, logger)
        assert.NoError(t, err)
        assert.NotNil(t, engine)
    })

    t.Run("Create engine with valid in_memory type", func(t *testing.T) {
        cfg := &configuration.EngineConfig{
            Type: "in_memory",
        }
        engine, err := CreateEngine(cfg, logger)
        assert.NoError(t, err)
        assert.NotNil(t, engine)
    })

    t.Run("Create engine with unsupported type", func(t *testing.T) {
        cfg := &configuration.EngineConfig{
            Type: "unsupported_type",
        }
        engine, err := CreateEngine(cfg, logger)
        assert.Error(t, err)
        assert.Nil(t, engine)
        assert.Equal(t, "engine type is incorrect", err.Error())
    })
}