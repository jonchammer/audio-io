package wave

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/jonchammer/audio-io/bytes"
)

// ------------------------------------------------------------------------- //
// NewWriter
// ------------------------------------------------------------------------- //

func TestNewWriter_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeUint8, 44100,
	)
	require.NoError(t, err)

	// Verify that 'w' was initialized correctly
	require.Equal(t, baseWriter, w.baseWriter)
	require.Equal(t, SampleTypeUint8, w.sampleType)
	require.Equal(t, uint32(44100), w.formatChunkData.SampleRate)
	require.Equal(t, uint16(1), w.formatChunkData.ChannelCount)
	require.Nil(t, w.factChunkData)
	require.Equal(t, uint32(0), w.dataBytes)
}

func TestNewWriter_WithFactChunk(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeFloat32, 44100,
	)
	require.NoError(t, err)

	// Verify that 'w' was initialized correctly
	require.Equal(t, baseWriter, w.baseWriter)
	require.Equal(t, SampleTypeFloat32, w.sampleType)
	require.Equal(t, uint32(44100), w.formatChunkData.SampleRate)
	require.Equal(t, uint16(1), w.formatChunkData.ChannelCount)
	require.NotNil(t, w.factChunkData)
	require.Equal(t, uint32(0), w.factChunkData.SampleFrames)
	require.Equal(t, uint32(0), w.dataBytes)
}

func TestNewWriter_Errors(t *testing.T) {

	baseWriter := &bytes.Writer{}

	// Invalid option (e.g. bad sample type)
	_, err := NewWriter(
		baseWriter, SampleType(-1), 44100,
	)
	require.ErrorIs(t, err, ErrWriterInvalidSampleType)
}

// ------------------------------------------------------------------------- //
// Flush
// ------------------------------------------------------------------------- //

func TestWriter_Flush_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeFloat32, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInterleavedFloat32([]float32{0.0, 1.0, 0.0, -1.0})
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_Flush_NoData(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeUint8, 44100,
	)
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_Flush_WithPadding(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeUint8, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInterleavedUint8([]byte{0, 1, 2})
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_Flush_InvalidByteCount(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeUint8, 44100,
		WithChannelCount(2),
	)
	require.NoError(t, err)

	// 3 bytes with 2 channels is invalid
	err = w.WriteInterleavedUint8([]byte{0, 1, 2})
	require.NoError(t, err)

	err = w.Flush()
	require.ErrorIs(t, err, ErrWriterInvalidByteCount)
}

// ------------------------------------------------------------------------- //
// WriteInterleavedUint8
// ------------------------------------------------------------------------- //

func TestWriter_WriteInterleavedUint8_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeUint8, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInterleavedUint8([]byte{0, 1, 2, 3})
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_WriteInterleavedUint8_InvalidSampleType(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeFloat32, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInterleavedUint8([]byte{0, 1, 2, 3})
	require.ErrorIs(t, err, ErrWriterExpectedUint8)
}

// ------------------------------------------------------------------------- //
// WriteInterleavedInt16
// ------------------------------------------------------------------------- //

func TestWriter_WriteInterleavedInt16_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeInt16, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInterleavedInt16([]int16{0, 32737, 0, -32768})
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_WriteInterleavedInt16_InvalidSampleType(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeFloat32, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInterleavedInt16([]int16{0, 32737, 0, -32768})
	require.ErrorIs(t, err, ErrWriterExpectedInt16)
}

// ------------------------------------------------------------------------- //
// WriteInterleavedInt24
// ------------------------------------------------------------------------- //

func TestWriter_WriteInterleavedInt24_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeInt24, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInterleavedInt24([]int32{0, 8388607, 0, -8388608})
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_WriteInterleavedInt24_InvalidSampleType(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeFloat32, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInterleavedInt24([]int32{0, 8388607, 0, -8388608})
	require.ErrorIs(t, err, ErrWriterExpectedInt24)
}

// ------------------------------------------------------------------------- //
// WriteInterleavedInt32
// ------------------------------------------------------------------------- //

func TestWriter_WriteInterleavedInt32_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeInt32, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInterleavedInt32([]int32{0, 2147483647, 0, -2147483648})
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_WriteInterleavedInt32_InvalidSampleType(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeFloat32, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInterleavedInt32([]int32{0, 2147483647, 0, -2147483648})
	require.ErrorIs(t, err, ErrWriterExpectedInt32)
}

// ------------------------------------------------------------------------- //
// WriteInterleavedFloat32
// ------------------------------------------------------------------------- //

func TestWriter_WriteInterleavedFloat32_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeFloat32, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInterleavedFloat32([]float32{0.0, 1.0, 0.0, -1.0})
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_WriteInterleavedFloat32_InvalidSampleType(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeUint8, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInterleavedFloat32([]float32{0.0, 1.0, 0.0, -1.0})
	require.ErrorIs(t, err, ErrWriterExpectedFloat32)
}

// ------------------------------------------------------------------------- //
// WriteInterleavedFloat64
// ------------------------------------------------------------------------- //

func TestWriter_WriteInterleavedFloat64_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeFloat64, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInterleavedFloat64([]float64{0.0, 1.0, 0.0, -1.0})
	require.NoError(t, err)
	require.Greater(t, baseWriter.Len(), 0)
}

func TestWriter_WriteInterleavedFloat64_InvalidSampleType(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeUint8, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInterleavedFloat64([]float64{0.0, 1.0, 0.0, -1.0})
	require.ErrorIs(t, err, ErrWriterExpectedFloat64)
}
