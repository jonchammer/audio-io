package wave

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

var (
	ErrReaderUnexpectedUint8   = errors.New("wave header indicates that this file does not use uint8 samples")
	ErrReaderUnexpectedInt16   = errors.New("wave header indicates that this file does not use int16 samples")
	ErrReaderUnexpectedInt24   = errors.New("wave header indicates that this file does not use int24 samples")
	ErrReaderUnexpectedInt32   = errors.New("wave header indicates that this file does not use int32 samples")
	ErrReaderUnexpectedFloat32 = errors.New("wave header indicates that this file does not use float32 samples")
	ErrReaderUnexpectedFloat64 = errors.New("wave header indicates that this file does not use float64 samples")
)

// A Reader is used to extract raw audio samples from its .wav representation.
// A Reader is created using NewReader, and data can be extracted using one of
// the ReadXXX methods. The caller can choose to read the entire file into a
// single buffer (useful for small files), or to read blocks of samples (useful
// for streaming).
//
// The Reader type generally enforces type safety when working with audio
// samples. If a .wav file was originally created using 16-bit integer samples,
// that audio data can only be safely read using the ReadInt16 method.
// Similarly, if the file was created using 32-bit IEEE float samples, that
// audio data can only be safely read using the ReadFloat32 method.
//
// Reader.Header returns a Header struct that contains useful metadata about
// the file, including what type should be used when reading samples, the total
// number of samples, the number of channels, etc. The caller will probably
// want to use this information to determine exactly how samples should be
// read, but calling Reader.Header is entirely optional if you don't need this
// information.
//
// Example usage (error handling omitted):
//
//	// Prepare data source
//	file, _ := os.Open("example.wav")
//	defer func() {
//	 	_ = file.Close()
//	}()
//
//	// Create a reader and get the header
//	r := NewReader(file)
//	header, _ := r.Header()
//
//	// In this example, we'll assume that we know ahead of time that
//	// 'example.wav' uses 16-bit integer samples. r.ReadInt16() will return an
//	// error if that assumption is incorrect.
//	data := make([]int16, header.SampleCount())
//	_, _ = r.ReadInt16(data)
type Reader struct {
	baseReader io.ReadSeeker
	dataReader io.Reader
	header     *Header
	buffer     []byte
}

// NewReader is a constructor function, used to create Reader instances.
// 'baseReader' is an io.ReadSeeker that represents the raw .wav data. This
// will commonly be an os.File or a bytes.Reader.
func NewReader(
	baseReader io.ReadSeeker,
) *Reader {
	return &Reader{
		baseReader: baseReader,
		dataReader: nil,
		header:     nil,
		buffer:     nil,
	}
}

// Header returns a Header object containing the metadata for the file (e.g.
// sample type, sample count, channel count, etc.)
func (r *Reader) Header() (*Header, error) {

	// If we haven't yet read the header, do that first. Results will be cached
	// after the first invocation.
	if r.header == nil {

		// Read the raw RIFF chunk data from the base reader.
		fileSize, riffData, err := ReadRIFFChunk(r.baseReader)
		if err != nil {
			return nil, err
		}

		// Parse the RIFF chunk as a Header.
		header, err := parseHeaderFromRIFFChunk(fileSize, riffData)
		if err != nil {
			return nil, err
		}

		r.header = header

		// We'll set up a LimitedReader to ensure the user doesn't
		// inadvertently try to read more bytes from the 'data' chunk than are
		// actually present.
		r.dataReader = io.LimitReader(r.baseReader, int64(header.DataBytes))
	}

	return r.header, nil
}

// ReadUint8 reads a chunk of uint8 samples from the data source and places
// them into the provided buffer. As many as len(data) samples could be read
// in a single call. The actual number of samples read will be returned, along
// with an error if data could not be read or the EOF has been reached.
//
// ReadUint8 will return an ErrReaderUnexpectedUint8 error if the underlying
// audio data is not representable as a []uint8 (e.g. float32 samples). If
// the caller is not sure of the data representation, they should call
// Header.SampleType to determine which ReadXXX function to call.
//
// NOTE: Audio samples will be **interleaved** if the data source uses multiple
// channels. core.DeinterleaveSlices can be used to de-interleave (split into
// separate channels) if needed.
func (r *Reader) ReadUint8(data []uint8) (int, error) {

	// Make sure we've read the header already
	header, err := r.Header()
	if err != nil {
		return 0, err
	}

	// Verify that the sample type is correct
	sampleType, err := header.SampleType()
	if err != nil {
		return 0, err
	}
	if sampleType != SampleTypeUint8 {
		return 0, ErrReaderUnexpectedUint8
	}

	// For uint8 data, we can read directly into the slice provided by the
	// caller. No buffering is required.
	return r.dataReader.Read(data)
}

