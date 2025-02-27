package initialization

import (
	"errors"
	"kava/internal/configuration"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
	errorLevel = "error"
)

const (
	defaultEncoding   = "json"
	defaultLevel      = zapcore.InfoLevel
	defaultOutputPath = "kava.log"
)

// CreateLogger -- конструктор логгера
func CreateLogger(cfg *configuration.LoggingConfig) (*zap.Logger, error) {
	level := defaultLevel
	output := defaultOutputPath

	if cfg != nil {
		if cfg.Level != "" {
			supportedLoggingLevels := map[string]zapcore.Level{
				debugLevel: zapcore.DebugLevel,
				infoLevel:  zapcore.InfoLevel,
				warnLevel:  zapcore.WarnLevel,
				errorLevel: zapcore.ErrorLevel,
			}

			var exist bool
			if level, exist = supportedLoggingLevels[cfg.Level]; !exist {
				return nil, errors.New("logging level is incorrect")
			}
			
			if cfg.Output != "" {
				// TODO: need to create a
				// directory if it is missing
				output = cfg.Output
			}
		}
	}

	loggerCfg := zap.Config{
		Encoding:    defaultEncoding,
		Level:       zap.NewAtomicLevelAt(level),
		OutputPaths: []string{output},
	}

	return loggerCfg.Build()
}
