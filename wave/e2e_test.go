package wave

import (
	ioBytes "bytes"
	"encoding/binary"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
	"time"

	"github.com/jonchammer/audio-io/bytes"
)

// ------------------------------------------------------------------------- //
// End-to-end tests - These are used to ensure the writer consistently
// generates the correct .wav files and that the reader is capable of
// interpreting them.
// ------------------------------------------------------------------------- //

// ------------------------------------------------------------------------- //
// Misc
// ------------------------------------------------------------------------- //

func TestE2E_Empty(t *testing.T) {

	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeUint8, 44100, WithChannelCount(2),
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

	r := NewReader(ioBytes.NewReader(data))

	// Check header
	header, err := r.Header()
	require.NoError(t, err)

	// Verify header fields have proper values
	require.Equal(t, uint32(44), header.ReportedFileSizeBytes)
	require.Nil(t, header.FactData)
	require.Nil(t, header.CueData)
	require.Equal(t, uint32(0), header.DataBytes)
	require.Empty(t, header.AdditionalChunks)

	// Check the format chunk
	require.Equal(t, FormatCodePCM, header.FormatData.FormatCode)
	require.Equal(t, uint16(2), header.FormatData.ChannelCount)
	require.Equal(t, uint32(44100), header.FormatData.FrameRate)
	require.Equal(t, uint32(88200), header.FormatData.ByteRate)
	require.Equal(t, uint16(2), header.FormatData.BlockAlign)
	require.Equal(t, uint16(8), header.FormatData.BitsPerSample)
	require.Nil(t, header.FormatData.ValidBitsPerSample)
	require.Nil(t, header.FormatData.ChannelMask)
	require.Nil(t, header.FormatData.SubFormat)

	// Check header helpers
	require.NoError(t, header.Validate())
	st, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, SampleTypeUint8, st)
	require.Equal(t, uint32(44100), header.FrameRate())
	require.Equal(t, uint32(88200), header.ByteRate())
	require.Equal(t, uint64(88200*8), header.BitRate())
	require.Equal(t, uint16(2), header.ChannelCount())
	require.Equal(t, uint32(0), header.FrameCount())
	require.Equal(t, uint32(0), header.SampleCount())
	require.Equal(t, time.Duration(0), header.PlayTime())

	// Read the audio data. We expect to get an EOF, since there is no data to read.
	buffer := make([]uint8, header.SampleCount())
	n, err := r.ReadUint8(buffer)
	require.ErrorIs(t, err, io.EOF)
	require.Equal(t, 0, n)
}

// ------------------------------------------------------------------------- //
// Uint8
// ------------------------------------------------------------------------- //

func TestE2E_Uint8_Normal(t *testing.T) {

	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeUint8, 44100, WithChannelCount(2),
	)
	require.NoError(t, err)

	// Write the file
	err = w.WriteUint8([]uint8{0, 1, 127, 128, 254, 255})
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

	r := NewReader(ioBytes.NewReader(data))

	// Check header
	header, err := r.Header()
	require.NoError(t, err)

	// Verify header fields have proper values
	require.Equal(t, uint32(50), header.ReportedFileSizeBytes)
	require.Nil(t, header.FactData)
	require.Nil(t, header.CueData)
	require.Equal(t, uint32(6), header.DataBytes)
	require.Empty(t, header.AdditionalChunks)

	// Check the format chunk
	require.Equal(t, FormatCodePCM, header.FormatData.FormatCode)
	require.Equal(t, uint16(2), header.FormatData.ChannelCount)
	require.Equal(t, uint32(44100), header.FormatData.FrameRate)
	require.Equal(t, uint32(88200), header.FormatData.ByteRate)
	require.Equal(t, uint16(2), header.FormatData.BlockAlign)
	require.Equal(t, uint16(8), header.FormatData.BitsPerSample)
	require.Nil(t, header.FormatData.ValidBitsPerSample)
	require.Nil(t, header.FormatData.ChannelMask)
	require.Nil(t, header.FormatData.SubFormat)

	// Check header helpers
	require.NoError(t, header.Validate())
	st, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, SampleTypeUint8, st)
	require.Equal(t, uint32(44100), header.FrameRate())
	require.Equal(t, uint32(88200), header.ByteRate())
	require.Equal(t, uint64(88200*8), header.BitRate())
	require.Equal(t, uint16(2), header.ChannelCount())
	require.Equal(t, uint32(3), header.FrameCount())
	require.Equal(t, uint32(6), header.SampleCount())

	seconds := 3.0 / 44100.0
	require.Equal(t, time.Duration(seconds*1e9), header.PlayTime())

	// Read the audio data.
	buffer := make([]uint8, header.SampleCount())
	n, err := r.ReadUint8(buffer)
	require.NoError(t, err)
	require.Equal(t, []byte{0, 1, 127, 128, 254, 255}, buffer[:n])
}

