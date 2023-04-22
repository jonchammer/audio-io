package wave

import (
	ioBytes "bytes"
	"encoding/binary"
	"github.com/stretchr/testify/require"
	"testing"

	"audio-io/bytes"
)

// ------------------------------------------------------------------------- //
// End-to-end tests - These are used to ensure the writer consistently
// generates the correct .wav files.
// ------------------------------------------------------------------------- //

// ------------------------------------------------------------------------- //
// Misc
// ------------------------------------------------------------------------- //

func TestE2E_Empty(t *testing.T) {

	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		WithSampleRate(44100),
		WithSampleType(SampleTypeUint8),
		WithChannelCount(2),
		WithBaseWriter(baseWriter),
	)
	require.NoError(t, err)

	// Write an empty wav file
	err = w.Flush()
	require.NoError(t, err)

	// Verify the bytes written to the baseWriter
	data := baseWriter.Bytes()
	require.Equal(t, 44, len(data))

	// Check RIFF chunk
	require.Equal(t, []byte("RIFF"), data[:4])
	require.Equal(t, uint32(36), binary.LittleEndian.Uint32(data[4:8]))
	require.Equal(t, []byte("WAVE"), data[8:12])

	// Check fmt chunk
	require.Equal(t, []byte("fmt "), data[12:16])
	require.Equal(t, uint32(16), binary.LittleEndian.Uint32(data[16:20]))
	require.Equal(t, uint16(0x01), binary.LittleEndian.Uint16(data[20:22]))
	require.Equal(t, uint16(2), binary.LittleEndian.Uint16(data[22:24]))
	require.Equal(t, uint32(44100), binary.LittleEndian.Uint32(data[24:28]))
	require.Equal(t, uint32(88200), binary.LittleEndian.Uint32(data[28:32]))
	require.Equal(t, uint16(2), binary.LittleEndian.Uint16(data[32:34]))
	require.Equal(t, uint16(8), binary.LittleEndian.Uint16(data[34:36]))

	// Check data chunk
	require.Equal(t, []byte("data"), data[36:40])
	require.Equal(t, uint32(0), binary.LittleEndian.Uint32(data[40:44]))
}

// ------------------------------------------------------------------------- //
// Uint8
// ------------------------------------------------------------------------- //

func TestE2E_Uint8_Normal(t *testing.T) {

	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		WithSampleRate(44100),
		WithSampleType(SampleTypeUint8),
		WithChannelCount(2),
		WithBaseWriter(baseWriter),
	)
	require.NoError(t, err)

	// Write the file
	err = w.WriteInterleavedUint8([]uint8{0, 1, 127, 128, 254, 255})
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	// Verify the bytes written to the baseWriter
	data := baseWriter.Bytes()
	require.Equal(t, 50, len(data))

	// Check RIFF chunk
	require.Equal(t, []byte("RIFF"), data[:4])
	require.Equal(t, uint32(42), binary.LittleEndian.Uint32(data[4:8]))
	require.Equal(t, []byte("WAVE"), data[8:12])

	// Check fmt chunk
	require.Equal(t, []byte("fmt "), data[12:16])
	require.Equal(t, uint32(16), binary.LittleEndian.Uint32(data[16:20]))
	require.Equal(t, uint16(0x01), binary.LittleEndian.Uint16(data[20:22]))
	require.Equal(t, uint16(2), binary.LittleEndian.Uint16(data[22:24]))
	require.Equal(t, uint32(44100), binary.LittleEndian.Uint32(data[24:28]))
	require.Equal(t, uint32(88200), binary.LittleEndian.Uint32(data[28:32]))
	require.Equal(t, uint16(2), binary.LittleEndian.Uint16(data[32:34]))
	require.Equal(t, uint16(8), binary.LittleEndian.Uint16(data[34:36]))

	// Check data chunk
	require.Equal(t, []byte("data"), data[36:40])
	require.Equal(t, uint32(6), binary.LittleEndian.Uint32(data[40:44]))
	require.Equal(t, []byte{0, 1, 127, 128, 254, 255}, data[44:50])
}

