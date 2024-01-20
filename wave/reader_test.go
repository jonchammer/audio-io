package wave

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

// ------------------------------------------------------------------------- //
// Header()
// ------------------------------------------------------------------------- //

func TestReader_Header_InvalidHeader(t *testing.T) {

	// Invalid riff header
	payload := []byte{
		' ', ' ', ' ', ' ',
	}
	r := NewReader(bytes.NewReader(payload))
	_, err := r.Header()
	require.ErrorIs(t, err, ErrRIFFChunkCorruptedHeader)

	// Corrupted RIFF chunk data
	payload = []byte{
		'R', 'I', 'F', 'F',
		0x08, 0x00, 0x00, 0x00,
		'W', 'A', 'V', 'E',
		'a', 'b', 'c', 'd',
		0x00, 0x00, 0x00, 0x00,
	}
	r = NewReader(bytes.NewReader(payload))
	_, err = r.Header()
	require.ErrorIs(t, err, ErrHeaderMissingFmtChunk)
}

// ------------------------------------------------------------------------- //
// ReadUint8()
// ------------------------------------------------------------------------- //

func TestReader_ReadUint8_InvalidHeader(t *testing.T) {
	payload := []byte{
		' ', ' ', ' ', ' ',
	}
	r := NewReader(bytes.NewReader(payload))
	_, err := r.ReadUint8(make([]uint8, 8))
	require.ErrorIs(t, err, ErrRIFFChunkCorruptedHeader)
}

func TestReader_ReadUint8_InvalidSampleType(t *testing.T) {

	// Invalid format code
	payload := []byte{
		'R', 'I', 'F', 'F',
		20, 0x00, 0x00, 0x00,
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		0x10, 0x00, 0x00, 0x00, // 16 bytes
		0xFF, 0xFF, 0xFF, 0xFF, // Invalid format code
		0x44, 0xAC, 0x00, 0x00,
		0x10, 0xB1, 0x02, 0x00,
		0x04, 0x00, 0x10, 0x00,
	}
	r := NewReader(bytes.NewReader(payload))
	_, err := r.ReadUint8(make([]uint8, 8))
	require.ErrorContains(t, err, "invalid format code: 'FormatCode(65535)'")

	// Sample type is int16
	payload = []byte{
		'R', 'I', 'F', 'F',
		28, 0x00, 0x00, 0x00,
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		0x10, 0x00, 0x00, 0x00, // 16 bytes
		0x01, 0x00, // PCM data
		0x01, 0x00, // 1 channel
		0x44, 0xAC, 0x00, 0x00, // 44,100 samples/sec
		0x88, 0x58, 0x01, 0x00, // 88,200 bytes/sec
		0x02, 0x00, // 2 bytes / frame
		0x10, 0x00, // 16 Bits per sample
	}
	r = NewReader(bytes.NewReader(payload))
	_, err = r.ReadUint8(make([]uint8, 8))
	require.ErrorIs(t, err, ErrReaderUnexpectedUint8)
}

// ------------------------------------------------------------------------- //
// ReadInt16()
// ------------------------------------------------------------------------- //

func TestReader_ReadInt16_InvalidHeader(t *testing.T) {
	payload := []byte{
		' ', ' ', ' ', ' ',
	}
	r := NewReader(bytes.NewReader(payload))
	_, err := r.ReadInt16(make([]int16, 8))
	require.ErrorIs(t, err, ErrRIFFChunkCorruptedHeader)
}

func TestReader_ReadInt16_InvalidSampleType(t *testing.T) {

	// Invalid format code
	payload := []byte{
		'R', 'I', 'F', 'F',
		20, 0x00, 0x00, 0x00,
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		0x10, 0x00, 0x00, 0x00, // 16 bytes
		0xFF, 0xFF, 0xFF, 0xFF, // Invalid format code
		0x44, 0xAC, 0x00, 0x00,
		0x10, 0xB1, 0x02, 0x00,
		0x04, 0x00, 0x10, 0x00,
	}
	r := NewReader(bytes.NewReader(payload))
	_, err := r.ReadInt16(make([]int16, 8))
	require.ErrorContains(t, err, "invalid format code: 'FormatCode(65535)'")

	// Sample type is actually uint8
	payload = []byte{
		'R', 'I', 'F', 'F',
		28, 0x00, 0x00, 0x00,
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		0x10, 0x00, 0x00, 0x00, // 16 bytes
		0x01, 0x00, // PCM data
		0x01, 0x00, // 1 channel
		0x44, 0xAC, 0x00, 0x00, // 44,100 samples/sec
		0x44, 0xAC, 0x00, 0x00, // 44,100 bytes/sec
		0x01, 0x00, // 1 byte / frame
		0x08, 0x00, // 8 Bits per sample
	}
	r = NewReader(bytes.NewReader(payload))
	_, err = r.ReadInt16(make([]int16, 8))
	require.ErrorIs(t, err, ErrReaderUnexpectedInt16)
}

