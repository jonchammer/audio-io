package wave

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/jonchammer/audio-io/bytes"
	"github.com/jonchammer/audio-io/core"
)

// ------------------------------------------------------------------------- //
// NewWriter
// ------------------------------------------------------------------------- //

func TestNewWriter_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeUint8, 44100,
	)
	require.NoError(t, err)

	// Verify that 'w' was initialized correctly
	require.Equal(t, baseWriter, w.baseWriter)
	require.Equal(t, core.SampleTypeUint8, w.sampleType)
	require.Equal(t, uint32(44100), w.formatChunkData.FrameRate)
	require.Equal(t, uint16(1), w.formatChunkData.ChannelCount)
	require.Nil(t, w.factChunkData)
	require.Equal(t, uint32(0), w.dataBytes)
}

func TestNewWriter_WithFactChunk(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeFloat32, 44100,
	)
	require.NoError(t, err)

	// Verify that 'w' was initialized correctly
	require.Equal(t, baseWriter, w.baseWriter)
	require.Equal(t, core.SampleTypeFloat32, w.sampleType)
	require.Equal(t, uint32(44100), w.formatChunkData.FrameRate)
	require.Equal(t, uint16(1), w.formatChunkData.ChannelCount)
	require.NotNil(t, w.factChunkData)
	require.Equal(t, uint32(0), w.factChunkData.SampleFrames)
	require.Equal(t, uint32(0), w.dataBytes)
}

func TestNewWriter_Errors(t *testing.T) {

	baseWriter := &bytes.Writer{}

	// Invalid option (e.g. bad sample type)
	_, err := NewWriter(
		baseWriter, core.SampleType(-1), 44100,
	)
	require.ErrorIs(t, err, ErrWriterInvalidSampleType)
}

// ------------------------------------------------------------------------- //
// Flush
// ------------------------------------------------------------------------- //

func TestWriter_Flush_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeFloat32, 44100,
	)
	require.NoError(t, err)

	err = w.WriteFloat32([]float32{0.0, 1.0, 0.0, -1.0})
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_Flush_NoData(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeUint8, 44100,
	)
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_Flush_WithPadding(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeUint8, 44100,
	)
	require.NoError(t, err)

	err = w.WriteUint8([]byte{0, 1, 2})
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_Flush_InvalidByteCount(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeUint8, 44100,
		WithChannelCount(2),
	)
	require.NoError(t, err)

	// 3 bytes with 2 channels is invalid
	err = w.WriteUint8([]byte{0, 1, 2})
	require.NoError(t, err)

	err = w.Flush()
	require.ErrorIs(t, err, ErrWriterInvalidByteCount)
}

// ------------------------------------------------------------------------- //
// WriteUint8
// ------------------------------------------------------------------------- //

func TestWriter_WriteUint8_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeUint8, 44100,
	)
	require.NoError(t, err)

	err = w.WriteUint8([]byte{0, 1, 2, 3})
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_WriteUint8_InvalidSampleType(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeFloat32, 44100,
	)
	require.NoError(t, err)

	err = w.WriteUint8([]byte{0, 1, 2, 3})
	require.ErrorIs(t, err, ErrWriterExpectedUint8)
}

// ------------------------------------------------------------------------- //
// WriteInt16
// ------------------------------------------------------------------------- //

func TestWriter_WriteInt16_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeInt16, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInt16([]int16{0, 32737, 0, -32768})
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_WriteInt16_InvalidSampleType(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeFloat32, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInt16([]int16{0, 32737, 0, -32768})
	require.ErrorIs(t, err, ErrWriterExpectedInt16)
}

// ------------------------------------------------------------------------- //
// WriteInt24
// ------------------------------------------------------------------------- //

func TestWriter_WriteInt24_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeInt24, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInt24([]int32{0, 8388607, 0, -8388608})
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_WriteInt24_InvalidSampleType(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeFloat32, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInt24([]int32{0, 8388607, 0, -8388608})
	require.ErrorIs(t, err, ErrWriterExpectedInt24)
}

// ------------------------------------------------------------------------- //
// WriteInt32
// ------------------------------------------------------------------------- //

func TestWriter_WriteInt32_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeInt32, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInt32([]int32{0, 2147483647, 0, -2147483648})
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_WriteInt32_InvalidSampleType(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeFloat32, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInt32([]int32{0, 2147483647, 0, -2147483648})
	require.ErrorIs(t, err, ErrWriterExpectedInt32)
}

// ------------------------------------------------------------------------- //
// WriteFloat32
// ------------------------------------------------------------------------- //

func TestWriter_WriteFloat32_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeFloat32, 44100,
	)
	require.NoError(t, err)

	err = w.WriteFloat32([]float32{0.0, 1.0, 0.0, -1.0})
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_WriteFloat32_InvalidSampleType(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeUint8, 44100,
	)
	require.NoError(t, err)

	err = w.WriteFloat32([]float32{0.0, 1.0, 0.0, -1.0})
	require.ErrorIs(t, err, ErrWriterExpectedFloat32)
}

// ------------------------------------------------------------------------- //
// WriteFloat64
// ------------------------------------------------------------------------- //

func TestWriter_WriteFloat64_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeFloat64, 44100,
	)
	require.NoError(t, err)

	err = w.WriteFloat64([]float64{0.0, 1.0, 0.0, -1.0})
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_WriteFloat64_InvalidSampleType(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, core.SampleTypeUint8, 44100,
	)
	require.NoError(t, err)

	err = w.WriteFloat64([]float64{0.0, 1.0, 0.0, -1.0})
	require.ErrorIs(t, err, ErrWriterExpectedFloat64)
}
