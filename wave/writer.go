package wave

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrWriterInvalidSampleType = errors.New("provided sample type is invalid")
	ErrWriterInvalidByteCount  = errors.New("an invalid number of bytes were written before the writer was closed")

	ErrWriterExpectedUint8   = errors.New("sample type was not set to uint8 when the writer was constructed")
	ErrWriterExpectedInt16   = errors.New("sample type was not set to int16 when the writer was constructed")
	ErrWriterExpectedInt24   = errors.New("sample type was not set to int24 when the writer was constructed")
	ErrWriterExpectedInt32   = errors.New("sample type was not set to int32 when the writer was constructed")
	ErrWriterExpectedFloat32 = errors.New("sample type was not set to float32 when the writer was constructed")
	ErrWriterExpectedFloat64 = errors.New("sample type was not set to float64 when the writer was constructed")
)

// A Writer is used to generate .wav files from raw audio samples. A Writer
// is created using NewWriter, and samples are written using one of the
// WriteXXX methods. Samples can be written over the span of multiple calls
// (e.g. to allow the caller to generate audio samples on the fly). After all
// audio samples are written, the caller is expected to call Flush(). Flush()
// ensures that all wave metadata is set properly.
//
// Example usage (error handling omitted):
//
//	w, _ := NewWriter(
//	    output, SampleTypeInt16, 44100, WithChannelCount(2),
//	)
//	defer func() {
//	    _ = w.Flush()
//	}
//	var audioData []int16 = ...
//	_ = w.WriteInt16(audioData)
type Writer struct {

	// Handles writes to the final .wav file (or buffer)
	baseWriter io.WriteSeeker

	// Determines what types of audio data this writer should accept at runtime
	sampleType SampleType

	// Metadata chunks. Most of this information be calculated when the writer
	// is created, but some fields cannot be determined until runtime. These
	// chunks may be written multiple times as new information is made
	// available. We may or may not have a 'fact' chunk at all.
	formatChunkData FormatChunkData
	factChunkData   *FactChunkData

	// The number of bytes of audio data written to 'baseWriter' so far
	dataBytes uint32
}

// NewWriter is a constructor function, used to create Writer instances.
//   - baseWriter - The base writer can be an os.File or any other type that
//     implements the io.WriteSeeker interface in the Go standard library.
//   - sampleType - The sample type determines which of the WriteXXX APIs can
//     be used. Note that it is the caller's responsibility to ensure data is
//     in the correct format, though the quantizers/dequantizers in the core
//     package can help with that process.
//   - frameRate - The frame rate is measured in frames per second. Common
//     values are 44100 Hz (normal for CD audio) and 48000 Hz (common for
//     cinema). Note that the term "sample rate" is more common, but isn't as
//     well-defined in terms of multi-channel data. We use "frame rate" to be
//     more precise about what is expected of the caller.
//
// WriterOptions can be used to provide additional optional inputs (e.g.
// setting the number of channels).
func NewWriter(
	baseWriter io.WriteSeeker,
	sampleType SampleType,
	frameRate uint32,
	opts ...WriterOption,
) (*Writer, error) {

	// Validate the required inputs
	if !sampleType.IsValid() {
		return nil, ErrWriterInvalidSampleType
	}

	// Process any optional inputs
	options := &writerOptions{
		channelCount: 1,
	}
	for _, opt := range opts {
		err := opt(options)
		if err != nil {
			return nil, err
		}
	}

	// The user-provided fields will determine the format chunk fields
	formatChunkData := NewFormatChunkData(
		options.channelCount, frameRate, sampleType,
	)

	// It's generally agreed that regular PCM data doesn't require a 'fact'
	// chunk. We'll add one in all other cases.
	var factChunkData *FactChunkData
	if formatChunkData.FormatCode != FormatCodePCM {
		factChunkData = &FactChunkData{
			SampleFrames: 0,
		}
	}

	return &Writer{
		baseWriter:      baseWriter,
		sampleType:      sampleType,
		formatChunkData: formatChunkData,
		factChunkData:   factChunkData,
		dataBytes:       0,
	}, nil
}