// ------------------------------------------------------------------------- //
// ReadInt24()
// ------------------------------------------------------------------------- //

func TestReader_ReadInt24_InvalidHeader(t *testing.T) {
	payload := []byte{
		' ', ' ', ' ', ' ',
	}
	r := NewReader(bytes.NewReader(payload))
	_, err := r.ReadInt24(make([]int32, 8))
	require.ErrorIs(t, err, ErrRIFFChunkCorruptedHeader)
}

func TestReader_ReadInt24_InvalidSampleType(t *testing.T) {

	// Invalid format code
	payload := []byte{
		'R', 'I', 'F', 'F',
		20, 0x00, 0x00, 0x00,
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		0x10, 0x00, 0x00, 0x00, // 16 bytes
		0xFF, 0xFF, 0xFF, 0xFF, // Invalid format code
		0x44, 0xAC, 0x00, 0x00,
		0x10, 0xB1, 0x02, 0x00,
		0x04, 0x00, 0x10, 0x00,
	}
	r := NewReader(bytes.NewReader(payload))
	_, err := r.ReadInt24(make([]int32, 8))
	require.ErrorContains(t, err, "invalid format code: 'FormatCode(65535)'")

	// Sample type is actually uint8
	payload = []byte{
		'R', 'I', 'F', 'F',
		28, 0x00, 0x00, 0x00,
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		0x10, 0x00, 0x00, 0x00, // 16 bytes
		0x01, 0x00, // PCM data
		0x01, 0x00, // 1 channel
		0x44, 0xAC, 0x00, 0x00, // 44,100 samples/sec
		0x44, 0xAC, 0x00, 0x00, // 44,100 bytes/sec
		0x01, 0x00, // 1 byte / frame
		0x08, 0x00, // 8 Bits per sample
	}
	r = NewReader(bytes.NewReader(payload))
	_, err = r.ReadInt24(make([]int32, 8))
	require.ErrorIs(t, err, ErrReaderUnexpectedInt24)
}

// ------------------------------------------------------------------------- //
// ReadInt32()
// ------------------------------------------------------------------------- //

func TestReader_ReadInt32_InvalidHeader(t *testing.T) {
	payload := []byte{
		' ', ' ', ' ', ' ',
	}
	r := NewReader(bytes.NewReader(payload))
	_, err := r.ReadInt32(make([]int32, 8))
	require.ErrorIs(t, err, ErrRIFFChunkCorruptedHeader)
}

func TestReader_ReadInt32_InvalidSampleType(t *testing.T) {

	// Invalid format code
	payload := []byte{
		'R', 'I', 'F', 'F',
		20, 0x00, 0x00, 0x00,
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		0x10, 0x00, 0x00, 0x00, // 16 bytes
		0xFF, 0xFF, 0xFF, 0xFF, // Invalid format code
		0x44, 0xAC, 0x00, 0x00,
		0x10, 0xB1, 0x02, 0x00,
		0x04, 0x00, 0x10, 0x00,
	}
	r := NewReader(bytes.NewReader(payload))
	_, err := r.ReadInt32(make([]int32, 8))
	require.ErrorContains(t, err, "invalid format code: 'FormatCode(65535)'")

	// Sample type is actually uint8
	payload = []byte{
		'R', 'I', 'F', 'F',
		28, 0x00, 0x00, 0x00,
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		0x10, 0x00, 0x00, 0x00, // 16 bytes
		0x01, 0x00, // PCM data
		0x01, 0x00, // 1 channel
		0x44, 0xAC, 0x00, 0x00, // 44,100 samples/sec
		0x44, 0xAC, 0x00, 0x00, // 44,100 bytes/sec
		0x01, 0x00, // 1 byte / frame
		0x08, 0x00, // 8 Bits per sample
	}
	r = NewReader(bytes.NewReader(payload))
	_, err = r.ReadInt32(make([]int32, 8))
	require.ErrorIs(t, err, ErrReaderUnexpectedInt32)
}

