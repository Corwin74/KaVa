package server

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewConsole(t *testing.T) {
    // Arrange
    logger := zaptest.NewLogger(t)
    db := new(MockDatabase)
    in := os.Stdin
    out := os.Stdout

    // Act
    console, err := NewConsole(in, out, db, logger)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, console)
    assert.Equal(t, in, console.in)
    assert.Equal(t, out, console.out)
    assert.Equal(t, db, console.db)
    assert.Equal(t, logger, console.logger)
}

func TestConsole_Start(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        mockResp string
    }{
        {
            name:     "Simple query",
            input:    "GET course\n",
            expected: "test response\n",
            mockResp: "test response",
        },
        {
            name:     "Empty query",
            input:    "\n",
            expected: "\n",
            mockResp: "",
        },
	}

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            logger := zaptest.NewLogger(t)
            db := new(MockDatabase)
            
            // Создаем входной reader из строки
            in := strings.NewReader(tt.input)
            
            // Создаем выходной buffer для проверки результата
            out := &strings.Builder{}
            
            // Настраиваем мок
            db.On("HandleQuery", mock.Anything, mock.AnythingOfType("string")).Return(tt.mockResp)
            
            console, err := NewConsole(in, out, db, logger)
            assert.NoError(t, err)

            // Act
            ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
            defer cancel()
            
            console.Start(ctx)

            // Assert
            assert.Equal(t, tt.expected, out.String())
            db.AssertExpectations(t)
        })
    }
}

func TestConsole_Start_ContextCancel(t *testing.T) {
    logger := zaptest.NewLogger(t)
    db := new(MockDatabase)
    in := os.Stdin
    out := os.Stdout

    // Act
    console, _ := NewConsole(in, out, db, logger)
    
    // Создаем контекст с возможностью отмены
    ctx, cancel := context.WithCancel(context.Background())
    
    // Создаем канал для отслеживания завершения метода Start
    done := make(chan struct{})
    
    // Запускаем Start в отдельной горутине
    go func() {
        console.Start(ctx)
        close(done)
    }()
    
    // Небольшая пауза, чтобы горутина успела запуститься
    time.Sleep(100 * time.Millisecond)
    
    // Отменяем контекст
    cancel()
    
    // Ожидаем завершения метода Start с таймаутом
    select {
    case <-done:
        // Успешно - метод Start завершился после отмены контекста
    case <-time.After(1 * time.Second):
        t.Error("Start не завершился после отмены контекста")
    }
}


type OneTimeErrorReader struct {
    err       error
    errorSent bool
    mu        sync.Mutex
}

func NewOneTimeErrorReader(err error) *OneTimeErrorReader {
    return &OneTimeErrorReader{
        err:       err,
        errorSent: false,
    }
}

func (r *OneTimeErrorReader) Read(p []byte) (n int, err error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    if !r.errorSent {
        r.errorSent = true
        return 0, r.err
    }
    
    // Блокируемся навсегда после первой ошибки
    select {} 
}

func TestConsole_Start_ReadError(t *testing.T) {
    // Создаем observer core для перехвата логов
    observedZapCore, logs := observer.New(zap.ErrorLevel)
    zapLogger := zap.New(observedZapCore)

    errReader := NewOneTimeErrorReader(errors.New("read error"))

    logger := zapLogger
    db := new(MockDatabase)
    in := errReader
    out := os.Stdout

    // Act
    console, _ := NewConsole(in, out, db, logger)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    done := make(chan struct{})

    go func() {
        console.Start(ctx)
        close(done)
    }()

	// Небольшая пауза, чтобы горутина успела запуститься
	time.Sleep(100 * time.Millisecond)

	// Отменяем контекст
	cancel()

	// Проверяем, что была залогирована ошибка
	if logs.Len() == 0 {
		t.Error("no error was logged")
	}
	logEntry := logs.All()[0]
	if logEntry.Level != zap.ErrorLevel {
		t.Errorf("unexpected log level: got %v, want %v", logEntry.Level, zap.ErrorLevel)
	}
	if !strings.Contains(logEntry.Message, "failed to read query") {
		t.Errorf("unexpected error message: %v", logEntry.Message)
	}
}