func TestE2E_Uint8_Padding(t *testing.T) {

	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		WithSampleRate(44100),
		WithSampleType(SampleTypeUint8),
		WithChannelCount(1),
		WithBaseWriter(baseWriter),
	)
	require.NoError(t, err)

	// Write the file
	err = w.WriteInterleavedUint8([]uint8{0, 1, 2})
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	// Verify the bytes written to the baseWriter
	data := baseWriter.Bytes()
	require.Equal(t, 48, len(data))

	// Check RIFF chunk
	require.Equal(t, []byte("RIFF"), data[:4])
	require.Equal(t, uint32(40), binary.LittleEndian.Uint32(data[4:8]))
	require.Equal(t, []byte("WAVE"), data[8:12])

	// Check fmt chunk
	require.Equal(t, []byte("fmt "), data[12:16])
	require.Equal(t, uint32(16), binary.LittleEndian.Uint32(data[16:20]))
	require.Equal(t, uint16(0x01), binary.LittleEndian.Uint16(data[20:22]))
	require.Equal(t, uint16(1), binary.LittleEndian.Uint16(data[22:24]))
	require.Equal(t, uint32(44100), binary.LittleEndian.Uint32(data[24:28]))
	require.Equal(t, uint32(44100), binary.LittleEndian.Uint32(data[28:32]))
	require.Equal(t, uint16(1), binary.LittleEndian.Uint16(data[32:34]))
	require.Equal(t, uint16(8), binary.LittleEndian.Uint16(data[34:36]))

	// Check data chunk
	require.Equal(t, []byte("data"), data[36:40])
	require.Equal(t, uint32(3), binary.LittleEndian.Uint32(data[40:44]))
	require.Equal(t, []byte{0, 1, 2}, data[44:47])
	require.Equal(t, byte(0x00), data[47])
}

func TestE2E_Uint8_Extensible(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		WithSampleRate(44100),
		WithSampleType(SampleTypeUint8),
		WithChannelCount(4),
		WithBaseWriter(baseWriter),
	)
	require.NoError(t, err)

	err = w.WriteInterleavedUint8(
		[]uint8{0, 1, 2, 3, 127, 128, 129, 130, 252, 253, 254, 255},
	)
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	// Verify the bytes written to the baseWriter
	data := baseWriter.Bytes()
	require.Equal(t, 92, len(data))

	require.Equal(t, []byte("RIFF"), data[:4])
	require.Equal(t, uint32(84), binary.LittleEndian.Uint32(data[4:8]))
	require.Equal(t, []byte("WAVE"), data[8:12])

	require.Equal(t, []byte("fmt "), data[12:16])
	require.Equal(t, uint32(40), binary.LittleEndian.Uint32(data[16:20]))
	require.Equal(t, uint16(0xFFFE), binary.LittleEndian.Uint16(data[20:22]))
	require.Equal(t, uint16(4), binary.LittleEndian.Uint16(data[22:24]))
	require.Equal(t, uint32(44100), binary.LittleEndian.Uint32(data[24:28]))
	require.Equal(t, uint32(176400), binary.LittleEndian.Uint32(data[28:32]))
	require.Equal(t, uint16(4), binary.LittleEndian.Uint16(data[32:34]))
	require.Equal(t, uint16(8), binary.LittleEndian.Uint16(data[34:36]))
	require.Equal(t, uint16(22), binary.LittleEndian.Uint16(data[36:38]))
	require.Equal(t, uint16(8), binary.LittleEndian.Uint16(data[38:40]))
	require.Equal(t, uint32(0), binary.LittleEndian.Uint32(data[40:44]))
	require.Equal(t, uint16(0x1), binary.LittleEndian.Uint16(data[44:46]))
	require.Equal(t, []byte{
		0x00, 0x00,
		0x00, 0x00, 0x10, 0x00,
		0x80, 0x00, 0x00, 0xAA,
		0x00, 0x38, 0x9B, 0x71,
	}, data[46:60])

	require.Equal(t, []byte("fact"), data[60:64])
	require.Equal(t, uint32(4), binary.LittleEndian.Uint32(data[64:68]))
	require.Equal(t, uint32(3), binary.LittleEndian.Uint32(data[68:72]))

	require.Equal(t, []byte("data"), data[72:76])
	require.Equal(t, uint32(12), binary.LittleEndian.Uint32(data[76:80]))
	require.Equal(t, []uint8{0, 1, 2, 3, 127, 128, 129, 130, 252, 253, 254, 255}, data[80:92])
}

