package wave

import (
	"bytes"
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

func TestWritePackedInt24_Normal(t *testing.T) {
	input := []int32{-8388608, 0, 8388607}
	expected := []byte{0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x7F}

	output := &bytes.Buffer{}
	n, err := WritePackedInt24(output, input)
	require.NoError(t, err)
	require.Equal(t, 9, n)
	require.Equal(t, expected, output.Bytes())
}

func TestWritePackedInt24_Error(t *testing.T) {
	input := []int32{-8388608, 0, 8388607}
	output := &errWriter{Err: errors.New("something went wrong")}
	n, err := WritePackedInt24(output, input)
	require.Error(t, err)
	require.Equal(t, 0, n)
}

func TestReadPackedInt24_Normal(t *testing.T) {
	input := []byte{0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x7F}
	expected := []int32{-8388608, 0, 8388607}

	output, err := ReadPackedInt24(input)
	require.NoError(t, err)
	require.Equal(t, expected, output)
}

func TestReadPackedInt24_CorruptInput(t *testing.T) {
	input := []byte{0x00}
	_, err := ReadPackedInt24(input)
	require.Error(t, err)
	require.Equal(t, ErrIOInvalid24BitInput, err)
}
