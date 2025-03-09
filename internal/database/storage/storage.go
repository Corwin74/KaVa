package storage

import (
	"context"
	"errors"
	"kava/internal/common"
	"kava/internal/database/compute"
	"kava/internal/database/storage/wal"

	"go.uber.org/zap"
)

var ErrorNotExist = errors.New("key not exist")

// Storage - хранит данные используя engine
type Storage struct {
	engine    Engine
	wal       WAL
	stream    <-chan []wal.Log
	generator *IDGenerator
	logger    *zap.Logger
}

// NewStorage - конструктор
func NewStorage(engine Engine, wal WAL, logger *zap.Logger) (*Storage, error) {
	if engine == nil {
		return nil, errors.New("engine is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	storage := &Storage{
		engine: engine,
		logger: logger,
		wal:    wal,
	}

	var lastLSN int64
	if storage.wal != nil {
		logs, err := storage.wal.Recover()
		if err != nil {
			logger.Error("failed to recover data from WAL", zap.Error(err))
		} else {
			lastLSN = storage.applyData(logs)
		}
	}

	storage.generator = NewIDGenerator(lastLSN)
	return storage, nil
}

// Set - сохраняет данные используя engine
func (s *Storage) Set(ctx context.Context, key, value string) error {
	txID := s.generator.Generate()
	ctx = common.ContextWithTxID(ctx, txID)
	
	if s.wal != nil {
		futureResponse := s.wal.Set(ctx, key, value)
		if err := futureResponse.Get(); err != nil {
			return err
		}
	}

	s.engine.Set(ctx, key, value)
	return nil
}

// Get - получает данные, используя engine
func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	txID := s.generator.Generate()
	ctx = common.ContextWithTxID(ctx, txID)

	v, exist := s.engine.Get(ctx, key)
	if !exist {
		return "", ErrorNotExist
	}
	return v, nil
}

// Del - удаляет данные, используя движок
func (s *Storage) Del(ctx context.Context, key string) error {
	txID := s.generator.Generate()
	ctx = common.ContextWithTxID(ctx, txID)

	if s.wal != nil {
		futureResponse := s.wal.Del(ctx, key)
		if err := futureResponse.Get(); err != nil {
			return err
		}
	}

	s.engine.Del(ctx, key)
	return nil
}

func (s *Storage) applyData(logs []wal.Log) int64 {
	var lastLSN int64
	for _, log := range logs {
		lastLSN = max(lastLSN, log.LSN)
		ctx := common.ContextWithTxID(context.Background(), log.LSN)
		switch log.CommandID {
		case compute.SetCommandID:
			s.engine.Set(ctx, log.Arguments[0], log.Arguments[1])
		case compute.DelCommandID:
			s.engine.Del(ctx, log.Arguments[0])
		}
	}

	return lastLSN
}