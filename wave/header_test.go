package wave

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/jonchammer/audio-io/core"
)

func TestParseHeaderFromRIFFChunk_Normal(t *testing.T) {

	totalFileSize := uint32(42)
	riffChunkData := &RIFFChunkData{
		SubChunks: []Chunk{
			{
				ID:   FormatChunkID,
				Size: 16,
				Body: []byte{
					0x01, 0x00, 0x02, 0x00,
					0x44, 0xAC, 0x00, 0x00,
					0x10, 0xB1, 0x02, 0x00,
					0x04, 0x00, 0x10, 0x00,
				},
			},
			{
				ID:   FactChunkID,
				Size: 4,
				Body: []byte{0x80, 0x00, 0x00, 0x00},
			},
			{
				ID:   CueChunkID,
				Size: 28,
				Body: []byte{
					0x01, 0x00, 0x00, 0x00,
					0x01, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x64, 0x61, 0x74, 0x61,
					0x00, 0x00, 0x00, 0x00,
					0x2A, 0x00, 0x00, 0x00,
					0x0E, 0x00, 0x00, 0x00,
				},
			},
			{
				ID:   [4]byte{'z', 'e', 'r', 'o'},
				Size: 8,
				Body: []byte{
					0x00, 0x01, 0x02, 0x03,
					0x04, 0x05, 0x06, 0x07,
				},
			},
			{
				ID:   DataChunkID,
				Size: 8,
				Body: []byte{
					0x10, 0x11, 0x12, 0x13,
					0x14, 0x15, 0x16, 0x17,
				},
			},
		},
	}

	header, err := parseHeaderFromRIFFChunk(totalFileSize, riffChunkData)
	require.NoError(t, err)

	require.Equal(t, uint32(42), totalFileSize)

	// Check format section
	expectedFmt := NewFormatChunkData(2, 44100, core.SampleTypeInt16)
	require.Equal(t, expectedFmt, header.FormatData)

	// Check fact section
	expectedFact := FactChunkData{
		SampleFrames: 128,
	}
	require.NotNil(t, header.FactData)
	require.Equal(t, expectedFact, *header.FactData)

	// Check cue section
	expectedCue := CueChunkData{
		CuePoints: []CuePoint{
			{
				ID:           1,
				Position:     0,
				FCCChunk:     DataChunkID,
				ChunkStart:   0,
				BlockStart:   42,
				SampleOffset: 14,
			},
		},
	}
	require.NotNil(t, header.CueData)
	require.Equal(t, expectedCue, *header.CueData)

	// Check data section
	require.Equal(t, uint32(8), header.DataBytes)

	// Check additional chunks
	require.Equal(t, []Chunk{
		{
			ID:   [4]byte{'z', 'e', 'r', 'o'},
			Size: 8,
			Body: []byte{
				0x00, 0x01, 0x02, 0x03,
				0x04, 0x05, 0x06, 0x07,
			},
		},
	}, header.AdditionalChunks)
}

func TestParseHeaderFromRIFFChunk_Corrupted(t *testing.T) {

	// Corrupted format chunk
	riffChunkData := &RIFFChunkData{
		SubChunks: []Chunk{
			{
				ID:   FormatChunkID,
				Size: 16,
				Body: []byte{
					0x01, 0x00, 0x02, 0x00,
					0x44, 0xAC, 0x00, 0x00,
				},
			},
		},
	}
	_, err := parseHeaderFromRIFFChunk(42, riffChunkData)
	require.ErrorIs(t, err, ErrFmtChunkCorruptedPayload)

	// Corrupted fact chunk
	riffChunkData = &RIFFChunkData{
		SubChunks: []Chunk{
			{
				ID:   FactChunkID,
				Size: 4,
				Body: []byte{0x80, 0x00},
			},
		},
	}
	_, err = parseHeaderFromRIFFChunk(42, riffChunkData)
	require.ErrorIs(t, err, ErrFactChunkCorruptedPayload)

	// Corrupted cue chunk
	riffChunkData = &RIFFChunkData{
		SubChunks: []Chunk{
			{
				ID:   CueChunkID,
				Size: 28,
				Body: []byte{
					0x01, 0x00, 0x00, 0x00,
				},
			},
		},
	}
	_, err = parseHeaderFromRIFFChunk(42, riffChunkData)
	require.ErrorIs(t, err, ErrCueChunkCorruptedPayload)
}

func TestParseHeaderFromRIFFChunk_MissingFmt(t *testing.T) {
	riffChunkData := &RIFFChunkData{
		SubChunks: []Chunk{
			{
				ID:   DataChunkID,
				Size: 8,
				Body: []byte{
					0x10, 0x11, 0x12, 0x13,
					0x14, 0x15, 0x16, 0x17,
				},
			},
		},
	}
	_, err := parseHeaderFromRIFFChunk(42, riffChunkData)
	require.ErrorIs(t, err, ErrHeaderMissingFmtChunk)
}

func TestHeader_Validate_InvalidFormatCode(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.FormatCode = FormatCode(99)
	header := getValidHeader(formatData)
	err := header.Validate()
	require.ErrorContains(t, err, "format code: '99' was not recognized")
}

