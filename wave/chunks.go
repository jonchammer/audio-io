package wave

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// ------------------------------------------------------------------------- //
// Chunk
// ------------------------------------------------------------------------- //

// A Chunk is the core unit of the .wav file spec. Each chunk has an 8-byte
// header and a variable length body. The first 4 bytes represent the chunk's
// identifier, and the next 4 bytes contain the size of the body (measured
// in bytes).
//
// Most chunks are simple, in that they only contain data, but chunks are
// permitted to act as containers for other chunks. The main 'RIFF' chunk, for
// example, contains all other chunks, including metadata and audio data.
type Chunk struct {
	ID   [4]byte
	Size uint32
	Body []byte
}

// Serialize transforms this chunk into a []byte according to the .wav
// specification.
func (c Chunk) Serialize() []byte {

	result := make([]byte, 0, 4+4+len(c.Body))
	result = append(result, c.ID[:]...)
	result = append(result, uint32ToBytes(c.Size)...)
	result = append(result, c.Body...)

	return result
}

func (c Chunk) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(c.Serialize())
	return int64(n), err
}

// ------------------------------------------------------------------------- //
// RIFF chunk
// ------------------------------------------------------------------------- //

var (
	RIFFChunkID = [4]byte{'R', 'I', 'F', 'F'}
	WaveID      = [4]byte{'W', 'A', 'V', 'E'}

	ErrRIFFChunkCorruptedHeader = errors.New("RIFF header is corrupted")
)

// NewRIFFChunk returns a 'RIFF' Chunk containing the given RIFFChunkData. The
// 'RIFF' chunk will be the root element of the chunk tree; it contains all
// other chunks.
func NewRIFFChunk(data *RIFFChunkData) Chunk {
	riffData, totalSizeBytes := data.Serialize()
	return Chunk{
		ID:   RIFFChunkID,
		Size: totalSizeBytes,
		Body: riffData,
	}
}

type RIFFChunkData struct {

	// The RIFF chunk includes the WAVE ID before the sub chunks begin, but as
	// this field is always set to the constant "WAVE", we don't include it in
	// the actual RIFFChunkData structure.
	// WaveID [4]byte

	SubChunks []Chunk
}

// Serialize returns 1) the serialized representation of the RIFF Chunk and 2)
// the total 'expected' size of the RIFF Chunk.
//
// NOTE: The size of the resulting []byte is NOT guaranteed to match the
// expected size. The expected size is calculated using the Chunk.Size field,
// which is allowed to be non-zero, even if Chunk.Body is empty. This is used
// internally to deal with the 'data' chunk. The preamble of the wav file will
// be written multiple times, but the data will only be written once.
func (d RIFFChunkData) Serialize() ([]byte, uint32) {

	// We use a reasonable buffer size based on the max size of a WAV preamble
	// with no audio data:
	//
	//    4 - "WAVE"
	//   48 - Format Chunk (using extension format)
	//   12 - Fact Chunk
	// +  8 - Data Chunk Header
	//   -- - -------------------
	//   72 - Total
	const maxExpectedSizeBytes = 72

	riffBody := make([]byte, 0, maxExpectedSizeBytes)
	riffBody = append(riffBody, WaveID[:]...)

	totalSizeBytes := uint32(len(riffBody))
	for _, chunk := range d.SubChunks {
		riffBody = append(riffBody, chunk.Serialize()...)

		// The total size includes 8 bytes for the chunk's header, the actual
		// chunk size as reported by the chunk itself, and a padding byte, to
		// be included if the chunk size is odd.
		totalSizeBytes += 8 + chunk.Size + (chunk.Size & 1)
	}

	return riffBody, totalSizeBytes
}

