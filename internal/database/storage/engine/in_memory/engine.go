package in_memory

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

// NewEngine - конструктор движка
func NewEngine(logger *zap.Logger) (*Engine, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}
	mb := make(map[string]string)
	engine := &Engine{
		data:   mb,
		logger: logger,
	}

	return engine, nil
}

// Engine - хранит данные в памяти используя map
type Engine struct {
	data   map[string]string
	logger *zap.Logger
}

// Set - сохраняет значение по ключу
func (e *Engine) Set(ctx context.Context, key, value string) {
	e.data[key] = value
	e.logger.Debug("successfull set query")
}

// Get - возвращает значение по ключу
func (e *Engine) Get(ctx context.Context, key string) (string, bool) {
	v, exist := e.data[key]
	if exist {
		e.logger.Debug("successfull get query")
	} else {
		e.logger.Debug("key not found")
	}

	return v, exist
}

// Del - удаляет значение по ключу
func (e *Engine) Del(ctx context.Context, key string) {
	delete(e.data, key)
	e.logger.Debug("successfull del query")
}