func TestHeader_Validate_InvalidByteRate(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.ByteRate = 100
	header := getValidHeader(formatData)
	err := header.Validate()
	require.ErrorContains(t, err, "byte rate: '100' did not match expected result: '44100'")
}

func TestHeader_Validate_InvalidBlockAlign(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.BlockAlign = 100
	header := getValidHeader(formatData)
	err := header.Validate()
	require.ErrorContains(t, err, "block align: '100' did not match expected result: '1'")
}

func TestHeader_Validate_InvalidSampleFrames(t *testing.T) {
	header := getValidHeader(getValidFormatChunkData())
	header.FactData = &FactChunkData{
		SampleFrames: 100, // Should be 42
	}
	err := header.Validate()
	require.ErrorContains(t, err, "sample frames: '100' did not match expected result: '42'")
}

func TestHeader_Validate_InvalidExtensions(t *testing.T) {

	// Valid bits per sample
	formatData := getValidFormatChunkData()
	bps := uint16(8)
	formatData.ValidBitsPerSample = &bps
	header := getValidHeader(formatData)
	err := header.Validate()
	require.ErrorContains(t, err, "valid bits per sample should only be set if format code is extensible")

	// Channel mask
	formatData = getValidFormatChunkData()
	cm := uint32(0)
	formatData.ChannelMask = &cm
	header = getValidHeader(formatData)
	err = header.Validate()
	require.ErrorContains(t, err, "channel mask should only be set if format code is extensible")

	// Sub format
	formatData = getValidFormatChunkData()
	sub := FormatCodePCM
	formatData.SubFormat = &sub
	header = getValidHeader(formatData)
	err = header.Validate()
	require.ErrorContains(t, err, "sub format should only be set if format code is extensible")
}

func TestHeader_SampleType_Uint8(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.FormatCode = FormatCodePCM
	formatData.BitsPerSample = 8
	header := getValidHeader(formatData)

	sampleType, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, core.SampleTypeUint8, sampleType)
}

func TestHeader_SampleType_Int16(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.FormatCode = FormatCodePCM
	formatData.BitsPerSample = 16
	header := getValidHeader(formatData)

	sampleType, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, core.SampleTypeInt16, sampleType)
}

func TestHeader_SampleType_Int24(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.FormatCode = FormatCodePCM
	formatData.BitsPerSample = 24
	header := getValidHeader(formatData)

	sampleType, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, core.SampleTypeInt24, sampleType)
}

func TestHeader_SampleType_Int32(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.FormatCode = FormatCodePCM
	formatData.BitsPerSample = 32
	header := getValidHeader(formatData)

	sampleType, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, core.SampleTypeInt32, sampleType)
}

func TestHeader_SampleType_Float32(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.FormatCode = FormatCodeIEEEFloat
	formatData.BitsPerSample = 32
	header := getValidHeader(formatData)

	sampleType, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, core.SampleTypeFloat32, sampleType)
}

func TestHeader_SampleType_Float64(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.FormatCode = FormatCodeIEEEFloat
	formatData.BitsPerSample = 64
	header := getValidHeader(formatData)

	sampleType, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, core.SampleTypeFloat64, sampleType)
}

func TestHeader_SampleType_InvalidFormatCode(t *testing.T) {

	// Missing sub format
	formatData := getValidFormatChunkData()
	formatData.FormatCode = FormatCodeExtensible
	formatData.SubFormat = nil
	header := getValidHeader(formatData)
	_, err := header.SampleType()
	require.ErrorIs(t, err, ErrFmtChunkMissingSubFormat)

	// Format Code itself is invalid
	formatData = getValidFormatChunkData()
	formatData.FormatCode = FormatCode(99)
	formatData.BitsPerSample = 8
	header = getValidHeader(formatData)
	_, err = header.SampleType()
	require.ErrorContains(t, err, "invalid format code: 'FormatCode(99)'")
}

func TestHeader_SampleType_InvalidPCMBitsPerSample(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.FormatCode = FormatCodePCM
	formatData.BitsPerSample = 12
	header := getValidHeader(formatData)

	_, err := header.SampleType()
	require.ErrorContains(t, err, "unknown PCM type: '12' bits per sample")
}

func TestHeader_SampleType_InvalidIEEEFloatBitsPerSample(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.FormatCode = FormatCodeIEEEFloat
	formatData.BitsPerSample = 12
	header := getValidHeader(formatData)

	_, err := header.SampleType()
	require.ErrorContains(t, err, "unknown IEEE float type: '12' bits per sample")
}

// ------------------------------------------------------------------------- //
// Helpers
// ------------------------------------------------------------------------- //

func getValidFormatChunkData() FormatChunkData {
	return FormatChunkData{
		FormatCode:         FormatCodePCM,
		ChannelCount:       1,
		FrameRate:          44100,
		ByteRate:           44100,
		BlockAlign:         1,
		BitsPerSample:      8,
		ValidBitsPerSample: nil,
		ChannelMask:        nil,
		SubFormat:          nil,
	}
}

func getValidHeader(formatData FormatChunkData) *Header {
	return &Header{
		ReportedFileSizeBytes: 100,
		FormatData:            formatData,
		FactData:              nil,
		CueData:               nil,
		DataBytes:             42,
		AdditionalChunks:      nil,
	}
}
