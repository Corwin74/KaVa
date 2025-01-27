package compute

import (
	"errors"
	"strings"

	"go.uber.org/zap"
)

var (
	errInvalidQuery     = errors.New("invalid query")
	errInvalidCommand   = errors.New("invalid command")
	errInvalidArguments = errors.New("invalid arguments")
)


// Compute - парсит и валидирует запрос, передают дальше в storage
type Compute struct {
	logger *zap.Logger
}


// NewCompute - конструктор
func NewCompute(logger *zap.Logger) (*Compute, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Compute{
		logger: logger,
	}, nil
}

// Parse - парсит запрос на команду и аргументы
func (d *Compute) Parse(queryStr string) (Query, error) {
	tokens := strings.Fields(queryStr)
	if len(tokens) < 2 {
		d.logger.Debug("invalid query", zap.String("query", queryStr))
		return Query{}, errInvalidQuery
	}
	commandID, exist := commandTextToID[tokens[0]]
	if !exist {
		d.logger.Debug("command not found", zap.String("query", queryStr))
		return Query{}, errInvalidCommand
	}
	if commandArgumentsCount[commandID] != len(tokens) - 1 {
		d.logger.Debug("invalid number of arguments for the query", zap.String("query", queryStr))
		return Query{}, errInvalidArguments
	}
	query := Query{
		commandID: commandID,
		key: tokens[1],
	}
	if len(tokens) == 3 {
		query.value = tokens[2]
	}
	return query, nil
}