// ReadRIFFChunk reads a RIFF chunk from the given reader, returning the total
// file size (including the 8 bytes of header information) and a RIFFChunkData
// structure upon success.
//
// ReadRIFFChunk will scan through the entire reader, searching for any chunks
// within the file. After extracting all relevant metadata, the reader will be
// reset to the beginning of the 'data' chunk, ready for buffered reads.
func ReadRIFFChunk(r io.ReadSeeker) (uint32, *RIFFChunkData, error) {

	buffer := make([]byte, 4)

	// RIFF ID ("RIFF")
	_, err := io.ReadFull(r, buffer)
	if err != nil {
		return 0, nil, err
	}
	if !bytes.Equal(buffer, RIFFChunkID[:]) {
		return 0, nil, ErrRIFFChunkCorruptedHeader
	}

	// File size
	_, err = io.ReadFull(r, buffer)
	if err != nil {
		return 0, nil, err
	}
	fileSize := binary.LittleEndian.Uint32(buffer) + 8

	// Read the WAVE ID
	_, err = io.ReadFull(r, buffer)
	if err != nil {
		return 0, nil, err
	}
	if !bytes.Equal(buffer, WaveID[:]) {
		return 0, nil, ErrRIFFChunkCorruptedHeader
	}

	currentOffset := int64(12)
	dataChunkOffset := int64(0)

	// Read the sub chunks. For now, we assume that the data chunk will be the
	// last entry in the file, and we'll avoid reading the actual audio data.
	chunks := make([]Chunk, 0, 2)
	for {

		if currentOffset >= int64(fileSize) {
			break
		}

		// Chunk ID
		var chunkID [4]byte
		_, err = io.ReadFull(r, chunkID[:])
		if err != nil {
			return 0, nil, err
		}
		currentOffset += 4

		// Chunk size
		_, err = io.ReadFull(r, buffer)
		if err != nil {
			return 0, nil, err
		}
		chunkSize := binary.LittleEndian.Uint32(buffer)
		paddingByteCount := int64(chunkSize & 1)
		currentOffset += 4

		// Chunk body - For any chunk but the 'data' one, we'll read the chunk
		// body in full. For the 'data' chunk, we'll simply skip over those
		// bytes instead.
		var chunkBytes []byte
		if chunkID != DataChunkID {
			chunkBytes = make([]byte, chunkSize)
			_, err = io.ReadFull(r, chunkBytes)
			if err != nil {
				return 0, nil, err
			}
			currentOffset += int64(chunkSize)

			// If a padding byte is present, we'll need to account for it too.
			if paddingByteCount != 0 {
				_, err = r.Read(buffer[:1])
				if err != nil {
					return 0, nil, err
				}
				currentOffset++
			}

		} else {
			dataChunkOffset = currentOffset
			currentOffset, err = r.Seek(
				currentOffset+int64(chunkSize)+paddingByteCount,
				io.SeekStart,
			)
			if err != nil {
				return 0, nil, err
			}
		}

		chunks = append(chunks, Chunk{
			ID:   chunkID,
			Size: chunkSize,
			Body: chunkBytes,
		})
	}

	// Reset 'r' to the beginning of the data chunk
	_, err = r.Seek(dataChunkOffset, io.SeekStart)
	if err != nil {
		return 0, nil, err
	}

	return fileSize, &RIFFChunkData{
		SubChunks: chunks,
	}, nil
}

// ------------------------------------------------------------------------- //
// Format chunk
// ------------------------------------------------------------------------- //

var (
	FormatChunkID                = [4]byte{'f', 'm', 't', ' '}
	ErrFmtChunkMissingSubFormat  = errors.New("sub format is expected, but not present")
	ErrFmtChunkInvalidExtensible = errors.New("extensible format requires that ValidBitsPerSample, ChannelMask, and SubFormat be set")
	ErrFmtChunkCorruptedPayload  = errors.New("detected corrupted 'fmt' payload")
)

// NewFormatChunk returns a 'fmt' Chunk containing the given FormatChunkData.
// The 'fmt' chunk contains most of the metadata about the audio file,
// including channel count, sample rate, and bits per sample.
func NewFormatChunk(data *FormatChunkData) (Chunk, error) {
	fmtData, err := data.Serialize()
	if err != nil {
		return Chunk{}, err
	}

	return Chunk{
		ID:   FormatChunkID,
		Size: uint32(len(fmtData)),
		Body: fmtData,
	}, nil
}

