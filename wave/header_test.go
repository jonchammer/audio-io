package wave

import (
	"github.com/stretchr/testify/require"
	"testing"
)

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
	require.Equal(t, SampleTypeUint8, sampleType)
}

func TestHeader_SampleType_Int16(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.FormatCode = FormatCodePCM
	formatData.BitsPerSample = 16
	header := getValidHeader(formatData)

	sampleType, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, SampleTypeInt16, sampleType)
}

func TestHeader_SampleType_Int24(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.FormatCode = FormatCodePCM
	formatData.BitsPerSample = 24
	header := getValidHeader(formatData)

	sampleType, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, SampleTypeInt24, sampleType)
}

func TestHeader_SampleType_Int32(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.FormatCode = FormatCodePCM
	formatData.BitsPerSample = 32
	header := getValidHeader(formatData)

	sampleType, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, SampleTypeInt32, sampleType)
}

func TestHeader_SampleType_Float32(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.FormatCode = FormatCodeIEEEFloat
	formatData.BitsPerSample = 32
	header := getValidHeader(formatData)

	sampleType, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, SampleTypeFloat32, sampleType)
}

func TestHeader_SampleType_Float64(t *testing.T) {
	formatData := getValidFormatChunkData()
	formatData.FormatCode = FormatCodeIEEEFloat
	formatData.BitsPerSample = 64
	header := getValidHeader(formatData)

	sampleType, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, SampleTypeFloat64, sampleType)
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