// ------------------------------------------------------------------------- //
// ReadFloat32()
// ------------------------------------------------------------------------- //

func TestReader_ReadFloat32_InvalidHeader(t *testing.T) {
	payload := []byte{
		' ', ' ', ' ', ' ',
	}
	r := NewReader(bytes.NewReader(payload))
	_, err := r.ReadFloat32(make([]float32, 8))
	require.ErrorIs(t, err, ErrRIFFChunkCorruptedHeader)
}

func TestReader_ReadFloat32_InvalidSampleType(t *testing.T) {

	// Invalid format code
	payload := []byte{
		'R', 'I', 'F', 'F',
		20, 0x00, 0x00, 0x00,
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		0x10, 0x00, 0x00, 0x00, // 16 bytes
		0xFF, 0xFF, 0xFF, 0xFF, // Invalid format code
		0x44, 0xAC, 0x00, 0x00,
		0x10, 0xB1, 0x02, 0x00,
		0x04, 0x00, 0x10, 0x00,
	}
	r := NewReader(bytes.NewReader(payload))
	_, err := r.ReadFloat32(make([]float32, 8))
	require.ErrorContains(t, err, "invalid format code: 'FormatCode(65535)'")

	// Sample type is actually uint8
	payload = []byte{
		'R', 'I', 'F', 'F',
		28, 0x00, 0x00, 0x00,
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		0x10, 0x00, 0x00, 0x00, // 16 bytes
		0x01, 0x00, // PCM data
		0x01, 0x00, // 1 channel
		0x44, 0xAC, 0x00, 0x00, // 44,100 samples/sec
		0x44, 0xAC, 0x00, 0x00, // 44,100 bytes/sec
		0x01, 0x00, // 1 byte / frame
		0x08, 0x00, // 8 Bits per sample
	}
	r = NewReader(bytes.NewReader(payload))
	_, err = r.ReadFloat32(make([]float32, 8))
	require.ErrorIs(t, err, ErrReaderUnexpectedFloat32)
}

// ------------------------------------------------------------------------- //
// ReadFloat64()
// ------------------------------------------------------------------------- //

func TestReader_ReadFloat64_InvalidHeader(t *testing.T) {
	payload := []byte{
		' ', ' ', ' ', ' ',
	}
	r := NewReader(bytes.NewReader(payload))
	_, err := r.ReadFloat64(make([]float64, 8))
	require.ErrorIs(t, err, ErrRIFFChunkCorruptedHeader)
}

func TestReader_ReadFloat64_InvalidSampleType(t *testing.T) {

	// Invalid format code
	payload := []byte{
		'R', 'I', 'F', 'F',
		20, 0x00, 0x00, 0x00,
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		0x10, 0x00, 0x00, 0x00, // 16 bytes
		0xFF, 0xFF, 0xFF, 0xFF, // Invalid format code
		0x44, 0xAC, 0x00, 0x00,
		0x10, 0xB1, 0x02, 0x00,
		0x04, 0x00, 0x10, 0x00,
	}
	r := NewReader(bytes.NewReader(payload))
	_, err := r.ReadFloat64(make([]float64, 8))
	require.ErrorContains(t, err, "invalid format code: 'FormatCode(65535)'")

	// Sample type is actually uint8
	payload = []byte{
		'R', 'I', 'F', 'F',
		28, 0x00, 0x00, 0x00,
		'W', 'A', 'V', 'E',
		'f', 'm', 't', ' ',
		0x10, 0x00, 0x00, 0x00, // 16 bytes
		0x01, 0x00, // PCM data
		0x01, 0x00, // 1 channel
		0x44, 0xAC, 0x00, 0x00, // 44,100 samples/sec
		0x44, 0xAC, 0x00, 0x00, // 44,100 bytes/sec
		0x01, 0x00, // 1 byte / frame
		0x08, 0x00, // 8 Bits per sample
	}
	r = NewReader(bytes.NewReader(payload))
	_, err = r.ReadFloat64(make([]float64, 8))
	require.ErrorIs(t, err, ErrReaderUnexpectedFloat64)
}