// WriteUint8 is used to add uint8 audio samples. Audio data is assumed to be
// organized into frames consisting of multiple samples, one sample per channel.
// WriteUint8 will fail if the SampleType of the Writer is not set to
// SampleTypeUint8.
func (w *Writer) WriteUint8(data []uint8) error {
	if w.sampleType != SampleTypeUint8 {
		return ErrWriterExpectedUint8
	}

	err := w.write(data)
	if err != nil {
		return err
	}

	w.dataBytes += uint32(w.sampleType.Size() * len(data))
	return nil
}

// WriteInt16 is used to add uint8 audio samples. Audio data is assumed to be
// organized into frames consisting of multiple samples, one sample per channel.
// WriteInt16 will fail if the SampleType of the Writer is not set to
// SampleTypeInt16.
func (w *Writer) WriteInt16(data []int16) error {

	if w.sampleType != SampleTypeInt16 {
		return ErrWriterExpectedInt16
	}

	err := w.write(data)
	if err != nil {
		return err
	}

	w.dataBytes += uint32(w.sampleType.Size() * len(data))
	return nil
}

// WriteInt24 is used to add uint8 audio samples. Audio data is assumed to be
// organized into frames consisting of multiple samples, one sample per channel.
// WriteInt24 will fail if the SampleType of the Writer is not set to
// SampleTypeInt24.
//
// NOTE: most programming languages (including Go) don't provide a native
// 24-bit integer type, so we usually use `int32` as a container type with the
// understanding that values are expected to fall in the range
// [-8388608, 8388607]. There is special logic where needed to pack and unpack
// values as 24-bit integers.
func (w *Writer) WriteInt24(data []int32) error {

	if w.sampleType != SampleTypeInt24 {
		return ErrWriterExpectedInt24
	}

	// NOTE: Specialized logic is needed for int24 compared to the other data
	// types.
	err := w.writeInt24(data)
	if err != nil {
		return err
	}

	w.dataBytes += uint32(w.sampleType.Size() * len(data))
	return nil
}

// WriteInt32 is used to add uint8 audio samples. Audio data is assumed to be
// organized into frames consisting of multiple samples, one sample per channel.
// WriteInt32 will fail if the SampleType of the Writer is not set to
// SampleTypeInt32.
func (w *Writer) WriteInt32(data []int32) error {

	if w.sampleType != SampleTypeInt32 {
		return ErrWriterExpectedInt32
	}

	err := w.write(data)
	if err != nil {
		return err
	}

	w.dataBytes += uint32(w.sampleType.Size() * len(data))
	return nil
}

// WriteFloat32 is used to add uint8 audio samples. Audio data is assumed to be
// organized into frames consisting of multiple samples, one sample per channel.
// WriteFloat32 will fail if the SampleType of the Writer is not set to
// SampleTypeFloat32.
func (w *Writer) WriteFloat32(data []float32) error {

	if w.sampleType != SampleTypeFloat32 {
		return ErrWriterExpectedFloat32
	}

	err := w.write(data)
	if err != nil {
		return err
	}

	w.dataBytes += uint32(w.sampleType.Size() * len(data))
	return nil
}

// WriteFloat64 is used to add uint8 audio samples. Audio data is assumed to be
// organized into frames consisting of multiple samples, one sample per channel.
// WriteFloat64 will fail if the SampleType of the Writer is not set to
// SampleTypeFloat64.
func (w *Writer) WriteFloat64(data []float64) error {

	if w.sampleType != SampleTypeFloat64 {
		return ErrWriterExpectedFloat64
	}

	err := w.write(data)
	if err != nil {
		return err
	}

	w.dataBytes += uint32(w.sampleType.Size() * len(data))
	return nil
}

// write is a common helper for most of the WriteXXX
// methods declared above.
func (w *Writer) write(data any) error {
	if w.dataBytes == 0 {
		err := w.writePreamble()
		if err != nil {
			return err
		}
	}

	// Seek to the end of the file and append the new block
	_, err := w.baseWriter.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	err = binary.Write(w.baseWriter, binary.LittleEndian, data)
	if err != nil {
		return err
	}

	return nil
}

// writeInt24 is a specialization of write to be used with int24 data.
func (w *Writer) writeInt24(data []int32) error {
	if w.dataBytes == 0 {
		err := w.writePreamble()
		if err != nil {
			return err
		}
	}

	// Seek to the end of the file and append the new block
	_, err := w.baseWriter.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	_, err = WritePackedInt24(w.baseWriter, data)
	if err != nil {
		return err
	}

	return nil
}

