package wave

import (
	"encoding/binary"
	"github.com/stretchr/testify/require"
	"testing"
)

// ------------------------------------------------------------------------- //
// Chunk
// ------------------------------------------------------------------------- //

func TestChunk_Serialize(t *testing.T) {
	chunk := Chunk{
		ID:   [4]byte{'a', 'b', 'c', 'd'},
		Size: 4,
		Body: []byte{0x01, 0x02, 0x03, 0x04},
	}

	require.Equal(t, []byte{
		'a', 'b', 'c', 'd',
		0x04, 0x00, 0x00, 0x00,
		0x01, 0x02, 0x03, 0x04,
	}, chunk.Serialize())
}

// ------------------------------------------------------------------------- //
// RIFF Chunk Data
// ------------------------------------------------------------------------- //

func TestRIFFChunkData_Serialize_Normal(t *testing.T) {
	data := RIFFChunkData{
		subChunks: []Chunk{
			{
				ID:   [4]byte{'a', 'b', 'c', 'd'},
				Size: 4,
				Body: []byte{0x01, 0x02, 0x03, 0x04},
			},
			{
				ID:   [4]byte{'e', 'f', 'g', 'h'},
				Size: 4,
				Body: []byte{0x05, 0x06, 0x07, 0x08},
			},
		},
	}
	result, totalSizeBytes := data.Serialize()

	require.Equal(t, []byte{
		'W', 'A', 'V', 'E',
		'a', 'b', 'c', 'd',
		0x04, 0x00, 0x00, 0x00,
		0x01, 0x02, 0x03, 0x04,
		'e', 'f', 'g', 'h',
		0x04, 0x00, 0x00, 0x00,
		0x05, 0x06, 0x07, 0x08,
	}, result)
	require.Equal(t, uint32(28), totalSizeBytes)
}

func TestRIFFChunkData_Serialize_Empty(t *testing.T) {
	result, totalSizeBytes := RIFFChunkData{}.Serialize()
	require.Equal(t, []byte("WAVE"), result)
	require.Equal(t, uint32(4), totalSizeBytes)
}

// ------------------------------------------------------------------------- //
// Format Chunk Data
// ------------------------------------------------------------------------- //

func TestNewFormatChunkData_Uint8(t *testing.T) {

	// <= 2 channels, regular PCM
	data := NewFormatChunkData(2, 44100, SampleTypeUint8)
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodePCM,
		ChannelCount:       2,
		SampleRate:         44100,
		ByteRate:           88200,
		BlockAlign:         2,
		BitsPerSample:      8,
		ValidBitsPerSample: nil,
		ChannelMask:        nil,
		SubFormat:          nil,
	}, data)

	// > 2 channels, Extensible
	data = NewFormatChunkData(4, 44100, SampleTypeUint8)
	validBitsPerSample := uint16(8)
	channelMask := uint32(0)
	subFormat := FormatCodePCM
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeExtensible,
		ChannelCount:       4,
		SampleRate:         44100,
		ByteRate:           176400,
		BlockAlign:         4,
		BitsPerSample:      8,
		ValidBitsPerSample: &validBitsPerSample,
		ChannelMask:        &channelMask,
		SubFormat:          &subFormat,
	}, data)
}

func TestNewFormatChunkData_Int16(t *testing.T) {

	// <= 2 channels, regular PCM
	data := NewFormatChunkData(2, 44100, SampleTypeInt16)
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodePCM,
		ChannelCount:       2,
		SampleRate:         44100,
		ByteRate:           176400,
		BlockAlign:         4,
		BitsPerSample:      16,
		ValidBitsPerSample: nil,
		ChannelMask:        nil,
		SubFormat:          nil,
	}, data)

	// > 2 channels, Extensible
	data = NewFormatChunkData(4, 44100, SampleTypeInt16)
	validBitsPerSample := uint16(16)
	channelMask := uint32(0)
	subFormat := FormatCodePCM
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeExtensible,
		ChannelCount:       4,
		SampleRate:         44100,
		ByteRate:           352800,
		BlockAlign:         8,
		BitsPerSample:      16,
		ValidBitsPerSample: &validBitsPerSample,
		ChannelMask:        &channelMask,
		SubFormat:          &subFormat,
	}, data)
}

func TestNewFormatChunkData_Int24(t *testing.T) {

	// 24-bit int is always extensible
	data := NewFormatChunkData(2, 44100, SampleTypeInt24)
	validBitsPerSample := uint16(24)
	channelMask := uint32(0)
	subFormat := FormatCodePCM
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeExtensible,
		ChannelCount:       2,
		SampleRate:         44100,
		ByteRate:           264600,
		BlockAlign:         6,
		BitsPerSample:      24,
		ValidBitsPerSample: &validBitsPerSample,
		ChannelMask:        &channelMask,
		SubFormat:          &subFormat,
	}, data)
}

