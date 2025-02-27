package client

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewTCPClient(t *testing.T) {
	// Создаем тестовый TCP сервер
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer listener.Close()

	t.Run("Successful connection", func(t *testing.T) {
		client, err := NewTCPClient(listener.Addr().String(), 1024, time.Second*30)
		if err != nil {
			t.Errorf("Expected successful connection, got error: %v", err)
		}
		if client == nil {
			t.Error("Expected client instance, got nil")
		}
		client.Close()
	})

	t.Run("Invalid address", func(t *testing.T) {
		client, err := NewTCPClient("invalid:address", 1024, time.Second*30)
		if err == nil {
			t.Error("Expected error for invalid address, got nil")
		}
		if client != nil {
			t.Error("Expected nil client for invalid address")
			client.Close()
		}
	})

	t.Run("Zero buffer size", func(t *testing.T) {
		client, err := NewTCPClient(listener.Addr().String(), 0, time.Second*30)
		if err != nil {
			t.Errorf("Expected successful connection with zero buffer, got error: %v", err)
		}
		client.Close()
	})
}

func TestSend(t *testing.T) {
    // Создаем тестовый TCP сервер
    listener, err := net.Listen("tcp", "localhost:8787")
    if err != nil {
        t.Fatalf("Failed to create test server: %v", err)
    }
    defer listener.Close()

    // Запускаем горутину для обработки подключений
    go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			defer conn.Close()

			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			require.NoError(t, err)

			// Эхо-ответ
			_, err = conn.Write(buffer[:n])
			require.NoError(t, err)
		}
    }()

	t.Run("Buffer size too small", func(t *testing.T) {
        client, err := NewTCPClient(listener.Addr().String(), 1, time.Second*30)
        if err != nil {
            t.Fatalf("Failed to create client: %v", err)
        }
        defer client.Close()

        testData := []byte("test message")
        _, err = client.Send(testData)
        if err == nil || err.Error() != "small buffer size" {
            t.Errorf("Expected 'small buffer size' error, got: %v", err)
        }
    })

    t.Run("Successful send and receive", func(t *testing.T) {
        client, err := NewTCPClient(listener.Addr().String(), 1024, time.Second*30)
        if err != nil {
            t.Fatalf("Failed to create client: %v", err)
        }
        defer client.Close()

        testData := []byte("test message")
        response, err := client.Send(testData)
        if err != nil {
            t.Errorf("Expected successful send, got error: %v", err)
        }
        if string(response) != string(testData) {
            t.Errorf("Expected response %q, got %q", testData, response)
        }
    })
}

func TestClose(t *testing.T) {
    listener, err := net.Listen("tcp", "localhost:0")
    if err != nil {
        t.Fatalf("Failed to create test server: %v", err)
    }
    defer listener.Close()

    t.Run("Close connection", func(t *testing.T) {
        client, err := NewTCPClient(listener.Addr().String(), 1024, time.Second*30)
        if err != nil {
            t.Fatalf("Failed to create client: %v", err)
        }

        client.Close()
        
        // Попытка отправить данные после закрытия должна вернуть ошибку
        _, err = client.Send([]byte("test"))
        if err == nil {
            t.Error("Expected error when sending after close, got nil")
        }
    })

    t.Run("Double close", func(t *testing.T) {
        client, err := NewTCPClient(listener.Addr().String(), 1024, time.Second*30)
        if err != nil {
            t.Fatalf("Failed to create client: %v", err)
        }

        // Двойное закрытие не должно вызывать панику
        client.Close()
        client.Close()
    })
}

func TestTimeout(t *testing.T) {
    listener, err := net.Listen("tcp", "localhost:0")
    if err != nil {
        t.Fatalf("Failed to create test server: %v", err)
    }
    defer listener.Close()

    // Запускаем медленный сервер
    go func() {
        conn, err := listener.Accept()
        if err != nil {
            return
        }
        defer conn.Close()
        
        // Имитируем задержку
        time.Sleep(time.Second * 2)
    }()

    t.Run("Connection timeout", func(t *testing.T) {
        client, err := NewTCPClient(listener.Addr().String(), 1024, time.Second)
        if err != nil {
            t.Fatalf("Failed to create client: %v", err)
        }
        defer client.Close()

        _, err = client.Send([]byte("test"))
        if err == nil {
            t.Error("Expected timeout error, got nil")
        }
    })
}