func TestE2E_Uint8_Padding(t *testing.T) {

	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeUint8, 44100,
	)
	require.NoError(t, err)

	// Write the file
	err = w.WriteUint8([]uint8{0, 1, 2})
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

	r := NewReader(ioBytes.NewReader(data))

	// Check header
	header, err := r.Header()
	require.NoError(t, err)

	// Verify header fields have proper values
	require.Equal(t, uint32(48), header.ReportedFileSizeBytes)
	require.Nil(t, header.FactData)
	require.Nil(t, header.CueData)
	require.Equal(t, uint32(3), header.DataBytes)
	require.Empty(t, header.AdditionalChunks)

	// Check the format chunk
	require.Equal(t, FormatCodePCM, header.FormatData.FormatCode)
	require.Equal(t, uint16(1), header.FormatData.ChannelCount)
	require.Equal(t, uint32(44100), header.FormatData.FrameRate)
	require.Equal(t, uint32(44100), header.FormatData.ByteRate)
	require.Equal(t, uint16(1), header.FormatData.BlockAlign)
	require.Equal(t, uint16(8), header.FormatData.BitsPerSample)
	require.Nil(t, header.FormatData.ValidBitsPerSample)
	require.Nil(t, header.FormatData.ChannelMask)
	require.Nil(t, header.FormatData.SubFormat)

	// Check header helpers
	require.NoError(t, header.Validate())
	st, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, SampleTypeUint8, st)
	require.Equal(t, uint32(44100), header.FrameRate())
	require.Equal(t, uint32(44100), header.ByteRate())
	require.Equal(t, uint64(44100*8), header.BitRate())
	require.Equal(t, uint16(1), header.ChannelCount())
	require.Equal(t, uint32(3), header.FrameCount())
	require.Equal(t, uint32(3), header.SampleCount())

	seconds := 3.0 / 44100.0
	require.Equal(t, time.Duration(seconds*1e9), header.PlayTime())

	// Read the audio data.
	buffer := make([]uint8, header.SampleCount())
	n, err := r.ReadUint8(buffer)
	require.NoError(t, err)
	require.Equal(t, []byte{0, 1, 2}, buffer[:n])
}

func TestE2E_Uint8_Extensible(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeUint8, 44100, WithChannelCount(4),
	)
	require.NoError(t, err)

	err = w.WriteUint8(
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

	r := NewReader(ioBytes.NewReader(data))

	// Check header
	header, err := r.Header()
	require.NoError(t, err)

	// Verify header fields have proper values
	require.Equal(t, uint32(92), header.ReportedFileSizeBytes)
	require.NotNil(t, header.FactData)
	require.Nil(t, header.CueData)
	require.Equal(t, uint32(12), header.DataBytes)
	require.Empty(t, header.AdditionalChunks)

	// Check the format chunk
	require.Equal(t, FormatCodeExtensible, header.FormatData.FormatCode)
	require.Equal(t, uint16(4), header.FormatData.ChannelCount)
	require.Equal(t, uint32(44100), header.FormatData.FrameRate)
	require.Equal(t, uint32(176400), header.FormatData.ByteRate)
	require.Equal(t, uint16(4), header.FormatData.BlockAlign)
	require.Equal(t, uint16(8), header.FormatData.BitsPerSample)
	require.NotNil(t, header.FormatData.ValidBitsPerSample)
	require.Equal(t, uint16(8), *header.FormatData.ValidBitsPerSample)
	require.NotNil(t, header.FormatData.ChannelMask)
	require.Equal(t, uint32(0), *header.FormatData.ChannelMask)
	require.NotNil(t, header.FormatData.SubFormat)
	require.Equal(t, FormatCodePCM, *header.FormatData.SubFormat)

	// Check the fact chunk
	require.Equal(t, uint32(3), header.FactData.SampleFrames)

	// Check header helpers
	require.NoError(t, header.Validate())
	st, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, SampleTypeUint8, st)
	require.Equal(t, uint32(44100), header.FrameRate())
	require.Equal(t, uint32(176400), header.ByteRate())
	require.Equal(t, uint64(176400*8), header.BitRate())
	require.Equal(t, uint16(4), header.ChannelCount())
	require.Equal(t, uint32(3), header.FrameCount())
	require.Equal(t, uint32(12), header.SampleCount())

	seconds := 3.0 / 44100.0
	require.Equal(t, time.Duration(seconds*1e9), header.PlayTime())

	// Read the audio data.
	buffer := make([]uint8, header.SampleCount())
	n, err := r.ReadUint8(buffer)
	require.NoError(t, err)
	require.Equal(
		t,
		[]uint8{0, 1, 2, 3, 127, 128, 129, 130, 252, 253, 254, 255},
		buffer[:n],
	)
}

