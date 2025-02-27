package server

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockDatabase - мок для Database
type MockDatabase struct {
    mock.Mock
}

// HandleQuery - Мок обработки запроса
func (m *MockDatabase) HandleQuery(ctx context.Context, query string) string {
    args := m.Called(ctx, query)
    return args.String(0)
}
