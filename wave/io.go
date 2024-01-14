package wave

import (
	"errors"
	"io"
)

var (
	ErrIOInvalid24BitInput      = errors.New("length of input byte slice is not divisible by 3")
	ErrIOInvalid24BitOutputSize = errors.New("length of output does not equal len(input) / 3")
)

// WriteMany is a general extension of the io.Writer interface that allows
// multiple buffers to be written a single writer in series.
func WriteMany(w io.Writer, buffers ...[]byte) (int, error) {
	totalBytesWritten := 0
	for _, buffer := range buffers {
		n, err := w.Write(buffer)
		totalBytesWritten += n
		if err != nil {
			return totalBytesWritten, err
		}
	}
	return totalBytesWritten, nil
}

// WritePackedInt24 is a replacement for binary.Write() designed for int24 data
// that explicitly packs data from int32 elements into 3-byte chunks in
// little-endian byte order, without separations between elements. On success,
// 3 * len(data) bytes will be written to 'w'. On failure, WritePackedInt24
// will return the total number of bytes written so far and an error.
func WritePackedInt24(w io.Writer, data []int32) (int, error) {

	// Input (int32 - visualized in big endian order):
	// 00000000 xxxxxxxx yyyyyyyy zzzzzzzz
	//
	// Output (int24 - little endian order):
	// zzzzzzzz yyyyyyyy xxxxxxxx
	totalBytesWritten := 0
	for _, x := range data {
		tmp := [3]byte{
			byte(x & 0x000000FF),
			byte((x & 0x0000FF00) >> 8),
			byte((x & 0x00FF0000) >> 16),
		}

		n, err := w.Write(tmp[:])
		totalBytesWritten += n
		if err != nil {
			return totalBytesWritten, err
		}
	}
	return totalBytesWritten, nil
}

// ReadPackedInt24 is a replacement for binary.Read() designed for int24 data
// (in little-endian byte order) that explicitly unpacks the source values into
// int32 elements as they are read.
//
// If successful, this function will return a []int32 of size len(input) / 3.
//
// This function will return ErrIOInvalid24BitInput if the size of the input is
// not evenly divisible by 3.
func ReadPackedInt24(input []byte) ([]int32, error) {
	output := make([]int32, len(input)/3)
	err := ReadPackedInt24Into(input, output)
	return output, err
}

// ReadPackedInt24Into is a replacement for binary.Read() designed for int24
// data (in little-endian byte order) that explicitly unpacks the source values
// into int32 elements as they are read.
//
// As opposed to ReadPackedInt24, this function assumes that the caller has
// already allocated space for the output.
//
// Errors:
//   - ErrIOInvalid24BitInput - The size of the input is not evenly divisible
//     by 3.
//   - ErrIOInvalid24BitOutputSize - The size of the output is not equal to
//     the size of the input / 3.
func ReadPackedInt24Into(input []byte, output []int32) error {
	const mask = 0x01 << (24 - 1)

	if len(input)%3 != 0 {
		return ErrIOInvalid24BitInput
	}
	if len(output) != len(input)/3 {
		return ErrIOInvalid24BitOutputSize
	}

	// Input (int24 - little endian order):
	// zzzzzzzz yyyyyyyy xxxxxxxx
	//
	// Output (int32 - visualized in big endian order). Note that we have to
	// handle sign extension ourselves when converting back to int32.
	// SSSSSSSS xxxxxxxx yyyyyyyy zzzzzzzz
	//
	// NOTE: This is a very old trick for implementing sign extension in a high
	// level language.
	// References:
	//   - http://graphics.stanford.edu/~seander/bithacks.html#FixedSignExtend
	j := 0
	for i := 0; i < len(input); i += 3 {

		// The lower 24 bits of 'x' will be correct, but 'x' is not sign-extended
		x := (int32(input[i+2]) << 16) | (int32(input[i+1]) << 8) | int32(input[i])

		// Use the mask to handle sign extension
		output[j] = (x ^ mask) - mask
		j += 1
	}

	return nil
}
