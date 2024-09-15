package core

import (
	"errors"
	"io"
	"strconv"
)

const (
	MinInt24 = -1 << 23
	MaxInt24 = 1<<23 - 1
)

var (
	ErrIOInvalid24BitInput      = errors.New("length of input byte slice is not divisible by 3")
	ErrIOInvalid24BitOutputSize = errors.New("length of output does not equal len(input) / 3")
)

// Int24 is a packed representation of a 3-byte signed integer with values
// ranging from [-8388608, 8388607]. Int24 instances are guaranteed to take
// exactly 3 bytes of space. The bytes are of the Int24 instance are explicitly
// stored in little-endian order, regardless of the endian-ness of the host
// machine.
type Int24 [3]byte

func (i Int24) String() string {
	return strconv.Itoa(int(i.AsInt32()))
}

// AsInt32 returns an int32 representation of i.
func (i Int24) AsInt32() int32 {
	const mask = 0x01 << (24 - 1)
	x := (int32(i[2]) << 16) | (int32(i[1]) << 8) | int32(i[0])
	return (x ^ mask) - mask
}

// AsInt64 returns an int64 representation of i.
func (i Int24) AsInt64() int64 {
	const mask = 0x01 << (24 - 1)
	x := (int64(i[2]) << 16) | (int64(i[1]) << 8) | int64(i[0])
	return (x ^ mask) - mask
}

// Int24FromInt32 creates a new Int24 instance using the lowest 3 bytes of 'i'.
func Int24FromInt32(i int32) Int24 {
	return [3]byte{
		byte(i & 0x000000FF),
		byte((i & 0x0000FF00) >> 8),
		byte((i & 0x00FF0000) >> 16),
	}
}

// Int24FromInt64 creates a new Int24 instance using the lowest 3 bytes of 'i'.
func Int24FromInt64(i int64) Int24 {
	return [3]byte{
		byte(i & 0x000000FF),
		byte((i & 0x0000FF00) >> 8),
		byte((i & 0x00FF0000) >> 16),
	}
}

// WriteInt24Slice is a replacement for binary.Write() designed to work with
// slices of core.Int24. The data will be written in little-endian byte order
// with no additional padding between elements. On success, 3 * len(data)
// bytes will be written to 'w'. On failure, WriteInt24Slice will return the
// total number of bytes written so far and an error.
func WriteInt24Slice(w io.Writer, data []Int24) (int, error) {
	totalBytesWritten := 0
	for _, x := range data {
		n, err := w.Write(x[:])
		totalBytesWritten += n
		if err != nil {
			return totalBytesWritten, err
		}
	}
	return totalBytesWritten, nil
}

// ReadInt24 is a replacement for binary.Read() designed for core.Int24 data.
//
// If successful, this function will return a []int32 of size len(input) / 3.
//
// This function will return ErrIOInvalid24BitInput if the size of the input is
// not evenly divisible by 3.
func ReadInt24(input []byte) ([]Int24, error) {
	output := make([]Int24, len(input)/3)
	err := ReadInt24Into(input, output)
	return output, err
}

// ReadInt24Into is a replacement for binary.Read() designed for core.Int24
// data.
//
// As opposed to ReadInt24, this function assumes that the caller has already
// allocated space for the output.
//
// Errors:
//   - ErrIOInvalid24BitInput - The size of the input is not evenly divisible
//     by 3.
//   - ErrIOInvalid24BitOutputSize - The size of the output is not equal to
//     the size of the input / 3.
func ReadInt24Into(input []byte, output []Int24) error {
	if len(input)%3 != 0 {
		return ErrIOInvalid24BitInput
	}
	if len(output) != len(input)/3 {
		return ErrIOInvalid24BitOutputSize
	}

	j := 0
	for i := 0; i < len(input); i += 3 {
		output[j] = [3]byte{input[i], input[i+1], input[i+2]}
		j++
	}
	return nil
}
