package server

import (
	"context"
	"fmt"
	"kava/internal/configuration"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// TestNewTCPServer - тест конструктора сервера
func TestNewTCPServer(t *testing.T) {
    logger := zap.NewNop()
    mockDB := new(MockDatabase)

    tests := []struct {
        name    string
        cfg     *configuration.TCPServerConfig
        wantErr bool
    }{
        {
            name: "Valid configuration",
            cfg: &configuration.TCPServerConfig{
                Host:           "localhost",
                Port:           8080,
                MaxConnections: 100,
                MaxMessageSize: 1024,
                IdleTimeout:    time.Second * 30,
            },
            wantErr: false,
        },
        {
            name:    "Nil configuration",
            cfg:     nil,
            wantErr: true,
        },
        {
            name: "Invalid port",
            cfg: &configuration.TCPServerConfig{
                Host:           "localhost",
                Port:           -1,
                MaxConnections: 100,
                MaxMessageSize: 1024,
                IdleTimeout:    time.Second * 30,
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            server, err := NewTCPServer(tt.cfg, mockDB, logger)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, server)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, server)
            }
        })
    }
}

// TestTCPServer_HandleConnection - тест обработки соединения
func TestTCPServer_HandleConnection(t *testing.T) {
    logger := zap.NewNop()
    mockDB := new(MockDatabase)
    
    cfg := &configuration.TCPServerConfig{
        Host:           "localhost",
        Port:           0, // Использование порта 0 позволит системе выбрать свободный порт
        MaxConnections: 10,
        MaxMessageSize: 1024,
        IdleTimeout:    time.Second * 30,
    }

    server, err := NewTCPServer(cfg, mockDB, logger)
    assert.NoError(t, err)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()

    // Запуск сервера
    go server.Start(ctx)

    // Получение реального порта
    addr := server.listener.Addr().(*net.TCPAddr)
    
    // Тест успешного соединения и обмена данными
    t.Run("Successful connection and data exchange", func(t *testing.T) {
        conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", addr.Port))
        assert.NoError(t, err)
        defer conn.Close()

        // Настройка ожидаемого поведения мока
        mockDB.On("HandleQuery", mock.Anything, "test query").Return("test response").Once()

        // Отправка тестового запроса
        _, err = conn.Write([]byte("test query"))
        assert.NoError(t, err)

        // Чтение ответа
        buffer := make([]byte, 1024)
        n, err := conn.Read(buffer)
        assert.NoError(t, err)
        
        // Проверка ответа (с учетом добавленного \n)
        assert.Equal(t, "test response\n", string(buffer[:n]))
    })

    // Тест максимального количества соединений
    t.Run("Max connections limit", func(t *testing.T) {
        connections := make([]net.Conn, 0, cfg.MaxConnections+1)
        defer func() {
            for _, conn := range connections {
                conn.Close()
            }
        }()

        // Установка максимального количества соединений
        for i := 0; i < cfg.MaxConnections; i++ {
            conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", addr.Port))
            assert.NoError(t, err)
            connections = append(connections, conn)
        }

        // Попытка получить ответ в дополнительном соединении
        // Должно завершиться по таймауту
        conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", addr.Port))
		// Отправка тестового запроса
		_, err = conn.Write([]byte("test query"))
		assert.NoError(t, err)

		// Чтение ответа
		buffer := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(time.Second * 3))
		_, err = conn.Read(buffer)
		assert.Error(t, err)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				conn.Close()
				return
			}
		}
		conn.Close()
		t.Error("Expected connection to fail due to max connections limit")
    })
}

// TestTCPServer_Shutdown - тест корректного завершения работы сервера
func TestTCPServer_Shutdown(t *testing.T) {
    logger := zap.NewNop()
    mockDB := new(MockDatabase)
    
    cfg := &configuration.TCPServerConfig{
        Host:           "localhost",
        Port:           0,
        MaxConnections: 10,
        MaxMessageSize: 1024,
        IdleTimeout:    time.Second * 30,
    }

    server, err := NewTCPServer(cfg, mockDB, logger)
    assert.NoError(t, err)

    ctx, cancel := context.WithCancel(context.Background())
    
    // Запуск сервера
    go server.Start(ctx)

    // Получение адреса сервера
    addr := server.listener.Addr().(*net.TCPAddr)

    // Установка соединения
    conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", addr.Port))
    assert.NoError(t, err)
    defer conn.Close()

    // Отправка запроса на завершение работы сервера
    cancel()

    // Проверка, что новые соединения больше не принимаются
    time.Sleep(time.Millisecond * 100) // Даем серверу время на завершение
    _, err = net.Dial("tcp", fmt.Sprintf("localhost:%d", addr.Port))
    assert.Error(t, err)
}
