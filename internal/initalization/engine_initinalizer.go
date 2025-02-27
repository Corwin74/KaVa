package initialization

import (
	"errors"
	"kava/internal/configuration"
	"kava/internal/database/storage/engine/in_memory"
	"go.uber.org/zap"
)

// CreateEngine -- создание движка базы данных
func CreateEngine(cfg *configuration.EngineConfig, logger *zap.Logger) (*in_memory.Engine, error) {
	

	if cfg == nil {
		return in_memory.NewEngine(logger)
	}

	if cfg.Type != "" {
		supportedTypes := map[string]struct{}{
			"in_memory": {},
		}

		if _, found := supportedTypes[cfg.Type]; !found {
			return nil, errors.New("engine type is incorrect")
		}
	}

	return in_memory.NewEngine(logger)

}