// ------------------------------------------------------------------------- //
// Int16
// ------------------------------------------------------------------------- //

func TestE2E_Int16_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		WithSampleRate(44100),
		WithSampleType(SampleTypeInt16),
		WithChannelCount(2),
		WithBaseWriter(baseWriter),
	)
	require.NoError(t, err)

	err = w.WriteInterleavedInt16([]int16{-32768, -32767, 0, 1, 32766, 32767})
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	// Verify the bytes written to the baseWriter
	data := baseWriter.Bytes()
	require.Equal(t, 56, len(data))

	require.Equal(t, []byte("RIFF"), data[:4])
	require.Equal(t, uint32(48), binary.LittleEndian.Uint32(data[4:8]))
	require.Equal(t, []byte("WAVE"), data[8:12])

	require.Equal(t, []byte("fmt "), data[12:16])
	require.Equal(t, uint32(16), binary.LittleEndian.Uint32(data[16:20]))
	require.Equal(t, uint16(0x01), binary.LittleEndian.Uint16(data[20:22]))
	require.Equal(t, uint16(2), binary.LittleEndian.Uint16(data[22:24]))
	require.Equal(t, uint32(44100), binary.LittleEndian.Uint32(data[24:28]))
	require.Equal(t, uint32(176400), binary.LittleEndian.Uint32(data[28:32]))
	require.Equal(t, uint16(4), binary.LittleEndian.Uint16(data[32:34]))
	require.Equal(t, uint16(16), binary.LittleEndian.Uint16(data[34:36]))

	require.Equal(t, []byte("data"), data[36:40])
	require.Equal(t, uint32(12), binary.LittleEndian.Uint32(data[40:44]))

	readData := make([]int16, 6)
	err = binary.Read(ioBytes.NewReader(data[44:56]), binary.LittleEndian, readData)
	require.NoError(t, err)
	require.Equal(t, []int16{-32768, -32767, 0, 1, 32766, 32767}, readData)
}

func TestE2E_Int16_Extensible(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		WithSampleRate(44100),
		WithSampleType(SampleTypeInt16),
		WithChannelCount(4),
		WithBaseWriter(baseWriter),
	)
	require.NoError(t, err)

	err = w.WriteInterleavedInt16([]int16{
		-32768, -32767, -32766, -32765, 0, 1, 2, 3, 32764, 32765, 32766, 32767,
	})
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	// Verify the bytes written to the baseWriter
	data := baseWriter.Bytes()
	require.Equal(t, 104, len(data))

	require.Equal(t, []byte("RIFF"), data[:4])
	require.Equal(t, uint32(96), binary.LittleEndian.Uint32(data[4:8]))
	require.Equal(t, []byte("WAVE"), data[8:12])

	require.Equal(t, []byte("fmt "), data[12:16])
	require.Equal(t, uint32(40), binary.LittleEndian.Uint32(data[16:20]))
	require.Equal(t, uint16(0xFFFE), binary.LittleEndian.Uint16(data[20:22]))
	require.Equal(t, uint16(4), binary.LittleEndian.Uint16(data[22:24]))
	require.Equal(t, uint32(44100), binary.LittleEndian.Uint32(data[24:28]))
	require.Equal(t, uint32(352800), binary.LittleEndian.Uint32(data[28:32]))
	require.Equal(t, uint16(8), binary.LittleEndian.Uint16(data[32:34]))
	require.Equal(t, uint16(16), binary.LittleEndian.Uint16(data[34:36]))
	require.Equal(t, uint16(22), binary.LittleEndian.Uint16(data[36:38]))
	require.Equal(t, uint16(16), binary.LittleEndian.Uint16(data[38:40]))
	require.Equal(t, uint32(0), binary.LittleEndian.Uint32(data[40:44]))
	require.Equal(t, uint16(0x1), binary.LittleEndian.Uint16(data[44:46]))
	require.Equal(t, []byte{
		0x00, 0x00,
		0x00, 0x00, 0x10, 0x00,
		0x80, 0x00, 0x00, 0xAA,
		0x00, 0x38, 0x9B, 0x71,
	}, data[46:60])

	require.Equal(t, []byte("fact"), data[60:64])
	require.Equal(t, uint32(4), binary.LittleEndian.Uint32(data[64:68]))
	require.Equal(t, uint32(3), binary.LittleEndian.Uint32(data[68:72]))

	require.Equal(t, []byte("data"), data[72:76])
	require.Equal(t, uint32(24), binary.LittleEndian.Uint32(data[76:80]))
	readData := make([]int16, 12)
	err = binary.Read(ioBytes.NewReader(data[80:104]), binary.LittleEndian, readData)
	require.NoError(t, err)
	require.Equal(
		t,
		[]int16{-32768, -32767, -32766, -32765, 0, 1, 2, 3, 32764, 32765, 32766, 32767},
		readData,
	)
}