func TestNewFormatChunkData_Int32(t *testing.T) {

	// 32-bit int is always extensible
	data := NewFormatChunkData(2, 44100, SampleTypeInt32)
	validBitsPerSample := uint16(32)
	channelMask := uint32(0)
	subFormat := FormatCodePCM
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeExtensible,
		ChannelCount:       2,
		SampleRate:         44100,
		ByteRate:           352800,
		BlockAlign:         8,
		BitsPerSample:      32,
		ValidBitsPerSample: &validBitsPerSample,
		ChannelMask:        &channelMask,
		SubFormat:          &subFormat,
	}, data)
}

func TestNewFormatChunkData_Float32(t *testing.T) {

	// <= 2 channels, regular IEEE float
	data := NewFormatChunkData(2, 44100, SampleTypeFloat32)
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeIEEEFloat,
		ChannelCount:       2,
		SampleRate:         44100,
		ByteRate:           352800,
		BlockAlign:         8,
		BitsPerSample:      32,
		ValidBitsPerSample: nil,
		ChannelMask:        nil,
		SubFormat:          nil,
	}, data)

	// > 2 channels, Extensible
	data = NewFormatChunkData(4, 44100, SampleTypeFloat32)
	validBitsPerSample := uint16(32)
	channelMask := uint32(0)
	subFormat := FormatCodeIEEEFloat
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeExtensible,
		ChannelCount:       4,
		SampleRate:         44100,
		ByteRate:           705600,
		BlockAlign:         16,
		BitsPerSample:      32,
		ValidBitsPerSample: &validBitsPerSample,
		ChannelMask:        &channelMask,
		SubFormat:          &subFormat,
	}, data)
}

func TestNewFormatChunkData_Float64(t *testing.T) {

	// <= 2 channels, regular IEEE float
	data := NewFormatChunkData(2, 44100, SampleTypeFloat64)
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeIEEEFloat,
		ChannelCount:       2,
		SampleRate:         44100,
		ByteRate:           705600,
		BlockAlign:         16,
		BitsPerSample:      64,
		ValidBitsPerSample: nil,
		ChannelMask:        nil,
		SubFormat:          nil,
	}, data)

	// > 2 channels, Extensible
	data = NewFormatChunkData(4, 44100, SampleTypeFloat64)
	validBitsPerSample := uint16(64)
	channelMask := uint32(0)
	subFormat := FormatCodeIEEEFloat
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeExtensible,
		ChannelCount:       4,
		SampleRate:         44100,
		ByteRate:           1411200,
		BlockAlign:         32,
		BitsPerSample:      64,
		ValidBitsPerSample: &validBitsPerSample,
		ChannelMask:        &channelMask,
		SubFormat:          &subFormat,
	}, data)
}

func TestFormatChunkData_EffectiveFormatCode(t *testing.T) {

	// Non-extensible format code - return directly
	data := FormatChunkData{
		FormatCode: FormatCodePCM,
		SubFormat:  nil,
	}
	formatCode, err := data.EffectiveFormatCode()
	require.NoError(t, err)
	require.Equal(t, FormatCodePCM, formatCode)

	// Extensible format code - return sub format
	subFormat := FormatCodePCM
	data = FormatChunkData{
		FormatCode: FormatCodeExtensible,
		SubFormat:  &subFormat,
	}
	formatCode, err = data.EffectiveFormatCode()
	require.NoError(t, err)
	require.Equal(t, FormatCodePCM, formatCode)

	// Extensible format code provided, but missing sub format
	data = FormatChunkData{
		FormatCode: FormatCodeExtensible,
		SubFormat:  nil,
	}
	_, err = data.EffectiveFormatCode()
	require.ErrorIs(t, err, ErrFmtChunkMissingSubFormat)
}

func TestFormatChunkData_ChunkSize(t *testing.T) {

	// PCM
	data := FormatChunkData{
		FormatCode: FormatCodePCM,
	}
	require.Equal(t, uint32(16), data.ChunkSize())

	// IEEE Float
	data = FormatChunkData{
		FormatCode: FormatCodeIEEEFloat,
	}
	require.Equal(t, uint32(18), data.ChunkSize())

	// Extensible
	data = FormatChunkData{
		FormatCode: FormatCodeExtensible,
	}
	require.Equal(t, uint32(40), data.ChunkSize())
}

