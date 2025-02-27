package server

import "context"

// Database -- интерфейс базы данных
type Database interface {
	HandleQuery(ctx context.Context, queryStr string) string
}
