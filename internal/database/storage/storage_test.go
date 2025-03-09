package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"kava/pkg/concurrency"
)

func TestNewStorage(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	writeAheadLog := NewMockWAL(ctrl)
	writeAheadLog.EXPECT().
		Recover().
		Return(nil, nil)

	tests := map[string]struct {
		engine Engine
		logger *zap.Logger
		wal    WAL

		expectedErr    error
		expectedNilObj bool
	}{
		"create storage without engine": {
			expectedErr:    errors.New("engine is invalid"),
			expectedNilObj: true,
		},
		"create storage without logger": {
			engine:         NewMockEngine(ctrl),
			expectedErr:    errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create engine with wal": {
			engine:      NewMockEngine(ctrl),
			wal:         writeAheadLog,
			logger:      zap.NewNop(),
			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			storage, err := NewStorage(test.engine, test.wal, test.logger)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, storage)
			} else {
				assert.NotNil(t, storage)
			}
		})
	}
}

func TestStorageSet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	tests := map[string]struct {
		engine func() Engine
		wal    func() WAL

		expectedErr error
	}{
		"set without wal": {
			engine: func() Engine {
				engine := NewMockEngine(ctrl)
				engine.EXPECT().
					Set(gomock.Any(), "key", "value")
				return engine
			},
			wal: func() WAL { return nil },
		},
		"set with error from wal": {
			engine: func() Engine { return NewMockEngine(ctrl) },
			wal: func() WAL {
				result := make(chan error, 1)
				result <- errors.New("wal error")
				future := concurrency.NewFuture(result)

				wal := NewMockWAL(ctrl)
				wal.EXPECT().
					Recover().
					Return(nil, nil)
				wal.EXPECT().
					Set(gomock.Any(), "key", "value").
					Return(future)
				return wal
			},
			expectedErr: errors.New("wal error"),
		},
		"set with wal": {
			engine: func() Engine {
				engine := NewMockEngine(ctrl)
				engine.EXPECT().
					Set(gomock.Any(), "key", "value")
				return engine
			},
			wal: func() WAL {
				result := make(chan error, 1)
				result <- nil
				future := concurrency.NewFuture(result)

				wal := NewMockWAL(ctrl)
				wal.EXPECT().
					Recover().
					Return(nil, nil)
				wal.EXPECT().
					Set(gomock.Any(), "key", "value").
					Return(future)
				return wal
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			storage, err := NewStorage(test.engine(), test.wal(), zap.NewNop())
			require.NoError(t, err)

			err = storage.Set(context.Background(), "key", "value")
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestStorageDel(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	tests := map[string]struct {
		engine func() Engine
		wal    func() WAL

		expectedErr error
	}{
		"del without wal": {
			engine: func() Engine {
				engine := NewMockEngine(ctrl)
				engine.EXPECT().
					Del(gomock.Any(), "key")
				return engine
			},
			wal: func() WAL { return nil },
		},
		"del with error from wal": {
			engine: func() Engine { return NewMockEngine(ctrl) },
			wal: func() WAL {
				result := make(chan error, 1)
				result <- errors.New("wal error")
				future := concurrency.NewFuture(result)

				wal := NewMockWAL(ctrl)
				wal.EXPECT().
					Recover().
					Return(nil, nil)
				wal.EXPECT().
					Del(gomock.Any(), "key").
					Return(future)
				return wal
			},
			expectedErr: errors.New("wal error"),
		},
		"del with wal": {
			engine: func() Engine {
				engine := NewMockEngine(ctrl)
				engine.EXPECT().
					Del(gomock.Any(), "key")
				return engine
			},
			wal: func() WAL {
				result := make(chan error, 1)
				result <- nil
				future := concurrency.NewFuture(result)

				wal := NewMockWAL(ctrl)
				wal.EXPECT().
					Recover().
					Return(nil, nil)
				wal.EXPECT().
					Del(gomock.Any(), "key").
					Return(future)
				return wal
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			storage, err := NewStorage(test.engine(), test.wal(), zap.NewNop())
			require.NoError(t, err)

			err = storage.Del(context.Background(), "key")
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestStorageGet(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	tests := map[string]struct {
		engine func() Engine
		wal    WAL

		expectedValue string
		expectedErr   error
	}{
		"get with unexisting element": {
			engine: func() Engine {
				engine := NewMockEngine(ctrl)
				engine.EXPECT().
					Get(gomock.Any(), "key").
					Return("", false)
				return engine
			},
			expectedErr: ErrorNotExist,
		},
		"get with existing element": {
			engine: func() Engine {
				engine := NewMockEngine(ctrl)
				engine.EXPECT().
					Get(gomock.Any(), "key").
					Return("value", true)
				return engine
			},
			expectedValue: "value",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			storage, err := NewStorage(test.engine(), test.wal, zap.NewNop())
			require.NoError(t, err)

			value, err := storage.Get(context.Background(), "key")
			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedValue, value)
		})
	}
}