type FormatChunkData struct {

	// The FormatCode specifies how the audio data should be interpreted.
	FormatCode FormatCode

	// ChannelCount represents the number of channels of audio in the file.
	ChannelCount uint16

	// The FrameRate is a count of the number of frames to be played per
	// second (e.g. 44100), independent of the number of channels. Note that
	// the term "sample rate" is perhaps a bit more common, but we try to avoid
	// using it because of the ambiguity with multi-channel data.
	FrameRate uint32

	// ByteRate is a count of the number of bytes to be played per second. It
	// should be equal to the FrameRate * the number of bytes per sample *
	// ChannelCount.
	ByteRate uint32

	// BlockAlign describes the number of bytes in a single frame of audio. It
	// should be equal to the number of bytes per sample * ChannelCount.
	BlockAlign uint16

	// BitsPerSample represents the number of bits present in each sample,
	// rounded up to the nearest multiple of 8. (e.g. 24-bit PCM data would
	// have a value of 24, while IEEE float32 data would have a value of 32).
	BitsPerSample uint16

	// ValidBitsPerSample should be defined (non-nil) if:
	//   FormatCode == FormatCodeExtensible
	//
	// See the wav specification for exact details.
	ValidBitsPerSample *uint16

	// ChannelMask should be defined (non-nil) if:
	//   FormatCode == FormatCodeExtensible
	//
	// It is typically used to represent the 'speaker position mask'. See the
	// wave specification for exact details.
	ChannelMask *uint32

	// SubFormat should be defined (non-nil) if:
	//   FormatCode == FormatCodeExtensible
	//
	// It will have the same function as FormatCode.
	SubFormat *FormatCode
}

func NewFormatChunkData(
	channelCount uint16,
	frameRate uint32,
	sampleType SampleType,
) FormatChunkData {

	// Consolidate the required information from the sample type
	effectiveFormatCode := sampleType.EffectiveFormatCode()
	sampleSizeBytes := sampleType.Size()

	// Calculate some common values needed for the format chunk
	bitsPerSample := uint16(sampleSizeBytes * 8)
	blockAlign := uint16(sampleSizeBytes) * channelCount
	bytesPerSecond := frameRate * uint32(blockAlign)

	// We'll use the extensible format when it applies
	isExtensible := channelCount > 2 ||
		sampleType == SampleTypeInt24 ||
		sampleType == SampleTypeInt32

	// Set the extension fields (as appropriate)
	formatCode := effectiveFormatCode
	var validBitsPerSample *uint16
	var channelMask *uint32
	var subFormat *FormatCode
	if isExtensible {
		formatCode = FormatCodeExtensible
		validBitsPerSample = &bitsPerSample
		channelMask = new(uint32)
		subFormat = &effectiveFormatCode
	}

	return FormatChunkData{
		FormatCode:         formatCode,
		ChannelCount:       channelCount,
		FrameRate:          frameRate,
		ByteRate:           bytesPerSecond,
		BlockAlign:         blockAlign,
		BitsPerSample:      bitsPerSample,
		ValidBitsPerSample: validBitsPerSample,
		ChannelMask:        channelMask,
		SubFormat:          subFormat,
	}
}

// EffectiveFormatCode will return either c.FormatCode or c.SubFormat,
// depending on whether Extensible mode is enabled or not. The result will
// always be either FormatCodePCM or FormatCodeIEEEFloat on success.
func (c *FormatChunkData) EffectiveFormatCode() (FormatCode, error) {
	if c.FormatCode == FormatCodeExtensible {
		if c.SubFormat == nil {
			return FormatCode(0xFF), ErrFmtChunkMissingSubFormat
		}

		return *c.SubFormat, nil
	}

	return c.FormatCode, nil
}

// ChunkSize returns the total size of this chunk in bytes. The chunk size does
// not include the 8 byte header associated with all chunks.
func (c *FormatChunkData) ChunkSize() uint32 {
	size := uint32(16)
	if c.FormatCode == FormatCodeIEEEFloat {
		size += 2
	} else if c.FormatCode == FormatCodeExtensible {
		size += 24
	}

	return size
}

