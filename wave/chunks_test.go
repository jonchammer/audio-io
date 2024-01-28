package wave

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/require"
	"io"
	"testing"

	"github.com/jonchammer/audio-io/core"
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
		SubChunks: []Chunk{
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

func TestReadRIFFChunk_Normal(t *testing.T) {

	var payload bytes.Buffer
	payload.Write(RIFFChunkID[:])    // "RIFF"
	payload.Write(uint32ToBytes(78)) // Example file size
	payload.Write(WaveID[:])         // "WAVE"
	payload.Write([]byte{            // Add a few example chunks
		'a', 'b', 'c', 'd',
		0x04, 0x00, 0x00, 0x00,
		0x01, 0x02, 0x03, 0x04,
		'e', 'f', 'g', 'h',
		0x04, 0x00, 0x00, 0x00,
		0x05, 0x06, 0x07, 0x08,
	})
	payload.Write(DataChunkID[:])    // "data"
	payload.Write(uint32ToBytes(42)) // Data size in bytes
	payload.Write(make([]byte, 42))

	fileSize, riffChunkData, err := ReadRIFFChunk(bytes.NewReader(payload.Bytes()))
	require.NoError(t, err)
	require.Equal(t, uint32(78+8), fileSize) // +8 for the RIFF header
	require.NotNil(t, riffChunkData)
	require.Equal(t, 3, len(riffChunkData.SubChunks))

	chunk := riffChunkData.SubChunks[0]
	require.Equal(t, [4]byte{'a', 'b', 'c', 'd'}, chunk.ID)
	require.Equal(t, uint32(4), chunk.Size)
	require.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, chunk.Body)

	chunk = riffChunkData.SubChunks[1]
	require.Equal(t, [4]byte{'e', 'f', 'g', 'h'}, chunk.ID)
	require.Equal(t, uint32(4), chunk.Size)
	require.Equal(t, []byte{0x05, 0x06, 0x07, 0x08}, chunk.Body)

	chunk = riffChunkData.SubChunks[2]
	require.Equal(t, [4]byte{'d', 'a', 't', 'a'}, chunk.ID)
	require.Equal(t, uint32(42), chunk.Size)
	require.Empty(t, chunk.Body)
}

func TestReadRIFFChunk_InvalidHeader(t *testing.T) {

	var payload bytes.Buffer

	// Incomplete ID
	payload.Write([]byte{'R', 'I', 'F'})
	_, _, err := ReadRIFFChunk(bytes.NewReader(payload.Bytes()))
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	payload.Reset()

	// Invalid RIFF ID
	payload.Write([]byte{'B', 'A', 'D', ' '})
	_, _, err = ReadRIFFChunk(bytes.NewReader(payload.Bytes()))
	require.ErrorIs(t, err, ErrRIFFChunkCorruptedHeader)
	payload.Reset()

	// Incomplete file size
	payload.Write(RIFFChunkID[:]) // "RIFF"
	payload.Write([]byte{0x04, 0x00, 0x00})
	_, _, err = ReadRIFFChunk(bytes.NewReader(payload.Bytes()))
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	payload.Reset()

	// Incomplete WAVE ID
	payload.Write(RIFFChunkID[:])     // "RIFF"
	payload.Write(uint32ToBytes(100)) // Example file size
	payload.Write([]byte{'W', 'A', 'V'})
	_, _, err = ReadRIFFChunk(bytes.NewReader(payload.Bytes()))
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	payload.Reset()

	// Invalid Wave ID
	payload.Write(RIFFChunkID[:])     // "RIFF"
	payload.Write(uint32ToBytes(100)) // Example file size
	payload.Write([]byte{'B', 'A', 'D', ' '})
	_, _, err = ReadRIFFChunk(bytes.NewReader(payload.Bytes()))
	require.ErrorIs(t, err, ErrRIFFChunkCorruptedHeader)
	payload.Reset()
}