// ReadInt16 reads a chunk of int16 samples from the data source and places
// them into the provided buffer. As many as len(data) samples could be read
// in a single call. The actual number of samples read will be returned, along
// with an error if data could not be read or the EOF has been reached.
//
// ReadInt16 will return an ErrReaderUnexpectedInt16 error if the underlying
// audio data is not representable as a []int16 (e.g. float32 samples). If
// the caller is not sure of the data representation, they should call
// Header.SampleType to determine which ReadXXX function to call.
//
// NOTE: Audio samples will be **interleaved** if the data source uses multiple
// channels. core.DeinterleaveSlices can be used to de-interleave (split into
// separate channels) if needed.
func (r *Reader) ReadInt16(data []int16) (int, error) {

	// Make sure we've read the header already
	header, err := r.Header()
	if err != nil {
		return 0, err
	}

	// Verify that the sample type is correct
	sampleType, err := header.SampleType()
	if err != nil {
		return 0, err
	}
	if sampleType != SampleTypeInt16 {
		return 0, ErrReaderUnexpectedInt16
	}

	n := sampleType.Size()
	bytesRead, err := r.readChunk(len(data) * n)
	samplesRead := bytesRead / n
	for i := 0; i < samplesRead; i++ {
		data[i] = int16(binary.LittleEndian.Uint16(r.buffer[n*i:]))
	}
	return samplesRead, err
}

// ReadInt24 reads a chunk of 24-bit samples from the data source (where each
// individual sample is represented as an int32 in the range
// [-8388608, 8388607]) and places those samples into the provided buffer. As
// many as len(data) samples could be read in a single call. The actual number
// of samples read will be returned, along with an error if data could not be
// read or the EOF has been reached.
//
// ReadInt24 will return an ErrReaderUnexpectedInt24 error if the underlying
// audio data is not representable as a []int32 (e.g. float32 samples). If
// the caller is not sure of the data representation, they should call
// Header.SampleType to determine which ReadXXX function to call.
//
// NOTE: Audio samples will be **interleaved** if the data source uses multiple
// channels. core.DeinterleaveSlices can be used to de-interleave (split into
// separate channels) if needed.
func (r *Reader) ReadInt24(data []int32) (int, error) {

	// Make sure we've read the header already
	header, err := r.Header()
	if err != nil {
		return 0, err
	}

	// Verify that the sample type is correct
	sampleType, err := header.SampleType()
	if err != nil {
		return 0, err
	}
	if sampleType != SampleTypeInt24 {
		return 0, ErrReaderUnexpectedInt24
	}

	n := sampleType.Size()
	bytesRead, err := r.readChunk(len(data) * n)
	samplesRead := bytesRead / n
	extraBytes := bytesRead % n
	_ = ReadPackedInt24Into(r.buffer[:(bytesRead-extraBytes)], data[:samplesRead])
	return samplesRead, err
}

// ReadInt32 reads a chunk of int32 samples from the data source and places
// them into the provided buffer. As many as len(data) samples could be read
// in a single call. The actual number of samples read will be returned, along
// with an error if data could not be read or the EOF has been reached.
//
// ReadInt32 will return an ErrReaderUnexpectedInt32 error if the underlying
// audio data is not representable as a []int32 (e.g. float32 samples). If
// the caller is not sure of the data representation, they should call
// Header.SampleType to determine which ReadXXX function to call.
//
// NOTE: Audio samples will be **interleaved** if the data source uses multiple
// channels. core.DeinterleaveSlices can be used to de-interleave (split into
// separate channels) if needed.
func (r *Reader) ReadInt32(data []int32) (int, error) {

	// Make sure we've read the header already
	header, err := r.Header()
	if err != nil {
		return 0, err
	}

	// Verify that the sample type is correct
	sampleType, err := header.SampleType()
	if err != nil {
		return 0, err
	}
	if sampleType != SampleTypeInt32 {
		return 0, ErrReaderUnexpectedInt32
	}

	n := sampleType.Size()
	bytesRead, err := r.readChunk(len(data) * n)
	samplesRead := bytesRead / n
	for i := 0; i < samplesRead; i++ {
		data[i] = int32(binary.LittleEndian.Uint32(r.buffer[n*i:]))
	}
	return samplesRead, err
}