// ------------------------------------------------------------------------- //
// Int16
// ------------------------------------------------------------------------- //

func TestE2E_Int16_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeInt16, 44100, WithChannelCount(2),
	)
	require.NoError(t, err)

	err = w.WriteInt16([]int16{-32768, -32767, 0, 1, 32766, 32767})
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

	r := NewReader(ioBytes.NewReader(data))

	// Check header
	header, err := r.Header()
	require.NoError(t, err)

	// Verify header fields have proper values
	require.Equal(t, uint32(56), header.ReportedFileSizeBytes)
	require.Nil(t, header.FactData)
	require.Nil(t, header.CueData)
	require.Equal(t, uint32(12), header.DataBytes)
	require.Empty(t, header.AdditionalChunks)

	// Check the format chunk
	require.Equal(t, FormatCodePCM, header.FormatData.FormatCode)
	require.Equal(t, uint16(2), header.FormatData.ChannelCount)
	require.Equal(t, uint32(44100), header.FormatData.FrameRate)
	require.Equal(t, uint32(176400), header.FormatData.ByteRate)
	require.Equal(t, uint16(4), header.FormatData.BlockAlign)
	require.Equal(t, uint16(16), header.FormatData.BitsPerSample)
	require.Nil(t, header.FormatData.ValidBitsPerSample)
	require.Nil(t, header.FormatData.ChannelMask)
	require.Nil(t, header.FormatData.SubFormat)

	// Check header helpers
	require.NoError(t, header.Validate())
	st, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, SampleTypeInt16, st)
	require.Equal(t, uint32(44100), header.FrameRate())
	require.Equal(t, uint32(176400), header.ByteRate())
	require.Equal(t, uint64(176400*8), header.BitRate())
	require.Equal(t, uint16(2), header.ChannelCount())
	require.Equal(t, uint32(3), header.FrameCount())
	require.Equal(t, uint32(6), header.SampleCount())

	seconds := 3.0 / 44100.0
	require.Equal(t, time.Duration(seconds*1e9), header.PlayTime())

	// Read the audio data.
	buffer := make([]int16, header.SampleCount())
	n, err := r.ReadInt16(buffer)
	require.NoError(t, err)
	require.Equal(t, []int16{-32768, -32767, 0, 1, 32766, 32767}, buffer[:n])
}

