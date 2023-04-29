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

type Reader struct {
	baseReader io.Reader
	header     *Header
	buffer     []byte
}

func NewReader(
	baseReader io.Reader,
) (*Reader, error) {
	return &Reader{
		baseReader: baseReader,
		header:     nil,
		buffer:     nil,
	}, nil
}

func (r *Reader) Header() (*Header, error) {

	// When the header has already been read, we can return it directly
	if r.header != nil {
		return r.header, nil
	}

	// Otherwise, we'll have to read it from the base reader
	// TODO: Fill in

	return r.header, nil
}

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
	return r.baseReader.Read(data)
}

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

// readChunk pulls up to 'maxBytes' from the base reader into this reader's
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
		r.buffer = append(make([]byte, maxBytes), r.buffer...)
	}

	// Read as many as 'maxBytes' elements into the internal buffer. Note that
	// the buffer may actually have space for more elements if 'readChunk'
	// was called earlier with a larger value for 'maxBytes'.
	return io.ReadFull(r.baseReader, r.buffer[:maxBytes])
}