// ------------------------------------------------------------------------- //
// Int24
// ------------------------------------------------------------------------- //

func TestE2E_Int24_Normal(t *testing.T) {

	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		WithSampleRate(44100),
		WithSampleType(SampleTypeInt24),
		WithChannelCount(2),
		WithBaseWriter(baseWriter),
	)
	require.NoError(t, err)

	err = w.WriteInterleavedInt24([]int32{-8388608, -8388607, 0, 1, 8388606, 8388607})
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	// Verify the bytes written to the baseWriter
	data := baseWriter.Bytes()
	require.Equal(t, 98, len(data))

	require.Equal(t, []byte("RIFF"), data[:4])
	require.Equal(t, uint32(90), binary.LittleEndian.Uint32(data[4:8]))
	require.Equal(t, []byte("WAVE"), data[8:12])

	require.Equal(t, []byte("fmt "), data[12:16])
	require.Equal(t, uint32(40), binary.LittleEndian.Uint32(data[16:20]))
	require.Equal(t, uint16(0xFFFE), binary.LittleEndian.Uint16(data[20:22]))
	require.Equal(t, uint16(2), binary.LittleEndian.Uint16(data[22:24]))
	require.Equal(t, uint32(44100), binary.LittleEndian.Uint32(data[24:28]))
	require.Equal(t, uint32(264600), binary.LittleEndian.Uint32(data[28:32]))
	require.Equal(t, uint16(6), binary.LittleEndian.Uint16(data[32:34]))
	require.Equal(t, uint16(24), binary.LittleEndian.Uint16(data[34:36]))
	require.Equal(t, uint16(22), binary.LittleEndian.Uint16(data[36:38]))
	require.Equal(t, uint16(24), binary.LittleEndian.Uint16(data[38:40]))
	require.Equal(t, uint32(0), binary.LittleEndian.Uint32(data[40:44]))
	require.Equal(t, uint16(0x1), binary.LittleEndian.Uint16(data[44:46]))
	require.Equal(t, []byte{
		0x00, 0x00,
		0x00, 0x00, 0x10, 0x00,
		0x80, 0x00, 0x00, 0xAA,
		0x00, 0x38, 0x9B, 0x71,
	}, data[46:60])

	require.Equal(t, []byte("fact"), data[60:64])
	require.Equal(t, uint32(4), binary.LittleEndian.Uint32(data[64:68]))
	require.Equal(t, uint32(3), binary.LittleEndian.Uint32(data[68:72]))

	require.Equal(t, []byte("data"), data[72:76])
	require.Equal(t, uint32(18), binary.LittleEndian.Uint32(data[76:80]))

	// Read 24-bit data from the buffer into an []int32 to reconstruct the data
	readData, err := ReadPackedInt24(data[80:98])
	require.NoError(t, err)
	require.Equal(
		t,
		[]int32{-8388608, -8388607, 0, 1, 8388606, 8388607},
		readData,
	)
}

