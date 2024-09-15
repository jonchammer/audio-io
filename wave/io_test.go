package wave

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

type errWriter struct {
	Err error
}

func (e *errWriter) Write(_ []byte) (int, error) {
	return 0, e.Err
}

// ------------------------------------------------------------------------- //
// Test cases
// ------------------------------------------------------------------------- //

func TestWriteMany_Error(t *testing.T) {
	output := &errWriter{Err: errors.New("something went wrong")}
	n, err := WriteMany(output, []byte{0x00})
	require.Error(t, err)
	require.Equal(t, 0, n)
}
