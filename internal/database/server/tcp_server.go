package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"kava/internal/configuration"
	"kava/pkg/concurency"
	"net"
	"time"

	"go.uber.org/zap"
)

// TCPServer -- структура сервера
type TCPServer struct {
	semaphore      concurency.Semaphore
	listener       net.Listener
	maxConnections int
	bufferSize     configuration.MessageSize
	idleTimeout    time.Duration
	database       Database
	logger         *zap.Logger
}

// NewTCPServer -- конструктор сервера
func NewTCPServer(cfg *configuration.TCPServerConfig, database Database, logger *zap.Logger) (*TCPServer, error) {

	if cfg == nil {
		return nil, errors.New("config is invalid")
	}

	server := &TCPServer{
		logger:     logger,
		database:   database,
		bufferSize: cfg.MaxMessageSize,
		idleTimeout: cfg.IdleTimeout,
	}

	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, errors.New("failed to listen")
	}

	server.listener = listener

	if cfg.MaxConnections != 0 {
		server.maxConnections = cfg.MaxConnections
	}
	server.semaphore = concurency.NewSemaphore(cfg.MaxConnections)

	return server, nil
}

// Start - запуск сервера
func (s *TCPServer) Start(ctx context.Context) {
	go func() {
		for {
			connection, err := s.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				s.logger.Error("failed to accept", zap.Error(err))
				continue
			}

			s.semaphore.Acquire()
			go func(connection net.Conn) {
				defer s.semaphore.Release()
				s.handleConnection(ctx, connection)
			}(connection)
		}
	}()
	<-ctx.Done()
	s.listener.Close()
}

func (s *TCPServer) handleConnection(ctx context.Context, connection net.Conn) {
	defer func() {
		if v := recover(); v != nil {
			s.logger.Error("captured panic", zap.Any("panic", v))
		}

		if err := connection.Close(); err != nil {
			s.logger.Warn("failed to close connection", zap.Error(err))
		}
	}()

	request := make([]byte, s.bufferSize)

	// Обработка запросов в одном соединении с клиентом
	for {
		count, err := connection.Read(request)
		if err != nil && err != io.EOF {
			s.logger.Warn(
				"failed to read data",
				zap.String("address", connection.RemoteAddr().String()),
				zap.Error(err),
			)
			break
		}

		response := s.database.HandleQuery(ctx, string(request[:count]))
		if _, err := connection.Write([]byte(response + "\n")); err != nil {
			s.logger.Warn(
				"failed to write data",
				zap.String("address", connection.RemoteAddr().String()),
				zap.Error(err),
			)
			break
		}
	}
}
