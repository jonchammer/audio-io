package bytes

import (
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestWriter_Init(t *testing.T) {
	w := Writer{}
	require.Equal(t, 0, w.Cap())
	require.Equal(t, 0, w.Len())
	require.Nil(t, w.Bytes())
}

func TestWriter_Normal(t *testing.T) {
	w := Writer{}

	// Basic append
	n, _ := w.Write([]byte{0x00, 0x01, 0x02})
	require.Equal(t, 3, n)
	require.GreaterOrEqual(t, w.Cap(), 3)
	require.Equal(t, 3, w.Len())
	require.Equal(t, []byte{0x00, 0x01, 0x02}, w.Bytes())

	// Seek and overwrite existing bytes
	idx, err := w.Seek(1, io.SeekStart)
	require.NoError(t, err)
	require.Equal(t, int64(1), idx)

	n, _ = w.Write([]byte{0xFF})
	require.Equal(t, 1, n)
	require.GreaterOrEqual(t, w.Cap(), 3)
	require.Equal(t, 3, w.Len())
	require.Equal(t, []byte{0x00, 0xFF, 0x02}, w.Bytes())

	// Partially overwrite
	n, _ = w.Write([]byte{0xFE, 0xFD})
	require.Equal(t, 2, n)
	require.GreaterOrEqual(t, w.Cap(), 4)
	require.Equal(t, 4, w.Len())
	require.Equal(t, []byte{0x00, 0xFF, 0xFE, 0xFD}, w.Bytes())
}

func TestWriter_Grow(t *testing.T) {
	w := Writer{}
	require.Equal(t, 0, w.Cap())
	require.Equal(t, 0, w.Len())

	w.Grow(10)
	require.GreaterOrEqual(t, w.Cap(), 10)
	require.Equal(t, 0, w.Len())
}

func TestSeek_Whence_Normal(t *testing.T) {
	w := Writer{}
	n, _ := w.Write([]byte{0x00, 0x01, 0x02, 0x03})
	require.Equal(t, 4, n)

	idx, err := w.Seek(4, io.SeekStart)
	require.NoError(t, err)
	require.Equal(t, int64(4), idx)

	idx, err = w.Seek(-2, io.SeekCurrent)
	require.NoError(t, err)
	require.Equal(t, int64(2), idx)

	idx, err = w.Seek(-1, io.SeekEnd)
	require.NoError(t, err)
	require.Equal(t, int64(3), idx)
}

func TestSeek_Whence_OOB(t *testing.T) {
	w := Writer{}
	n, _ := w.Write([]byte{0x00, 0x01, 0x02, 0x03})
	require.Equal(t, 4, n)

	_, err := w.Seek(-1, io.SeekStart)
	require.Error(t, err)

	_, err = w.Seek(5, io.SeekStart)
	require.Error(t, err)
}
