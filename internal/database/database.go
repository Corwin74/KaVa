package database

import (
	"context"
	"errors"
	"fmt"
	"kava/internal/database/compute"
	"kava/internal/database/storage"

	"go.uber.org/zap"
)

type computeLayer interface {
	Parse(string) (compute.Query, error)
}

type storageLayer interface {
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, error)
	Del(context.Context, string) error
}

// Database -- состав по слоям
type Database struct {
	computeLayer computeLayer
	storageLayer storageLayer
	logger       *zap.Logger
}

// NewDatabase -- конструктор Database
func NewDatabase(computeLayer computeLayer, storageLayer storageLayer, logger *zap.Logger) (*Database, error) {
	if computeLayer == nil {
		return nil, errors.New("compute layer is invalid")
	}

	if storageLayer == nil {
		return nil, errors.New("storage layer is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Database{
		computeLayer: computeLayer,
		storageLayer: storageLayer,
		logger:       logger,
	}, nil
}

// HandleQuery -- выполняет запрос от клиента
func (d *Database) HandleQuery(ctx context.Context, queryStr string) string {
	d.logger.Debug("handling query", zap.String("query", queryStr))
	query, err := d.computeLayer.Parse(queryStr)
	if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}
	switch query.CommandID() {
	case compute.DelCommandID:
		return d.handleDelQuery(ctx, query)
	case compute.GetCommandID:
		return d.handleGetQuery(ctx, query)
	case compute.SetCommandID:
		return d.handleSetQuery(ctx, query)
	}
	d.logger.Error(
		"compute layer is incorrect",
		zap.Int("command_id", query.CommandID()),
	)

	return "[error] internal error"
}

func (d *Database) handleSetQuery(ctx context.Context, query compute.Query) string {
	if err := d.storageLayer.Set(ctx, query.GetKey(), query.GetValue()); err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}
	return "[ok]"
}

func (d *Database) handleGetQuery(ctx context.Context, query compute.Query) string {
	value, err := d.storageLayer.Get(ctx, query.GetKey())
	if err == nil {
		return fmt.Sprintf("[ok] %s", value)
	}
	if err == storage.ErrorNotExist {
		return fmt.Sprintf("[error] %s", err.Error())
	}
	return fmt.Sprintf("[error] %s", err.Error())
}

func (d *Database) handleDelQuery(ctx context.Context, query compute.Query) string {
	if err := d.storageLayer.Del(ctx, query.GetKey()); err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return "[ok]"
}
