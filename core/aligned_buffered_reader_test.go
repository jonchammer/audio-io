package core

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestAlignedBufferedReader_ReadBuffer_Easy(t *testing.T) {

	in := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
	}
	r := NewAlignedBufferedReader(bytes.NewReader(in), 4)

	buffer, err := r.ReadBuffer(4)
	require.NoError(t, err)
	require.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, buffer)
	require.Zero(t, r.LeftoverSize())

	buffer, err = r.ReadBuffer(4)
	require.NoError(t, err)
	require.Equal(t, []byte{0x05, 0x06, 0x07, 0x08}, buffer)
	require.Zero(t, r.LeftoverSize())

	buffer, err = r.ReadBuffer(4)
	require.NoError(t, err)
	require.Empty(t, buffer)
	require.Equal(t, 1, r.LeftoverSize())

	buffer, err = r.ReadBuffer(4)
	require.ErrorIs(t, err, io.EOF)
	require.Empty(t, buffer)
}

func TestAlignedBufferedReader_ReadBuffer_Hard(t *testing.T) {

	in := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
	}
	r := NewAlignedBufferedReader(bytes.NewReader(in), 4)

	buffer, err := r.ReadBuffer(5)
	require.NoError(t, err)
	require.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, buffer)
	require.Equal(t, 1, r.LeftoverSize())

	buffer, err = r.ReadBuffer(5)
	require.NoError(t, err)
	require.Equal(t, []byte{0x05, 0x06, 0x07, 0x08}, buffer)
	require.Equal(t, 1, r.LeftoverSize())

	buffer, err = r.ReadBuffer(5)
	require.ErrorIs(t, err, io.EOF)
	require.Empty(t, buffer)
	require.Equal(t, 1, r.LeftoverSize())
}