func TestE2E_Int16_Extensible(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeInt16, 44100, WithChannelCount(4),
	)
	require.NoError(t, err)

	err = w.WriteInt16([]int16{
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

	r := NewReader(ioBytes.NewReader(data))

	// Check header
	header, err := r.Header()
	require.NoError(t, err)

	// Verify header fields have proper values
	require.Equal(t, uint32(104), header.ReportedFileSizeBytes)
	require.NotNil(t, header.FactData)
	require.Nil(t, header.CueData)
	require.Equal(t, uint32(24), header.DataBytes)
	require.Empty(t, header.AdditionalChunks)

	// Check the format chunk
	require.Equal(t, FormatCodeExtensible, header.FormatData.FormatCode)
	require.Equal(t, uint16(4), header.FormatData.ChannelCount)
	require.Equal(t, uint32(44100), header.FormatData.FrameRate)
	require.Equal(t, uint32(352800), header.FormatData.ByteRate)
	require.Equal(t, uint16(8), header.FormatData.BlockAlign)
	require.Equal(t, uint16(16), header.FormatData.BitsPerSample)
	require.NotNil(t, header.FormatData.ValidBitsPerSample)
	require.Equal(t, uint16(16), *header.FormatData.ValidBitsPerSample)
	require.NotNil(t, header.FormatData.ChannelMask)
	require.Equal(t, uint32(0), *header.FormatData.ChannelMask)
	require.NotNil(t, header.FormatData.SubFormat)
	require.Equal(t, FormatCodePCM, *header.FormatData.SubFormat)

	// Check the fact chunk
	require.Equal(t, uint32(3), header.FactData.SampleFrames)

	// Check header helpers
	require.NoError(t, header.Validate())
	st, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, SampleTypeInt16, st)
	require.Equal(t, uint32(44100), header.FrameRate())
	require.Equal(t, uint32(352800), header.ByteRate())
	require.Equal(t, uint64(352800*8), header.BitRate())
	require.Equal(t, uint16(4), header.ChannelCount())
	require.Equal(t, uint32(3), header.FrameCount())
	require.Equal(t, uint32(12), header.SampleCount())

	seconds := 3.0 / 44100.0
	require.Equal(t, time.Duration(seconds*1e9), header.PlayTime())

	// Read the audio data.
	buffer := make([]int16, header.SampleCount())
	n, err := r.ReadInt16(buffer)
	require.NoError(t, err)
	require.Equal(
		t,
		[]int16{-32768, -32767, -32766, -32765, 0, 1, 2, 3, 32764, 32765, 32766, 32767},
		buffer[:n],
	)
}

// ------------------------------------------------------------------------- //
// Int24
// ------------------------------------------------------------------------- //

func TestE2E_Int24_Normal(t *testing.T) {

	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeInt24, 44100, WithChannelCount(2),
	)
	require.NoError(t, err)

	err = w.WriteInt24([]int32{-8388608, -8388607, 0, 1, 8388606, 8388607})
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

	r := NewReader(ioBytes.NewReader(data))

	// Check header
	header, err := r.Header()
	require.NoError(t, err)

	// Verify header fields have proper values
	require.Equal(t, uint32(98), header.ReportedFileSizeBytes)
	require.NotNil(t, header.FactData)
	require.Nil(t, header.CueData)
	require.Equal(t, uint32(18), header.DataBytes)
	require.Empty(t, header.AdditionalChunks)

	// Check the format chunk
	require.Equal(t, FormatCodeExtensible, header.FormatData.FormatCode)
	require.Equal(t, uint16(2), header.FormatData.ChannelCount)
	require.Equal(t, uint32(44100), header.FormatData.FrameRate)
	require.Equal(t, uint32(264600), header.FormatData.ByteRate)
	require.Equal(t, uint16(6), header.FormatData.BlockAlign)
	require.Equal(t, uint16(24), header.FormatData.BitsPerSample)
	require.NotNil(t, header.FormatData.ValidBitsPerSample)
	require.Equal(t, uint16(24), *header.FormatData.ValidBitsPerSample)
	require.NotNil(t, header.FormatData.ChannelMask)
	require.Equal(t, uint32(0), *header.FormatData.ChannelMask)
	require.NotNil(t, header.FormatData.SubFormat)
	require.Equal(t, FormatCodePCM, *header.FormatData.SubFormat)

	// Check the fact chunk
	require.Equal(t, uint32(3), header.FactData.SampleFrames)

	// Check header helpers
	require.NoError(t, header.Validate())
	st, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, SampleTypeInt24, st)
	require.Equal(t, uint32(44100), header.FrameRate())
	require.Equal(t, uint32(264600), header.ByteRate())
	require.Equal(t, uint64(264600*8), header.BitRate())
	require.Equal(t, uint16(2), header.ChannelCount())
	require.Equal(t, uint32(3), header.FrameCount())
	require.Equal(t, uint32(6), header.SampleCount())

	seconds := 3.0 / 44100.0
	require.Equal(t, time.Duration(seconds*1e9), header.PlayTime())

	// Read the audio data.
	buffer := make([]int32, header.SampleCount())
	n, err := r.ReadInt24(buffer)
	require.NoError(t, err)
	require.Equal(
		t,
		[]int32{-8388608, -8388607, 0, 1, 8388606, 8388607},
		buffer[:n],
	)
}

