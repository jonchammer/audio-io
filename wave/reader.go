package wave

import (
	"io"
)

type Reader struct {
	baseReader io.Reader
}

func NewReader(
	baseReader io.Reader,
) (*Reader, error) {
	return &Reader{
		baseReader: baseReader,
	}, nil
}

func (r *Reader) Header() (*Header, error) {
	return nil, nil
}

func (r *Reader) ReadUint8(data []uint8) (int, error) {
	return 0, nil
}

func (r *Reader) ReadInt16(data []int16) (int, error) {
	return 0, nil
}

func (r *Reader) ReadInt24(data []int32) (int, error) {
	return 0, nil
}

func (r *Reader) ReadInt32(data []int32) (int, error) {
	return 0, nil
}

func (r *Reader) ReadFloat32(data []float32) (int, error) {
	return 0, nil
}

func (r *Reader) ReadFloat64(data []float64) (int, error) {
	return 0, nil
}
