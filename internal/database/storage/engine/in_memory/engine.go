package in_memory

import (
	"context"
	"errors"
	"sync"

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
	mu     sync.RWMutex
	data   map[string]string
	logger *zap.Logger
}

// Set - сохраняет значение по ключу
func (e *Engine) Set(ctx context.Context, key, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.data[key] = value
	e.logger.Debug(
		"successfull set query",
		zap.String("msg", "SET"),
		zap.String("key", key),
		zap.String("value", value))
}

// Get - возвращает значение по ключу
func (e *Engine) Get(ctx context.Context, key string) (string, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	v, exist := e.data[key]
	if exist {
		e.logger.Debug(
			"successfull get query",
			zap.String("msg", "GET"),
			zap.String("key", key),
		)
	} else {
		e.logger.Debug(
			"key not found",
			zap.String("msg", "key not found"),
			zap.String("key", key),
		)
	}

	return v, exist
}

// Del - удаляет значение по ключу
func (e *Engine) Del(ctx context.Context, key string) {
	e.mu.Lock()
	delete(e.data, key)
	e.mu.Unlock()
	e.logger.Debug(
		"successfull del query",
		zap.String("msg", "DEL"),
		zap.String("key", key),
	)
}