func TestE2E_Int24_Padding(t *testing.T) {

	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		WithSampleRate(44100),
		WithSampleType(SampleTypeInt24),
		WithChannelCount(1),
		WithBaseWriter(baseWriter),
	)
	require.NoError(t, err)

	err = w.WriteInterleavedInt24([]int32{-8388608})
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	// Verify the bytes written to the baseWriter
	data := baseWriter.Bytes()
	require.Equal(t, 84, len(data))

	require.Equal(t, []byte("RIFF"), data[:4])
	require.Equal(t, uint32(76), binary.LittleEndian.Uint32(data[4:8]))
	require.Equal(t, []byte("WAVE"), data[8:12])

	require.Equal(t, []byte("fmt "), data[12:16])
	require.Equal(t, uint32(40), binary.LittleEndian.Uint32(data[16:20]))
	require.Equal(t, uint16(0xFFFE), binary.LittleEndian.Uint16(data[20:22]))
	require.Equal(t, uint16(1), binary.LittleEndian.Uint16(data[22:24]))
	require.Equal(t, uint32(44100), binary.LittleEndian.Uint32(data[24:28]))
	require.Equal(t, uint32(132300), binary.LittleEndian.Uint32(data[28:32]))
	require.Equal(t, uint16(3), binary.LittleEndian.Uint16(data[32:34]))
	require.Equal(t, uint16(24), binary.LittleEndian.Uint16(data[34:36]))
	require.Equal(t, uint16(22), binary.LittleEndian.Uint16(data[36:38]))
	require.Equal(t, uint16(24), binary.LittleEndian.Uint16(data[38:40]))
	require.Equal(t, uint32(0), binary.LittleEndian.Uint32(data[40:44]))
	require.Equal(t, uint16(0x1), binary.LittleEndian.Uint16(data[44:46]))
	require.Equal(t, []byte{
		0x00, 0x00,
		0x00, 0x00, 0x10, 0x00,
		0x80, 0x00, 0x00, 0xAA,
		0x00, 0x38, 0x9B, 0x71,
	}, data[46:60])

	require.Equal(t, []byte("fact"), data[60:64])
	require.Equal(t, uint32(4), binary.LittleEndian.Uint32(data[64:68]))
	require.Equal(t, uint32(1), binary.LittleEndian.Uint32(data[68:72]))

	require.Equal(t, []byte("data"), data[72:76])
	require.Equal(t, uint32(3), binary.LittleEndian.Uint32(data[76:80]))

	// Read 24-bit data from the buffer into an []int32 to reconstruct the data
	readData, err := ReadPackedInt24(data[80:83])
	require.NoError(t, err)
	require.Equal(t, []int32{-8388608}, readData)
	require.Equal(t, uint8(0x00), data[83])
}

// ------------------------------------------------------------------------- //
// Int32
// ------------------------------------------------------------------------- //