// ReadFloat32 reads a chunk of float32 samples from the data source and places
// them into the provided buffer. As many as len(data) samples could be read
// in a single call. The actual number of samples read will be returned, along
// with an error if data could not be read or the EOF has been reached.
//
// ReadFloat32 will return an ErrReaderUnexpectedFloat32 error if the
// underlying audio data is not representable as a []float32 (e.g. int16
// samples). If the caller is not sure of the data representation, they should
// call Header.SampleType to determine which ReadXXX function to call.
//
// NOTE: Audio samples will be **interleaved** if the data source uses multiple
// channels. core.DeinterleaveSlices can be used to de-interleave (split into
// separate channels) if needed.
func (r *Reader) ReadFloat32(data []float32) (int, error) {

	// Make sure we've read the header already
	header, err := r.Header()
	if err != nil {
		return 0, err
	}

	// Verify that the sample type is correct
	sampleType, err := header.SampleType()
	if err != nil {
		return 0, err
	}
	if sampleType != SampleTypeFloat32 {
		return 0, ErrReaderUnexpectedFloat32
	}

	n := sampleType.Size()
	bytesRead, err := r.readChunk(len(data) * n)
	samplesRead := bytesRead / n
	for i := 0; i < samplesRead; i++ {
		data[i] = math.Float32frombits(binary.LittleEndian.Uint32(r.buffer[n*i:]))
	}
	return samplesRead, err
}

// ReadFloat64 reads a chunk of float64 samples from the data source and places
// them into the provided buffer. As many as len(data) samples could be read
// in a single call. The actual number of samples read will be returned, along
// with an error if data could not be read or the EOF has been reached.
//
// ReadFloat64 will return an ErrReaderUnexpectedFloat64 error if the
// underlying audio data is not representable as a []float64 (e.g. int16
// samples). If the caller is not sure of the data representation, they should
// call Header.SampleType to determine which ReadXXX function to call.
//
// NOTE: Audio samples will be **interleaved** if the data source uses multiple
// channels. core.DeinterleaveSlices can be used to de-interleave (split into
// separate channels) if needed.
func (r *Reader) ReadFloat64(data []float64) (int, error) {

	// Make sure we've read the header already
	header, err := r.Header()
	if err != nil {
		return 0, err
	}

	// Verify that the sample type is correct
	sampleType, err := header.SampleType()
	if err != nil {
		return 0, err
	}
	if sampleType != SampleTypeFloat64 {
		return 0, ErrReaderUnexpectedFloat64
	}

	n := sampleType.Size()
	bytesRead, err := r.readChunk(len(data) * n)
	samplesRead := bytesRead / n
	for i := 0; i < samplesRead; i++ {
		data[i] = math.Float64frombits(binary.LittleEndian.Uint64(r.buffer[n*i:]))
	}
	return samplesRead, err
}

// readChunk pulls up to 'maxBytes' from the data reader into this reader's
// internal buffer, returning the number of bytes actually read and an error.
//
// readChunk has the same semantics as io.ReadFull:
//   - If 'maxBytes' are read, 'maxBytes' is returned with no error
//   - If fewer than 'maxBytes' are read (but more than 0), the number of bytes
//     read will be returned with an io.ErrUnexpectedEOF error.
//   - If 0 bytes are read, 0 bytes will be returned with an io.EOF error.
func (r *Reader) readChunk(
	maxBytes int,
) (int, error) {

	// Buffer management. If we haven't yet allocated a buffer, we'll do so
	// now. If the user is now asking for more bytes than they have in the
	// past, we'll increase the size of the buffer.
	if r.buffer == nil {
		r.buffer = make([]byte, maxBytes)
	} else if len(r.buffer) < maxBytes {
		r.buffer = append(make([]byte, 0, maxBytes), r.buffer...)
	}

	// Read as many as 'maxBytes' elements into the internal buffer. Note that
	// the buffer may actually have space for more elements if 'readChunk'
	// was called earlier with a larger value for 'maxBytes'.
	return io.ReadFull(r.dataReader, r.buffer[:maxBytes])
}