func TestReadRIFFChunk_CorruptedChunk(t *testing.T) {

	var payload bytes.Buffer

	// Incomplete chunk ID
	payload.Write(RIFFChunkID[:])     // "RIFF"
	payload.Write(uint32ToBytes(100)) // Example file size
	payload.Write(WaveID[:])          // "WAVE"
	payload.Write([]byte{'a', 'b', 'c'})
	_, _, err := ReadRIFFChunk(bytes.NewReader(payload.Bytes()))
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	payload.Reset()

	// Incomplete chunk size
	payload.Write(RIFFChunkID[:])     // "RIFF"
	payload.Write(uint32ToBytes(100)) // Example file size
	payload.Write(WaveID[:])          // "WAVE"
	payload.Write([]byte{
		'a', 'b', 'c', 'd',
		0x04, 0x00, 0x00,
	})
	_, _, err = ReadRIFFChunk(bytes.NewReader(payload.Bytes()))
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	payload.Reset()

	// Incomplete chunk body
	payload.Write(RIFFChunkID[:])     // "RIFF"
	payload.Write(uint32ToBytes(100)) // Example file size
	payload.Write(WaveID[:])          // "WAVE"
	payload.Write([]byte{
		'a', 'b', 'c', 'd',
		0x04, 0x00, 0x00, 0x00,
		0x01, 0x02, 0x03,
	})
	_, _, err = ReadRIFFChunk(bytes.NewReader(payload.Bytes()))
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	payload.Reset()
}

func TestReadRIFFChunk_MissingDataChunk(t *testing.T) {
	var payload bytes.Buffer
	payload.Write(RIFFChunkID[:])     // "RIFF"
	payload.Write(uint32ToBytes(100)) // Example file size
	payload.Write(WaveID[:])          // "WAVE"
	payload.Write([]byte{             // Add a few example chunks
		'a', 'b', 'c', 'd',
		0x04, 0x00, 0x00, 0x00,
		0x01, 0x02, 0x03, 0x04,
		'e', 'f', 'g', 'h',
		0x04, 0x00, 0x00, 0x00,
		0x05, 0x06, 0x07, 0x08,
	})

	_, _, err := ReadRIFFChunk(bytes.NewReader(payload.Bytes()))
	require.ErrorIs(t, err, io.EOF)
}

// ------------------------------------------------------------------------- //
// Format Chunk Data
// ------------------------------------------------------------------------- //