func TestE2E_Int32_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		WithSampleRate(44100),
		WithSampleType(SampleTypeInt32),
		WithChannelCount(2),
		WithBaseWriter(baseWriter),
	)
	require.NoError(t, err)

	err = w.WriteInterleavedInt32([]int32{
		-2147483648, -2147483647, 0, 1, 2147483646, 2147483647,
	})
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	// Verify the bytes written to the baseWriter
	data := baseWriter.Bytes()
	require.Equal(t, 104, len(data))

	require.Equal(t, []byte("RIFF"), data[:4])
	require.Equal(t, uint32(96), binary.LittleEndian.Uint32(data[4:8]))
	require.Equal(t, []byte("WAVE"), data[8:12])

	require.Equal(t, []byte("fmt "), data[12:16])
	require.Equal(t, uint32(40), binary.LittleEndian.Uint32(data[16:20]))
	require.Equal(t, uint16(0xFFFE), binary.LittleEndian.Uint16(data[20:22]))
	require.Equal(t, uint16(2), binary.LittleEndian.Uint16(data[22:24]))
	require.Equal(t, uint32(44100), binary.LittleEndian.Uint32(data[24:28]))
	require.Equal(t, uint32(352800), binary.LittleEndian.Uint32(data[28:32]))
	require.Equal(t, uint16(8), binary.LittleEndian.Uint16(data[32:34]))
	require.Equal(t, uint16(32), binary.LittleEndian.Uint16(data[34:36]))
	require.Equal(t, uint16(22), binary.LittleEndian.Uint16(data[36:38]))
	require.Equal(t, uint16(32), binary.LittleEndian.Uint16(data[38:40]))
	require.Equal(t, uint32(0), binary.LittleEndian.Uint32(data[40:44]))
	require.Equal(t, uint16(0x1), binary.LittleEndian.Uint16(data[44:46]))
	require.Equal(t, []byte{
		0x00, 0x00,
		0x00, 0x00, 0x10, 0x00,
		0x80, 0x00, 0x00, 0xAA,
		0x00, 0x38, 0x9B, 0x71,
	}, data[46:60])

	require.Equal(t, []byte("fact"), data[60:64])
	require.Equal(t, uint32(4), binary.LittleEndian.Uint32(data[64:68]))
	require.Equal(t, uint32(3), binary.LittleEndian.Uint32(data[68:72]))

	require.Equal(t, []byte("data"), data[72:76])
	require.Equal(t, uint32(24), binary.LittleEndian.Uint32(data[76:80]))
	readData := make([]int32, 6)
	err = binary.Read(ioBytes.NewReader(data[80:104]), binary.LittleEndian, readData)
	require.NoError(t, err)
	require.Equal(
		t,
		[]int32{-2147483648, -2147483647, 0, 1, 2147483646, 2147483647},
		readData,
	)
}

// ------------------------------------------------------------------------- //
// Float32
// ------------------------------------------------------------------------- //

func TestE2E_Float32_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		WithSampleRate(44100),
		WithSampleType(SampleTypeFloat32),
		WithChannelCount(2),
		WithBaseWriter(baseWriter),
	)
	require.NoError(t, err)

	err = w.WriteInterleavedFloat32([]float32{
		-1.0, -0.99, 0, 0.1, 0.99, 1.0,
	})
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	// Verify the bytes written to the baseWriter
	data := baseWriter.Bytes()
	require.Equal(t, 82, len(data))

	require.Equal(t, []byte("RIFF"), data[:4])
	require.Equal(t, uint32(74), binary.LittleEndian.Uint32(data[4:8]))
	require.Equal(t, []byte("WAVE"), data[8:12])

	require.Equal(t, []byte("fmt "), data[12:16])
	require.Equal(t, uint32(18), binary.LittleEndian.Uint32(data[16:20]))
	require.Equal(t, uint16(0x03), binary.LittleEndian.Uint16(data[20:22]))
	require.Equal(t, uint16(2), binary.LittleEndian.Uint16(data[22:24]))
	require.Equal(t, uint32(44100), binary.LittleEndian.Uint32(data[24:28]))
	require.Equal(t, uint32(352800), binary.LittleEndian.Uint32(data[28:32]))
	require.Equal(t, uint16(8), binary.LittleEndian.Uint16(data[32:34]))
	require.Equal(t, uint16(32), binary.LittleEndian.Uint16(data[34:36]))
	require.Equal(t, uint16(0), binary.LittleEndian.Uint16(data[36:38]))

	require.Equal(t, []byte("fact"), data[38:42])
	require.Equal(t, uint32(4), binary.LittleEndian.Uint32(data[42:46]))
	require.Equal(t, uint32(3), binary.LittleEndian.Uint32(data[46:50]))

	require.Equal(t, []byte("data"), data[50:54])
	require.Equal(t, uint32(24), binary.LittleEndian.Uint32(data[54:58]))
	readData := make([]float32, 6)
	err = binary.Read(ioBytes.NewReader(data[58:82]), binary.LittleEndian, readData)
	require.NoError(t, err)
	require.Equal(t, []float32{-1.0, -0.99, 0, 0.1, 0.99, 1.0}, readData)
}

