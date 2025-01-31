package compute

import (
	"testing"

	"github.com/stretchr/testify/require"
)


func TestCommand(t *testing.T) {
	t.Parallel()

	require.Equal(t, SetCommandID, commandTextToID["SET"])
	require.Equal(t, GetCommandID, commandTextToID["GET"])
	require.Equal(t, DelCommandID, commandTextToID["DEL"])
}