func TestFormatChunkData_Serialize_PCM(t *testing.T) {

	data := NewFormatChunkData(2, 44100, SampleTypeUint8)
	result, err := data.Serialize()
	require.NoError(t, err)

	// Verify the length and fields are correct
	require.Equal(t, 16, len(result))
	require.Equal(t, uint16(data.FormatCode), binary.LittleEndian.Uint16(result[:2]))
	require.Equal(t, data.ChannelCount, binary.LittleEndian.Uint16(result[2:4]))
	require.Equal(t, data.SampleRate, binary.LittleEndian.Uint32(result[4:8]))
	require.Equal(t, data.ByteRate, binary.LittleEndian.Uint32(result[8:12]))
	require.Equal(t, data.BlockAlign, binary.LittleEndian.Uint16(result[12:14]))
	require.Equal(t, data.BitsPerSample, binary.LittleEndian.Uint16(result[14:16]))
}

func TestFormatChunkData_Serialize_IEEEFloat(t *testing.T) {

	data := NewFormatChunkData(2, 44100, SampleTypeFloat32)
	result, err := data.Serialize()
	require.NoError(t, err)

	// Verify the length and fields are correct
	require.Equal(t, 18, len(result))
	require.Equal(t, uint16(data.FormatCode), binary.LittleEndian.Uint16(result[:2]))
	require.Equal(t, data.ChannelCount, binary.LittleEndian.Uint16(result[2:4]))
	require.Equal(t, data.SampleRate, binary.LittleEndian.Uint32(result[4:8]))
	require.Equal(t, data.ByteRate, binary.LittleEndian.Uint32(result[8:12]))
	require.Equal(t, data.BlockAlign, binary.LittleEndian.Uint16(result[12:14]))
	require.Equal(t, data.BitsPerSample, binary.LittleEndian.Uint16(result[14:16]))
	require.Equal(t, uint16(0), binary.LittleEndian.Uint16(result[16:18]))
}

func TestFormatChunkData_Serialize_Extensible(t *testing.T) {

	data := NewFormatChunkData(4, 44100, SampleTypeUint8)
	result, err := data.Serialize()
	require.NoError(t, err)

	// Verify the length and fields are correct
	require.Equal(t, 40, len(result))
	require.Equal(t, uint16(data.FormatCode), binary.LittleEndian.Uint16(result[:2]))
	require.Equal(t, data.ChannelCount, binary.LittleEndian.Uint16(result[2:4]))
	require.Equal(t, data.SampleRate, binary.LittleEndian.Uint32(result[4:8]))
	require.Equal(t, data.ByteRate, binary.LittleEndian.Uint32(result[8:12]))
	require.Equal(t, data.BlockAlign, binary.LittleEndian.Uint16(result[12:14]))
	require.Equal(t, data.BitsPerSample, binary.LittleEndian.Uint16(result[14:16]))
	require.Equal(t, uint16(22), binary.LittleEndian.Uint16(result[16:18]))
	require.Equal(t, *data.ValidBitsPerSample, binary.LittleEndian.Uint16(result[18:20]))
	require.Equal(t, *data.ChannelMask, binary.LittleEndian.Uint32(result[20:24]))
	require.Equal(t, uint16(*data.SubFormat), binary.LittleEndian.Uint16(result[24:26]))
	require.Equal(t, []byte{
		0x00, 0x00,
		0x00, 0x00, 0x10, 0x00,
		0x80, 0x00, 0x00, 0xAA,
		0x00, 0x38, 0x9B, 0x71,
	}, result[26:40])
}

func TestFormatChunkData_Serialize_InvalidExtensible(t *testing.T) {

	// Missing 'ValidBitsPerSample'
	data := NewFormatChunkData(4, 44100, SampleTypeUint8)
	data.ValidBitsPerSample = nil
	_, err := data.Serialize()
	require.ErrorIs(t, err, ErrFmtChunkInvalidExtensible)

	// Missing 'ChannelMask'
	data = NewFormatChunkData(4, 44100, SampleTypeUint8)
	data.ChannelMask = nil
	_, err = data.Serialize()
	require.ErrorIs(t, err, ErrFmtChunkInvalidExtensible)

	// Missing 'SubFormat'
	data = NewFormatChunkData(4, 44100, SampleTypeUint8)
	data.SubFormat = nil
	_, err = data.Serialize()
	require.ErrorIs(t, err, ErrFmtChunkInvalidExtensible)
}

// ------------------------------------------------------------------------- //
// Fact Chunk Data
// ------------------------------------------------------------------------- //

func TestFactChunkData_ChunkSize(t *testing.T) {
	require.Equal(t, uint32(4), FactChunkData{}.ChunkSize())
}

func TestFactChunkData_Serialize(t *testing.T) {
	data := FactChunkData{
		SampleFrames: 128,
	}
	result := data.Serialize()
	require.Equal(t, uint32(128), binary.LittleEndian.Uint32(result[:4]))
}