func TestE2E_Float32_Extensible(t *testing.T) {

	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		WithSampleRate(44100),
		WithSampleType(SampleTypeFloat32),
		WithChannelCount(4),
		WithBaseWriter(baseWriter),
	)
	require.NoError(t, err)

	err = w.WriteInterleavedFloat32([]float32{
		-1.0, -0.99, -0.98, -0.97, 0.0, 0.1, 0.2, 0.3, 0.97, 0.98, 0.99, 1.0,
	})
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	// Verify the bytes written to the baseWriter
	data := baseWriter.Bytes()
	require.Equal(t, 128, len(data))

	require.Equal(t, []byte("RIFF"), data[:4])
	require.Equal(t, uint32(120), binary.LittleEndian.Uint32(data[4:8]))
	require.Equal(t, []byte("WAVE"), data[8:12])

	require.Equal(t, []byte("fmt "), data[12:16])
	require.Equal(t, uint32(40), binary.LittleEndian.Uint32(data[16:20]))
	require.Equal(t, uint16(0xFFFE), binary.LittleEndian.Uint16(data[20:22]))
	require.Equal(t, uint16(4), binary.LittleEndian.Uint16(data[22:24]))
	require.Equal(t, uint32(44100), binary.LittleEndian.Uint32(data[24:28]))
	require.Equal(t, uint32(705600), binary.LittleEndian.Uint32(data[28:32]))
	require.Equal(t, uint16(16), binary.LittleEndian.Uint16(data[32:34]))
	require.Equal(t, uint16(32), binary.LittleEndian.Uint16(data[34:36]))
	require.Equal(t, uint16(22), binary.LittleEndian.Uint16(data[36:38]))
	require.Equal(t, uint16(32), binary.LittleEndian.Uint16(data[38:40]))
	require.Equal(t, uint32(0), binary.LittleEndian.Uint32(data[40:44]))
	require.Equal(t, uint16(0x3), binary.LittleEndian.Uint16(data[44:46]))
	require.Equal(t, []byte{
		0x00, 0x00,
		0x00, 0x00, 0x10, 0x00,
		0x80, 0x00, 0x00, 0xAA,
		0x00, 0x38, 0x9B, 0x71,
	}, data[46:60])

	require.Equal(t, []byte("fact"), data[60:64])
	require.Equal(t, uint32(4), binary.LittleEndian.Uint32(data[64:68]))
	require.Equal(t, uint32(3), binary.LittleEndian.Uint32(data[68:72]))

	require.Equal(t, []byte("data"), data[72:76])
	require.Equal(t, uint32(48), binary.LittleEndian.Uint32(data[76:80]))
	readData := make([]float32, 12)
	err = binary.Read(ioBytes.NewReader(data[80:128]), binary.LittleEndian, readData)
	require.NoError(t, err)
	require.Equal(
		t,
		[]float32{-1.0, -0.99, -0.98, -0.97, 0.0, 0.1, 0.2, 0.3, 0.97, 0.98, 0.99, 1.0},
		readData,
	)
}

// ------------------------------------------------------------------------- //
// Float64
// ------------------------------------------------------------------------- //

func TestE2E_Float64_Normal(t *testing.T) {

	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		WithSampleRate(44100),
		WithSampleType(SampleTypeFloat64),
		WithChannelCount(2),
		WithBaseWriter(baseWriter),
	)
	require.NoError(t, err)

	err = w.WriteInterleavedFloat64([]float64{
		-1.0, -0.99, 0, 0.1, 0.99, 1.0,
	})
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	// Verify the bytes written to the baseWriter
	data := baseWriter.Bytes()
	require.Equal(t, 106, len(data))

	require.Equal(t, []byte("RIFF"), data[:4])
	require.Equal(t, uint32(98), binary.LittleEndian.Uint32(data[4:8]))
	require.Equal(t, []byte("WAVE"), data[8:12])

	require.Equal(t, []byte("fmt "), data[12:16])
	require.Equal(t, uint32(18), binary.LittleEndian.Uint32(data[16:20]))
	require.Equal(t, uint16(0x03), binary.LittleEndian.Uint16(data[20:22]))
	require.Equal(t, uint16(2), binary.LittleEndian.Uint16(data[22:24]))
	require.Equal(t, uint32(44100), binary.LittleEndian.Uint32(data[24:28]))
	require.Equal(t, uint32(705600), binary.LittleEndian.Uint32(data[28:32]))
	require.Equal(t, uint16(16), binary.LittleEndian.Uint16(data[32:34]))
	require.Equal(t, uint16(64), binary.LittleEndian.Uint16(data[34:36]))
	require.Equal(t, uint16(0), binary.LittleEndian.Uint16(data[36:38]))

	require.Equal(t, []byte("fact"), data[38:42])
	require.Equal(t, uint32(4), binary.LittleEndian.Uint32(data[42:46]))
	require.Equal(t, uint32(3), binary.LittleEndian.Uint32(data[46:50]))

	require.Equal(t, []byte("data"), data[50:54])
	require.Equal(t, uint32(48), binary.LittleEndian.Uint32(data[54:58]))
	readData := make([]float64, 6)
	err = binary.Read(ioBytes.NewReader(data[58:106]), binary.LittleEndian, readData)
	require.NoError(t, err)
	require.Equal(t, []float64{-1.0, -0.99, 0, 0.1, 0.99, 1.0}, readData)
}