// Serialize packs this chunk into a []byte according to the wave spec.
func (c *FormatChunkData) Serialize() ([]byte, error) {

	buffer := &bytes.Buffer{}
	buffer.Grow(int(c.ChunkSize()))

	writeUint16(buffer, uint16(c.FormatCode))
	writeUint16(buffer, c.ChannelCount)
	writeUint32(buffer, c.FrameRate)
	writeUint32(buffer, c.ByteRate)
	writeUint16(buffer, c.BlockAlign)
	writeUint16(buffer, c.BitsPerSample)

	if c.FormatCode == FormatCodeIEEEFloat {
		writeUint16(buffer, 0) // Size of the extension. Must be given for non-PCM data
	} else if c.FormatCode == FormatCodeExtensible {

		// Verify that the required fields have been provided
		if c.ValidBitsPerSample == nil || c.ChannelMask == nil || c.SubFormat == nil {
			return nil, ErrFmtChunkInvalidExtensible
		}

		writeUint16(buffer, 22)
		writeUint16(buffer, *c.ValidBitsPerSample)
		writeUint32(buffer, *c.ChannelMask)

		guid := [16]byte{
			0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x10, 0x00,
			0x80, 0x00, 0x00, 0xAA,
			0x00, 0x38, 0x9B, 0x71,
		}
		binary.LittleEndian.PutUint16(guid[:2], uint16(*c.SubFormat))
		buffer.Write(guid[:])
	}

	return buffer.Bytes(), nil
}

// DeserializeFormatChunk reads a FormatChunkData structure from the provided
// []byte input. Errors will be thrown if the data is obviously structurally
// corrupted, but no checking is performed on the validity of the fields
// themselves.
func DeserializeFormatChunk(data []byte) (*FormatChunkData, error) {

	const (
		minFactPayloadSize              = 16
		minFactPayloadWithExtensionSize = 18
		minExtensibleSize               = 22
	)
	if len(data) < minFactPayloadSize {
		return nil, ErrFmtChunkCorruptedPayload
	}

	formatCode := FormatCode(readUint16(data[:2]))
	channelCount := readUint16(data[2:4])
	frameRate := readUint32(data[4:8])
	byteRate := readUint32(data[8:12])
	blockAlign := readUint16(data[12:14])
	bitsPerSample := readUint16(data[14:16])

	var validBitsPerSample *uint16
	var channelMask *uint32
	var subFormat *FormatCode

	// Process any extensions (if present)
	if len(data) >= minFactPayloadWithExtensionSize {
		extensionSize := readUint16(data[16:18])
		if extensionSize >= minExtensibleSize {

			bps := readUint16(data[18:20])
			validBitsPerSample = &bps

			cm := readUint32(data[20:24])
			channelMask = &cm

			sf := FormatCode(readUint16(data[24:26]))
			subFormat = &sf

			// We'll ignore the rest of the GUID (data[26:40]) and any other
			// extension fields.
			// ...
		}
	}

	return &FormatChunkData{
		FormatCode:         formatCode,
		ChannelCount:       channelCount,
		FrameRate:          frameRate,
		ByteRate:           byteRate,
		BlockAlign:         blockAlign,
		BitsPerSample:      bitsPerSample,
		ValidBitsPerSample: validBitsPerSample,
		ChannelMask:        channelMask,
		SubFormat:          subFormat,
	}, nil
}

// ------------------------------------------------------------------------- //
// Fact chunk
// ------------------------------------------------------------------------- //

var (
	FactChunkID = [4]byte{'f', 'a', 'c', 't'}

	ErrFactChunkCorruptedPayload = errors.New("detected corrupted 'fact' payload")
)

// NewFactChunk returns a 'fact' Chunk containing the given FactChunkData.
// The 'fact' chunk contains some additional metadata about the audio file
// not present in the format chunk and is typically required when non-PCM data
// is used.
func NewFactChunk(data *FactChunkData) Chunk {
	factData := data.Serialize()
	return Chunk{
		ID:   FactChunkID,
		Size: uint32(len(factData)),
		Body: factData,
	}
}

type FactChunkData struct {

	// SampleFrames is the number of audio frames present in the file. This
	// is usually used to calculate playing time, especially for compressed
	// data. Play time == Sample Frames * Sample Rate (Samples per second)
	SampleFrames uint32
}

