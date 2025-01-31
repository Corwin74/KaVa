package in_memory

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)


func TestNewEngine(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{
		logger *zap.Logger
		expectedError error
		expectedNilObject bool
	}{
		"create with nil logger": {
			logger: nil,
			expectedError: errors.New("logger is invalid"),
			expectedNilObject: true,
		},
		"create": {
			logger: zap.NewNop(),
			expectedNilObject: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			engine, err := NewEngine(test.logger)
			assert.Equal(t, test.expectedError, err)
			if test.expectedNilObject {
				assert.Nil(t, engine)
			} else {
				assert.NotNil(t, engine)
			}
		})
	}
}


func TestEngineSet(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{
		engine *Engine
		key string
		value string
	}{
		"set": {
			engine: func() *Engine {
				engine, err := NewEngine(zap.NewNop())
				require.NoError(t, err)
				return engine
			}(),
			key: "key",
			value: "value",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			test.engine.Set(ctx, test.key, test.value)
			value, exist := test.engine.Get(ctx, test.key)
			assert.True(t, exist)
			assert.Equal(t, test.value, value)
		})

	}
}

func TestEngineDel(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{
		engine *Engine
		key string
	}{
		"del": {
			engine: func() *Engine {
				engine, err := NewEngine(zap.NewNop())
				require.NoError(t, err)
				return engine
			}(),
			key: "key",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			testValue := "VaLuE"
			test.engine.Set(ctx, test.key, testValue)
			value, exist := test.engine.Get(ctx, test.key)
			assert.True(t, exist)
			assert.Equal(t, testValue, value)
			test.engine.Del(ctx, test.key)
			value, exist = test.engine.Get(ctx, test.key)
			assert.False(t, exist)
			assert.Empty(t, value)
		})

	}
}

func TestEngineGet(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{
		engine *Engine
		key string
	}{
		"get": {
			engine: func() *Engine {
				engine, err := NewEngine(zap.NewNop())
				require.NoError(t, err)
				return engine
			}(),
			key: "key",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			value, exist := test.engine.Get(ctx, test.key)
			assert.False(t, exist)
			assert.Empty(t, value)
			testValue := "get_get"
			test.engine.Set(ctx, test.key, testValue)
			value, exist = test.engine.Get(ctx, test.key)
			assert.True(t, exist)
			assert.Equal(t, testValue, value)
		})
	}
}