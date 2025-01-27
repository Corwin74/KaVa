package compute

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)



func TestNewCompute(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{
		logger *zap.Logger
		expectedErr error
		expectedNilObj bool
	}{
		"create compute without logger": {
			expectedErr: errors.New("logger is invalid"),
			expectedNilObj: true,
		},
		"create compute": {
			logger: zap.NewNop(),
			expectedErr: nil,
		},
	}
	
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			compute, err := NewCompute(test.logger)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, compute)
			} else {
				assert.NotNil(t, compute)
			}
		})

	}
}

func TestParse(t *testing.T) {
	t.Parallel() 

	tests := map[string]struct{
		queryStr string
		expectedQuery Query
		expectedErr error
	}{
		"empty query": {
			queryStr: "",
			expectedErr: errInvalidQuery,
		},
		"space query": {
			queryStr: "            ",
			expectedErr: errInvalidQuery,
		},
		"command with leading space": {
			queryStr: " GET key",
			expectedQuery: NewQuery(GetCommandID, "key", ""),
		},
		"command with UTF symbols": {
			queryStr: "🍻 key",
			expectedErr: errInvalidCommand,
		},
		"invalid command": {
			queryStr: "INFO key",
			expectedErr: errInvalidCommand,
		},
		"GET without key": {
			queryStr: "GET",
			expectedErr: errInvalidQuery,
		},
		"GET with extra argument": {
			queryStr: "GET key 22",
			expectedErr: errInvalidArguments,
		},		
		"DEL without key": {
			queryStr: "DEL",
			expectedErr: errInvalidQuery,
		},		
		"DEL with extra argument": {
			queryStr: "DEL key value",
			expectedErr: errInvalidArguments,
		},		
		"SET without value": {
			queryStr: "SET key",
			expectedErr: errInvalidArguments,
		},
		"SET without value with space": {
			queryStr: "SET key ",
			expectedErr: errInvalidArguments,
		},		
		"SET with extra argument": {
			queryStr: "SET key value value2",
			expectedErr: errInvalidArguments,
		},
		"SET without argument": {
			queryStr: "SET",
			expectedErr: errInvalidQuery,
		},			
		"SET query": {
			queryStr: "SET key:value value",
			expectedQuery: NewQuery(SetCommandID, "key:value", "value"),
		},
		"GET query": {
			queryStr: "GET key:suffix",
			expectedQuery: NewQuery(GetCommandID, "key:suffix", ""),
		},
		"DEL query": {
			queryStr: "DEL key:suffix",
			expectedQuery: NewQuery(DelCommandID, "key:suffix", ""),
		},
	}
	compute, err := NewCompute(zap.NewNop())
	require.NoError(t, err)
	for name, test := range tests{
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			query, err := compute.Parse(test.queryStr)
			assert.Equal(t, test.expectedErr, err)
			assert.True(t, reflect.DeepEqual(test.expectedQuery, query))
		})
	}
}