func TestNewFormatChunkData_Uint8(t *testing.T) {

	// <= 2 channels, regular PCM
	data := NewFormatChunkData(2, 44100, core.SampleTypeUint8)
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodePCM,
		ChannelCount:       2,
		FrameRate:          44100,
		ByteRate:           88200,
		BlockAlign:         2,
		BitsPerSample:      8,
		ValidBitsPerSample: nil,
		ChannelMask:        nil,
		SubFormat:          nil,
	}, data)

	// > 2 channels, Extensible
	data = NewFormatChunkData(4, 44100, core.SampleTypeUint8)
	validBitsPerSample := uint16(8)
	channelMask := uint32(0)
	subFormat := FormatCodePCM
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeExtensible,
		ChannelCount:       4,
		FrameRate:          44100,
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
	data := NewFormatChunkData(2, 44100, core.SampleTypeInt16)
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodePCM,
		ChannelCount:       2,
		FrameRate:          44100,
		ByteRate:           176400,
		BlockAlign:         4,
		BitsPerSample:      16,
		ValidBitsPerSample: nil,
		ChannelMask:        nil,
		SubFormat:          nil,
	}, data)

	// > 2 channels, Extensible
	data = NewFormatChunkData(4, 44100, core.SampleTypeInt16)
	validBitsPerSample := uint16(16)
	channelMask := uint32(0)
	subFormat := FormatCodePCM
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeExtensible,
		ChannelCount:       4,
		FrameRate:          44100,
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
	data := NewFormatChunkData(2, 44100, core.SampleTypeInt24)
	validBitsPerSample := uint16(24)
	channelMask := uint32(0)
	subFormat := FormatCodePCM
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeExtensible,
		ChannelCount:       2,
		FrameRate:          44100,
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
	data := NewFormatChunkData(2, 44100, core.SampleTypeInt32)
	validBitsPerSample := uint16(32)
	channelMask := uint32(0)
	subFormat := FormatCodePCM
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeExtensible,
		ChannelCount:       2,
		FrameRate:          44100,
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
	data := NewFormatChunkData(2, 44100, core.SampleTypeFloat32)
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeIEEEFloat,
		ChannelCount:       2,
		FrameRate:          44100,
		ByteRate:           352800,
		BlockAlign:         8,
		BitsPerSample:      32,
		ValidBitsPerSample: nil,
		ChannelMask:        nil,
		SubFormat:          nil,
	}, data)

	// > 2 channels, Extensible
	data = NewFormatChunkData(4, 44100, core.SampleTypeFloat32)
	validBitsPerSample := uint16(32)
	channelMask := uint32(0)
	subFormat := FormatCodeIEEEFloat
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeExtensible,
		ChannelCount:       4,
		FrameRate:          44100,
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
	data := NewFormatChunkData(2, 44100, core.SampleTypeFloat64)
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeIEEEFloat,
		ChannelCount:       2,
		FrameRate:          44100,
		ByteRate:           705600,
		BlockAlign:         16,
		BitsPerSample:      64,
		ValidBitsPerSample: nil,
		ChannelMask:        nil,
		SubFormat:          nil,
	}, data)

	// > 2 channels, Extensible
	data = NewFormatChunkData(4, 44100, core.SampleTypeFloat64)
	validBitsPerSample := uint16(64)
	channelMask := uint32(0)
	subFormat := FormatCodeIEEEFloat
	require.Equal(t, FormatChunkData{
		FormatCode:         FormatCodeExtensible,
		ChannelCount:       4,
		FrameRate:          44100,
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

	data := NewFormatChunkData(2, 44100, core.SampleTypeUint8)
	result, err := data.Serialize()
	require.NoError(t, err)

	// Verify the length and fields are correct
	require.Equal(t, 16, len(result))
	require.Equal(t, uint16(data.FormatCode), binary.LittleEndian.Uint16(result[:2]))
	require.Equal(t, data.ChannelCount, binary.LittleEndian.Uint16(result[2:4]))
	require.Equal(t, data.FrameRate, binary.LittleEndian.Uint32(result[4:8]))
	require.Equal(t, data.ByteRate, binary.LittleEndian.Uint32(result[8:12]))
	require.Equal(t, data.BlockAlign, binary.LittleEndian.Uint16(result[12:14]))
	require.Equal(t, data.BitsPerSample, binary.LittleEndian.Uint16(result[14:16]))
}

func TestFormatChunkData_Serialize_IEEEFloat(t *testing.T) {

	data := NewFormatChunkData(2, 44100, core.SampleTypeFloat32)
	result, err := data.Serialize()
	require.NoError(t, err)

	// Verify the length and fields are correct
	require.Equal(t, 18, len(result))
	require.Equal(t, uint16(data.FormatCode), binary.LittleEndian.Uint16(result[:2]))
	require.Equal(t, data.ChannelCount, binary.LittleEndian.Uint16(result[2:4]))
	require.Equal(t, data.FrameRate, binary.LittleEndian.Uint32(result[4:8]))
	require.Equal(t, data.ByteRate, binary.LittleEndian.Uint32(result[8:12]))
	require.Equal(t, data.BlockAlign, binary.LittleEndian.Uint16(result[12:14]))
	require.Equal(t, data.BitsPerSample, binary.LittleEndian.Uint16(result[14:16]))
	require.Equal(t, uint16(0), binary.LittleEndian.Uint16(result[16:18]))
}

func TestFormatChunkData_Serialize_Extensible(t *testing.T) {

	data := NewFormatChunkData(4, 44100, core.SampleTypeUint8)
	result, err := data.Serialize()
	require.NoError(t, err)

	// Verify the length and fields are correct
	require.Equal(t, 40, len(result))
	require.Equal(t, uint16(data.FormatCode), binary.LittleEndian.Uint16(result[:2]))
	require.Equal(t, data.ChannelCount, binary.LittleEndian.Uint16(result[2:4]))
	require.Equal(t, data.FrameRate, binary.LittleEndian.Uint32(result[4:8]))
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
	data := NewFormatChunkData(4, 44100, core.SampleTypeUint8)
	data.ValidBitsPerSample = nil
	_, err := data.Serialize()
	require.ErrorIs(t, err, ErrFmtChunkInvalidExtensible)

	// Missing 'ChannelMask'
	data = NewFormatChunkData(4, 44100, core.SampleTypeUint8)
	data.ChannelMask = nil
	_, err = data.Serialize()
	require.ErrorIs(t, err, ErrFmtChunkInvalidExtensible)

	// Missing 'SubFormat'
	data = NewFormatChunkData(4, 44100, core.SampleTypeUint8)
	data.SubFormat = nil
	_, err = data.Serialize()
	require.ErrorIs(t, err, ErrFmtChunkInvalidExtensible)
}

func TestDeserializeFormatChunk_PCM(t *testing.T) {

	// "fmt" payload for 2 channel, 16-bit PCM samples
	var payload bytes.Buffer
	payload.Write(uint16ToBytes(1))             // Format Code
	payload.Write(uint16ToBytes(2))             // Channel Count
	payload.Write(uint32ToBytes(44100))         // Frame Rate
	payload.Write(uint32ToBytes(44100 * 2 * 2)) // Byte Rate
	payload.Write(uint16ToBytes(2 * 2))         // Block Align
	payload.Write(uint16ToBytes(8 * 2))         // Bits per Sample

	formatChunkData, err := DeserializeFormatChunk(payload.Bytes())
	require.NoError(t, err)
	require.Equal(t, FormatCodePCM, formatChunkData.FormatCode)
	require.Equal(t, uint16(2), formatChunkData.ChannelCount)
	require.Equal(t, uint32(44100), formatChunkData.FrameRate)
	require.Equal(t, uint32(44100*2*2), formatChunkData.ByteRate)
	require.Equal(t, uint16(4), formatChunkData.BlockAlign)
	require.Equal(t, uint16(16), formatChunkData.BitsPerSample)
	require.Nil(t, formatChunkData.ValidBitsPerSample)
	require.Nil(t, formatChunkData.ChannelMask)
	require.Nil(t, formatChunkData.SubFormat)
}

func TestDeserializeFormatChunk_IEEEFloat(t *testing.T) {

	// "fmt" payload for 2 channel, 32-bit IEEE samples
	var payload bytes.Buffer
	payload.Write(uint16ToBytes(3))             // Format Code
	payload.Write(uint16ToBytes(2))             // Channel Count
	payload.Write(uint32ToBytes(44100))         // Frame Rate
	payload.Write(uint32ToBytes(44100 * 4 * 2)) // Byte Rate
	payload.Write(uint16ToBytes(4 * 2))         // Block Align
	payload.Write(uint16ToBytes(32))            // Bits per Sample
	payload.Write(uint16ToBytes(0))             // Extension Size

	formatChunkData, err := DeserializeFormatChunk(payload.Bytes())
	require.NoError(t, err)
	require.Equal(t, FormatCodeIEEEFloat, formatChunkData.FormatCode)
	require.Equal(t, uint16(2), formatChunkData.ChannelCount)
	require.Equal(t, uint32(44100), formatChunkData.FrameRate)
	require.Equal(t, uint32(44100*4*2), formatChunkData.ByteRate)
	require.Equal(t, uint16(8), formatChunkData.BlockAlign)
	require.Equal(t, uint16(32), formatChunkData.BitsPerSample)
	require.Nil(t, formatChunkData.ValidBitsPerSample)
	require.Nil(t, formatChunkData.ChannelMask)
	require.Nil(t, formatChunkData.SubFormat)
}

func TestDeserializeFormatChunk_Extensible(t *testing.T) {

	// "fmt" payload for 4 channel, 16-bit PCM samples
	var payload bytes.Buffer
	payload.Write(uint16ToBytes(0xFFFE))        // Format Code
	payload.Write(uint16ToBytes(4))             // Channel Count
	payload.Write(uint32ToBytes(44100))         // Frame Rate
	payload.Write(uint32ToBytes(44100 * 2 * 4)) // Byte Rate
	payload.Write(uint16ToBytes(2 * 4))         // Block Align
	payload.Write(uint16ToBytes(8 * 2))         // Bits per Sample
	payload.Write(uint16ToBytes(22))            // Extension Size
	payload.Write(uint16ToBytes(8 * 2))         // Valid Bits Per Sample
	payload.Write(uint32ToBytes(0x12345678))    // Speaker Mask
	payload.Write(uint16ToBytes(0x01))          // Sub-format
	payload.Write(make([]byte, 14))             // Remainder of GUID

	formatChunkData, err := DeserializeFormatChunk(payload.Bytes())
	require.NoError(t, err)
	require.Equal(t, FormatCodeExtensible, formatChunkData.FormatCode)
	require.Equal(t, uint16(4), formatChunkData.ChannelCount)
	require.Equal(t, uint32(44100), formatChunkData.FrameRate)
	require.Equal(t, uint32(44100*2*4), formatChunkData.ByteRate)
	require.Equal(t, uint16(8), formatChunkData.BlockAlign)
	require.Equal(t, uint16(16), formatChunkData.BitsPerSample)
	require.NotNil(t, formatChunkData.ValidBitsPerSample)
	require.Equal(t, uint16(16), *formatChunkData.ValidBitsPerSample)
	require.NotNil(t, formatChunkData.ChannelMask)
	require.Equal(t, uint32(0x12345678), *formatChunkData.ChannelMask)
	require.NotNil(t, formatChunkData.SubFormat)
	require.Equal(t, FormatCodePCM, *formatChunkData.SubFormat)
}

func TestDeserializeFormatChunk_Corrupted(t *testing.T) {
	payload := []byte{0x00, 0x01, 0x02}
	_, err := DeserializeFormatChunk(payload)
	require.ErrorIs(t, err, ErrFmtChunkCorruptedPayload)
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

func TestDeserializeFactChunk_Normal(t *testing.T) {
	payload := uint32ToBytes(42)
	data, err := DeserializeFactChunk(payload)
	require.NoError(t, err)
	require.Equal(t, uint32(42), data.SampleFrames)
}

func TestDeserializeFactChunk_Corrupted(t *testing.T) {
	payload := []byte{0x00, 0x01, 0x02}
	_, err := DeserializeFactChunk(payload)
	require.ErrorIs(t, err, ErrFactChunkCorruptedPayload)
}

// ------------------------------------------------------------------------- //
// Cue Chunk Data
// ------------------------------------------------------------------------- //

func TestDeserializeCueChunk_Normal(t *testing.T) {
	var payload bytes.Buffer
	payload.Write(uint32ToBytes(2)) // Number of cue points

	// Cue point 0
	payload.Write(uint32ToBytes(1))  // ID
	payload.Write(uint32ToBytes(0))  // Position
	payload.Write(DataChunkID[:])    // FCC Chunk ID (typically 'data')
	payload.Write(uint32ToBytes(0))  // Chunk start (normally 0)
	payload.Write(uint32ToBytes(42)) // Block start
	payload.Write(uint32ToBytes(14)) // Sample offset

	// Cue point 1
	payload.Write(uint32ToBytes(2))  // ID
	payload.Write(uint32ToBytes(0))  // Position
	payload.Write(DataChunkID[:])    // FCC Chunk ID (typically 'data')
	payload.Write(uint32ToBytes(0))  // Chunk start (normally 0)
	payload.Write(uint32ToBytes(88)) // Block start
	payload.Write(uint32ToBytes(64)) // Sample offset

	cueChunkData, err := DeserializeCueChunk(payload.Bytes())
	require.NoError(t, err)

	expected := CueChunkData{
		CuePoints: []CuePoint{
			{
				ID:           1,
				Position:     0,
				FCCChunk:     DataChunkID,
				ChunkStart:   0,
				BlockStart:   42,
				SampleOffset: 14,
			},
			{
				ID:           2,
				Position:     0,
				FCCChunk:     DataChunkID,
				ChunkStart:   0,
				BlockStart:   88,
				SampleOffset: 64,
			},
		},
	}
	require.Equal(t, expected, *cueChunkData)
}

func TestDeserializeCueChunk_Corrupted(t *testing.T) {

	// Unable to read the cue point count
	payload := []byte{0x00, 0x01, 0x02}
	_, err := DeserializeCueChunk(payload)
	require.ErrorIs(t, err, ErrCueChunkCorruptedPayload)

	// Missing cue point data
	payload = uint32ToBytes(2)
	_, err = DeserializeCueChunk(payload)
	require.ErrorIs(t, err, ErrCueChunkCorruptedPayload)

	// Corrupted cue point data
	payload = []byte{
		0x02, 0x00, 0x00, 0x00, // Number of cue points
		0x01, 0x00, 0x00, 0x00, // Start of a valid cue point
		0x00, 0x00, 0x00, 0x00,
		0x64, 0x61, 0x74, 0x61, // ... but truncated
	}
	_, err = DeserializeCueChunk(payload)
	require.ErrorIs(t, err, ErrCueChunkCorruptedPayload)
}
