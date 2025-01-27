package storage

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

var ErrorNotExist = errors.New("key not exist")

// Engine - интерфейс движка который умеет сохранять, запрашивать и удалять данные
type Engine interface {
	Set(context.Context, string, string)
	Get(context.Context, string) (string, bool)
	Del(context.Context, string)
}

// Storage - хранит данные используя engine
type Storage struct {
	engine Engine
	logger *zap.Logger
}


// NewStorage - конструктор
func NewStorage(engine Engine, logger *zap.Logger) (*Storage, error) {
	if engine == nil {
		return nil, errors.New("engine is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	storage := &Storage{
		engine: engine,
		logger: logger,
	}

	return storage, nil
}

// Set - сохраняет данные используя engine
func (s *Storage) Set(ctx context.Context, key, value string) error {
	s.engine.Set(ctx, key, value)
	return nil
}


// Get - получает данные, используя engine
func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	v, exist := s.engine.Get(ctx, key)
	if !exist {
		return "", ErrorNotExist
	}
	return v, nil
}


// Del - удаляет данные, используя движок
func (s *Storage) Del(ctx context.Context, key string) error {
	s.engine.Del(ctx, key)
	return nil
}