// ChunkSize returns the total size of this chunk in bytes. The chunk size does
// not include the 8 byte header associated with all chunks.
func (c FactChunkData) ChunkSize() uint32 {
	return 4
}

// Serialize packs this data into a []byte according to the wave spec.
func (c FactChunkData) Serialize() []byte {
	return uint32ToBytes(c.SampleFrames)
}

// DeserializeFactChunk reads a FactChunkData structure from the provided
// []byte input.
func DeserializeFactChunk(data []byte) (*FactChunkData, error) {

	if len(data) < 4 {
		return nil, ErrFactChunkCorruptedPayload
	}

	return &FactChunkData{
		SampleFrames: readUint32(data),
	}, nil
}

// ------------------------------------------------------------------------- //
// Data chunk
// ------------------------------------------------------------------------- //

var (
	DataChunkID = [4]byte{'d', 'a', 't', 'a'}
)

// NewDataChunkHeader returns the header for a 'data' chunk that contains the
// given number of bytes. The 'data' chunk normally contains all the actual
// audio data, but this function returns only the header.
//
// NOTE: 'dataSize' should NOT include any padding added to ensure that the
// body will have an even number of bytes.
func NewDataChunkHeader(dataSize uint32) Chunk {
	return Chunk{
		ID:   DataChunkID,
		Size: dataSize,
		Body: nil,
	}
}

// ------------------------------------------------------------------------- //
// Cue chunk
// ------------------------------------------------------------------------- //

var (
	CueChunkID = [4]byte{'c', 'u', 'e', ' '}

	ErrCueChunkCorruptedPayload = errors.New("detected corrupted 'cue ' payload")
)

type CuePoint struct {
	ID           uint32
	Position     uint32
	FCCChunk     [4]byte
	ChunkStart   uint32
	BlockStart   uint32
	SampleOffset uint32
}

type CueChunkData struct {
	CuePoints []CuePoint
}

// DeserializeCueChunkData reads a CueChunkData structure from the provided
// []byte input. Errors will be thrown if the data is obviously structurally
// corrupted, but no checking is performed on the validity of the fields
// themselves.
func DeserializeCueChunkData(data []byte) (*CueChunkData, error) {

	const (
		minCuePayloadSize = 4
		cuePointSize      = 24
	)

	if len(data) < minCuePayloadSize {
		return nil, ErrCueChunkCorruptedPayload
	}

	numCuePoints := int(readUint32(data[:4]))
	if len(data) < numCuePoints*cuePointSize {
		return nil, ErrCueChunkCorruptedPayload
	}

	cuePoints := make([]CuePoint, numCuePoints)
	for i := 0; i < numCuePoints; i++ {
		cuePoints[i].ID = readUint32(data[4:8])
		cuePoints[i].Position = readUint32(data[8:12])
		copy(cuePoints[i].FCCChunk[:], data[12:16])
		cuePoints[i].ChunkStart = readUint32(data[16:20])
		cuePoints[i].BlockStart = readUint32(data[20:24])
		cuePoints[i].SampleOffset = readUint32(data[24:28])
	}
	return &CueChunkData{
		CuePoints: cuePoints,
	}, nil
}

// ------------------------------------------------------------------------- //
// Helpers
// ------------------------------------------------------------------------- //

func uint32ToBytes(val uint32) []byte {
	scratch := make([]byte, 4)
	binary.LittleEndian.PutUint32(scratch, val)
	return scratch
}

func uint16ToBytes(val uint16) []byte {
	scratch := make([]byte, 2)
	binary.LittleEndian.PutUint16(scratch, val)
	return scratch
}

func writeUint16(buffer *bytes.Buffer, val uint16) {
	scratch := make([]byte, 2)
	binary.LittleEndian.PutUint16(scratch, val)
	_, _ = buffer.Write(scratch)
}

func writeUint32(buffer *bytes.Buffer, val uint32) {
	scratch := make([]byte, 4)
	binary.LittleEndian.PutUint32(scratch, val)
	_, _ = buffer.Write(scratch)
}

func readUint16(buffer []byte) uint16 {
	return binary.LittleEndian.Uint16(buffer)
}

func readUint32(buffer []byte) uint32 {
	return binary.LittleEndian.Uint32(buffer)
}