func TestE2E_Float64_Extensible(t *testing.T) {

	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		WithSampleRate(44100),
		WithSampleType(SampleTypeFloat64),
		WithChannelCount(4),
		WithBaseWriter(baseWriter),
	)
	require.NoError(t, err)

	err = w.WriteInterleavedFloat64([]float64{
		-1.0, -0.99, -0.98, -0.97, 0.0, 0.1, 0.2, 0.3, 0.97, 0.98, 0.99, 1.0,
	})
	require.NoError(t, err)
	err = w.Flush()
	require.NoError(t, err)

	// Verify the bytes written to the baseWriter
	data := baseWriter.Bytes()
	require.Equal(t, 176, len(data))

	require.Equal(t, []byte("RIFF"), data[:4])
	require.Equal(t, uint32(168), binary.LittleEndian.Uint32(data[4:8]))
	require.Equal(t, []byte("WAVE"), data[8:12])

	require.Equal(t, []byte("fmt "), data[12:16])
	require.Equal(t, uint32(40), binary.LittleEndian.Uint32(data[16:20]))
	require.Equal(t, uint16(0xFFFE), binary.LittleEndian.Uint16(data[20:22]))
	require.Equal(t, uint16(4), binary.LittleEndian.Uint16(data[22:24]))
	require.Equal(t, uint32(44100), binary.LittleEndian.Uint32(data[24:28]))
	require.Equal(t, uint32(1411200), binary.LittleEndian.Uint32(data[28:32]))
	require.Equal(t, uint16(32), binary.LittleEndian.Uint16(data[32:34]))
	require.Equal(t, uint16(64), binary.LittleEndian.Uint16(data[34:36]))
	require.Equal(t, uint16(22), binary.LittleEndian.Uint16(data[36:38]))
	require.Equal(t, uint16(64), binary.LittleEndian.Uint16(data[38:40]))
	require.Equal(t, uint32(0), binary.LittleEndian.Uint32(data[40:44]))
	require.Equal(t, uint16(0x3), binary.LittleEndian.Uint16(data[44:46]))
	require.Equal(t, []byte{
		0x00, 0x00,
		0x00, 0x00, 0x10, 0x00,
		0x80, 0x00, 0x00, 0xAA,
		0x00, 0x38, 0x9B, 0x71,
	}, data[46:60])

	require.Equal(t, []byte("fact"), data[60:64])
	require.Equal(t, uint32(4), binary.LittleEndian.Uint32(data[64:68]))
	require.Equal(t, uint32(3), binary.LittleEndian.Uint32(data[68:72]))

	require.Equal(t, []byte("data"), data[72:76])
	require.Equal(t, uint32(96), binary.LittleEndian.Uint32(data[76:80]))
	readData := make([]float64, 12)
	err = binary.Read(ioBytes.NewReader(data[80:176]), binary.LittleEndian, readData)
	require.NoError(t, err)
	require.Equal(
		t,
		[]float64{-1.0, -0.99, -0.98, -0.97, 0.0, 0.1, 0.2, 0.3, 0.97, 0.98, 0.99, 1.0},
		readData,
	)
}
