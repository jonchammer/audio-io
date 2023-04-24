package wave

import (
	"encoding/binary"
	"errors"
	"io"
)

type Reader struct {
	baseReader io.Reader
	header     *Header
}

func NewReader(
	baseReader io.Reader,
) (*Reader, error) {
	return &Reader{
		baseReader: baseReader,
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
		return 0, errors.New("wave format is not uint8")
	}

	// For uint8 data, we can read directly into the buffer given by the caller.
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
		return 0, errors.New("wave format is not int16")
	}

	n := sampleType.Size()
	maxBytes := len(data) * n
	buffer := make([]byte, maxBytes)
	bytesRead, err := io.ReadFull(r.baseReader, buffer)
	samplesRead := bytesRead / n
	for i := 0; i < samplesRead; i++ {
		data[i] = int16(binary.LittleEndian.Uint16(buffer[n*i:]))
	}
	if err == io.ErrUnexpectedEOF {
		err = nil
	}
	return samplesRead, err
}

func (r *Reader) ReadInt24(data []int32) (int, error) {
	return 0, nil
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
		return 0, errors.New("wave format is not int32")
	}

	n := sampleType.Size()
	maxBytes := len(data) * n
	buffer := make([]byte, maxBytes)
	bytesRead, err := io.ReadFull(r.baseReader, buffer)
	samplesRead := bytesRead / n
	for i := 0; i < samplesRead; i++ {
		data[i] = int32(binary.LittleEndian.Uint32(buffer[n*i:]))
	}
	if err == io.ErrUnexpectedEOF {
		err = nil
	}
	return samplesRead, err
}

func (r *Reader) ReadFloat32(data []float32) (int, error) {
	return 0, nil
}

func (r *Reader) ReadFloat64(data []float64) (int, error) {
	return 0, nil
}