func TestE2E_Int24_Padding(t *testing.T) {

	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeInt24, 44100,
	)
	require.NoError(t, err)

	err = w.WriteInt24([]int32{-8388608})
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

	r := NewReader(ioBytes.NewReader(data))

	// Check header
	header, err := r.Header()
	require.NoError(t, err)

	// Verify header fields have proper values
	require.Equal(t, uint32(84), header.ReportedFileSizeBytes)
	require.NotNil(t, header.FactData)
	require.Nil(t, header.CueData)
	require.Equal(t, uint32(3), header.DataBytes)
	require.Empty(t, header.AdditionalChunks)

	// Check the format chunk
	require.Equal(t, FormatCodeExtensible, header.FormatData.FormatCode)
	require.Equal(t, uint16(1), header.FormatData.ChannelCount)
	require.Equal(t, uint32(44100), header.FormatData.FrameRate)
	require.Equal(t, uint32(132300), header.FormatData.ByteRate)
	require.Equal(t, uint16(3), header.FormatData.BlockAlign)
	require.Equal(t, uint16(24), header.FormatData.BitsPerSample)
	require.NotNil(t, header.FormatData.ValidBitsPerSample)
	require.Equal(t, uint16(24), *header.FormatData.ValidBitsPerSample)
	require.NotNil(t, header.FormatData.ChannelMask)
	require.Equal(t, uint32(0), *header.FormatData.ChannelMask)
	require.NotNil(t, header.FormatData.SubFormat)
	require.Equal(t, FormatCodePCM, *header.FormatData.SubFormat)

	// Check the fact chunk
	require.Equal(t, uint32(1), header.FactData.SampleFrames)

	// Check header helpers
	require.NoError(t, header.Validate())
	st, err := header.SampleType()
	require.NoError(t, err)
	require.Equal(t, SampleTypeInt24, st)
	require.Equal(t, uint32(44100), header.FrameRate())
	require.Equal(t, uint32(132300), header.ByteRate())
	require.Equal(t, uint64(132300*8), header.BitRate())
	require.Equal(t, uint16(1), header.ChannelCount())
	require.Equal(t, uint32(1), header.FrameCount())
	require.Equal(t, uint32(1), header.SampleCount())

	seconds := 1.0 / 44100.0
	require.Equal(t, time.Duration(seconds*1e9), header.PlayTime())

	// Read the audio data.
	buffer := make([]int32, header.SampleCount())
	n, err := r.ReadInt24(buffer)
	require.NoError(t, err)
	require.Equal(
		t,
		[]int32{-8388608},
		buffer[:n],
	)
}

// ------------------------------------------------------------------------- //
// Int32
// ------------------------------------------------------------------------- //

func TestE2E_Int32_Normal(t *testing.T) {
	baseWriter := &bytes.Writer{}
	w, err := NewWriter(
		baseWriter, SampleTypeInt32, 44100, WithChannelCount(2),
	)
	require.NoError(t, err)

	err = w.WriteInt32([]int32{
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
		baseWriter, SampleTypeFloat32, 44100, WithChannelCount(2),
	)
	require.NoError(t, err)

	err = w.WriteFloat32([]float32{
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
		baseWriter, SampleTypeFloat32, 44100, WithChannelCount(4),
	)
	require.NoError(t, err)

	err = w.WriteFloat32([]float32{
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
		baseWriter, SampleTypeFloat64, 44100, WithChannelCount(2),
	)
	require.NoError(t, err)

	err = w.WriteFloat64([]float64{
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
		baseWriter, SampleTypeFloat64, 44100, WithChannelCount(4),
	)
	require.NoError(t, err)

	err = w.WriteFloat64([]float64{
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
