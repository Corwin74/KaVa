package server

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"go.uber.org/zap"
)

// Сonsole -- читает запросы из in, пишет ответы в out
type Сonsole struct {
	in     io.Reader
	out    io.Writer
	db     Database
	logger *zap.Logger
}

// NewConsole -- конструктор консоли
func NewConsole(in io.Reader, out io.Writer, db Database, logger *zap.Logger) (*Сonsole, error) {
	return &Сonsole{
		in:     in,
		out:    out,
		db:     db,
		logger: logger,
	}, nil
}

// Start -- запускает консоль
func (c *Сonsole) Start(ctx context.Context) {
	go func() {
		reader := bufio.NewReader(c.in)
		var res string
		for {
			query, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				c.logger.Error("failed to read query", zap.Error(err))
				continue
			}
			res = c.db.HandleQuery(ctx, query)

			fmt.Fprintln(c.out, res)
		}
	}()
	<-ctx.Done()
}
