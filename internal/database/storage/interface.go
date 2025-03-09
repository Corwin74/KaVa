package storage

import (
	"context"
	"kava/internal/database/storage/wal"
	"kava/pkg/concurrency"
)

//go:generate mockgen -destination=storage_mock.go -package=storage . Engine,WAL

// Engine - интерфейс движка который умеет сохранять, запрашивать и удалять данные
type Engine interface {
	Set(context.Context, string, string)
	Get(context.Context, string) (string, bool)
	Del(context.Context, string)
}

// WAL - интерфейс для WAL
type WAL interface {
	Recover() ([]wal.Log, error)
	Set(context.Context, string, string) concurrency.FutureError
	Del(context.Context, string) concurrency.FutureError
}
