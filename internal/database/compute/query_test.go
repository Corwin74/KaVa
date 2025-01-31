package compute

import (
	"testing"

	"github.com/stretchr/testify/require"
)


func TestNewQuery(t *testing.T) {
	t.Parallel()

	query := NewQuery(SetCommandID, "key_test", "value_test")
	require.Equal(t, SetCommandID, query.CommandID())
	require.Equal(t, "key_test", query.GetKey())
	require.Equal(t, "value_test", query.GetValue())
}