// Flush rewinds the underlying io.WriteSeeker back to the beginning of the
// file and overwrites the existing .wav file header. Flush must be called
// after all audio samples have been written to ensure that the file's metadata
// is up-to-date.
//
// Flush will commonly be placed in a defer statement, but this isn't a
// requirement of the API.
//
// Flush will fail if an invalid number of samples are written (e.g. an odd
// number of samples are written when the Writer is configured for two
// channels) with an ErrWriterInvalidByteCount.
func (w *Writer) Flush() error {

	// Validate that the total number of bytes written to the data chunk makes
	// sense in the context of this writer.
	remainder := w.dataBytes % uint32(w.formatChunkData.BlockAlign)
	if remainder != 0 {
		return ErrWriterInvalidByteCount
	}

	// Update the 'SampleFrames' field in the fact chunk (if necessary)
	if w.factChunkData != nil {
		w.factChunkData.SampleFrames = w.dataBytes / uint32(w.formatChunkData.BlockAlign)
	}

	// Add another byte to the data chunk for padding (if necessary)
	padding := w.dataBytes & 1
	if padding != 0 {
		_, err := w.baseWriter.Write(make([]byte, 1))
		if err != nil {
			return err
		}
	}

	// Rewind to the beginning of the file and rewrite the preamble with the
	// final (correct) values.
	err := w.writePreamble()
	if err != nil {
		return err
	}

	return nil
}

// writePreamble rewinds the base writer back to the beginning of the file and
// writes (or rewrites) the .wav preamble, leaving the write head at the first
// byte for audio data. It returns the total size of the preamble in bytes
// and an error (in case the write fails).
//
// Note that the size of the preamble is set when the Writer is created.
// It will not change at runtime (though individual values within the preamble
// might).
//
// The preamble will have this format:
//
//	Field    Length    Contents
//	ckID          4    "RIFF"
//	ckSize        4    Total number of remaining bytes in the file
//	  WAVEID      4    "WAVE"
//
//	  ckID        4    "fmt "
//	  ckSize      4    Size of format chunk (N). Usually 16, 18, or 40
//	    fmtData   N    Format chunk data
//
//	  ckID        4    "fact"                         <---+
//	  ckSize      4    Size of fact chunk (M). Usually 4. | Fact chunk optional
//	    factData  M    Fact chunk data                <---+
//
//	  ckID        4    "data"
//	  ckSize      4    Size of data (P)
func (w *Writer) writePreamble() error {

	// Seek to the beginning of the writer
	_, err := w.baseWriter.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// The preamble is represented by a hierarchy of chunks. The root chunk
	// describes (recursively) the entire file structure.
	preambleBytes := w.getRootChunk().Serialize()

	// Write the preamble to the writer
	_, err = w.baseWriter.Write(preambleBytes)
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) getRootChunk() Chunk {

	subChunks := make([]Chunk, 0, 4)

	// NOTE: It's safe to ignore the error here because we can be sure that the
	// format chunk is well-formed.
	formatChunk, _ := NewFormatChunk(&w.formatChunkData)
	subChunks = append(subChunks, formatChunk)

	// We may or may not have a fact chunk
	if w.factChunkData != nil {
		subChunks = append(subChunks, NewFactChunk(w.factChunkData))
	}

	// We'll only write the header for the data chunk. We won't touch any of
	// the audio data that's already been written.
	subChunks = append(subChunks, NewDataChunkHeader(w.dataBytes))

	return NewRIFFChunk(&RIFFChunkData{
		subChunks: subChunks,
	})
}

// ------------------------------------------------------------------------- //
// Writer Options
// ------------------------------------------------------------------------- //

type writerOptions struct {
	channelCount uint16
}

// WriterOption is a functional argument used as part of NewWriter.
type WriterOption func(*writerOptions) error

// WithChannelCount is used to set the number of audio channels as part of
// NewWriter. A channel count of 1 will be assumed as the default unless
// explicitly overwritten by the user.
//
// Note that all WriteXXX APIs assume that frames are contiguous, so all
// samples for a given frame should be placed next to one other in memory.
func WithChannelCount(channelCount uint16) WriterOption {
	return func(opts *writerOptions) error {
		opts.channelCount = channelCount
		return nil
	}
}
