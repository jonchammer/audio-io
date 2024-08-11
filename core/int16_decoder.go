package core

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

// An Int16Decoder is an io.Reader that serves to translate audio samples from
// any source SampleType to SampleTypeInt16 at runtime. For example, some audio
// playback libraries do not have native support for Int24 samples. In those
// situations, an Int16Decoder can be used to transform those samples into
// their int16 equivalents in real time.
//
// Note that when converting from a larger bit depth to int16, some information
// will necessarily be lost.
//
// The zero value for Int16Decoder is invalid. Instances should be created
// using NewInt16Decoder instead.
type Int16Decoder struct {
	baseReader io.Reader
	sourceType SampleType
	buffer     []byte
}

func NewInt16Decoder(
	r io.Reader,
	sourceType SampleType,
) *Int16Decoder {
	return &Int16Decoder{
		baseReader: r,
		sourceType: sourceType,
		buffer:     nil,
	}
}

func (d *Int16Decoder) Read(out []byte) (int, error) {

	// We can rely on ReadInt16 to do most of the heavy lifting. We'll give it
	// an []int16 view of 'out' to write into, since that will generally be
	// more ergonomic to work with.
	samples, err := d.ReadInt16(AliasAs[int16](out))

	// Read requires that we return the number of bytes read, not samples.
	// Note that this differs from the convention used by ReadInt16, which
	// returns the number of *elements* read.
	return samples * 2, err
}

func (d *Int16Decoder) ReadInt16(out []int16) (int, error) {

	if err := d.validateLittleEndian(); err != nil {
		return 0, err
	}

	// Read the input into a temporary buffer, perform the conversion to int16,
	// and write the results to 'out'.
	var samples int
	var err error

	switch d.sourceType {
	case SampleTypeUint8:
		samples, err = d.readUint8(out)
	case SampleTypeInt16:
		samples, err = d.readInt16(out)
	case SampleTypeInt24:
		samples, err = d.readInt24(out)
	case SampleTypeInt32:
		samples, err = d.readInt32(out)
	case SampleTypeFloat32:
		samples, err = d.readFloat32(out)
	default:
		samples, err = d.readFloat64(out)
	}
	return samples, err
}

func (d *Int16Decoder) validateLittleEndian() error {

	// This implementation currently relies on some 'unsafe' magic to cast
	// between the raw data bytes (little-endian) and normal types. Therefore,
	// we need to guard against big-endian host machines.
	//
	// A slower path might eventually be introduced for big-endian host
	// machines, but they are fairly uncommon in 20XX, so we'll ignore them for
	// now.
	isLittleEndian := binary.NativeEndian.Uint16([]byte{0x12, 0x34}) == uint16(0x3412)
	if !isLittleEndian {
		return errors.New("big-endian architectures currently not supported")
	}
	return nil
}

// readChunk pulls up to 'maxBytes' from the base reader into this decoder's
// internal buffer, returning the number of bytes actually read and an error.
//
// readChunk has the same semantics as io.ReadFull:
//   - If 'maxBytes' are read, 'maxBytes' is returned with no error
//   - If fewer than 'maxBytes' are read (but more than 0), the number of bytes
//     read will be returned with an io.ErrUnexpectedEOF error.
//   - If 0 bytes are read, 0 bytes will be returned with an io.EOF error.
func (d *Int16Decoder) readChunk(maxBytes int) (int, error) {

	// Buffer management. If we haven't yet allocated a buffer, we'll do so
	// now. If the user is now asking for more bytes than they have in the
	// past, we'll increase the size of the buffer.
	if d.buffer == nil {
		d.buffer = make([]byte, maxBytes)
	} else if len(d.buffer) < maxBytes {
		d.buffer = append(make([]byte, 0, maxBytes), d.buffer...)
	}

	// Read as many as 'maxBytes' elements into the internal buffer. Note that
	// the buffer may actually have space for more elements if 'readChunk'
	// was called earlier with a larger value for 'maxBytes'.
	return io.ReadFull(d.baseReader, d.buffer[:maxBytes])
}

func (d *Int16Decoder) readUint8(out []int16) (int, error) {

	// Read raw 'uint8' bytes into buffer
	n, err := d.readChunk(len(out))
	samplesRead := n

	// Map uint8 vals to int16 vals and write to 'out'
	for i := 0; i < samplesRead; i++ {
		out[i] = uint8ToInt16(d.buffer[i])
	}

	// Return the number of elements (not bytes) and the error
	return samplesRead, err
}

func (d *Int16Decoder) readInt16(out []int16) (int, error) {
	// Read raw 'int16' bytes into buffer (to maintain symmetry with other
	// source types)
	n, err := d.readChunk(len(out) * 2)
	samplesRead := n / 2

	// Copy the int16 vals directly into 'out'
	int16Vals := AliasAs[int16](d.buffer[:samplesRead])
	for i, int16Val := range int16Vals {
		out[i] = int16Val
	}

	// Return the number of elements (not bytes) and the error
	return samplesRead, err
}

func (d *Int16Decoder) readInt24(out []int16) (int, error) {

	// Read raw 'int24' bytes into buffer
	n, err := d.readChunk(len(out) * 3)
	samplesRead := n / 3

	// Map int24 vals to int16 vals and write to 'out'
	const mask = 0x01 << (24 - 1)
	for i := 0; i < samplesRead; i++ {

		// Reconstruct the int24 val from the raw bytes
		x := (int32(d.buffer[3*i+2]) << 16) | (int32(d.buffer[3*i+1]) << 8) | int32(d.buffer[3*i])
		int24Val := (x ^ mask) - mask

		out[i] = int24ToInt16(int24Val)
	}

	// Return the number of elements (not bytes) and the error
	return samplesRead, err
}

func (d *Int16Decoder) readInt32(out []int16) (int, error) {

	// Read raw 'int32' bytes into buffer
	n, err := d.readChunk(len(out) * 4)
	samplesRead := n / 4

	// Map int32 vals to int16 vals and write to 'out'
	int32Vals := AliasAs[int32](d.buffer[:samplesRead])
	for i, int32Val := range int32Vals {
		out[i] = int32ToInt16(int32Val)
	}

	// Return the number of elements (not bytes) and the error
	return samplesRead, err
}

func (d *Int16Decoder) readFloat32(out []int16) (int, error) {

	// Read raw 'float32' bytes into buffer
	n, err := d.readChunk(len(out) * 4)
	samplesRead := n / 4

	// Map float32 vals to int16 vals and write to 'out'
	float32Vals := AliasAs[float32](d.buffer[:samplesRead])
	for i, float32Val := range float32Vals {
		out[i] = float32ToInt16(float32Val)
	}

	// Return the number of elements (not bytes) and the error
	return samplesRead, err
}

func (d *Int16Decoder) readFloat64(out []int16) (int, error) {

	// Read raw 'float64' bytes into buffer
	n, err := d.readChunk(len(out) * 8)
	samplesRead := n / 8

	// Map float64 vals to int16 vals and write to 'out'
	float64Vals := AliasAs[float64](d.buffer[:samplesRead])
	for i, float64Val := range float64Vals {
		out[i] = float64ToInt16(float64Val)
	}

	// Return the number of elements (not bytes) and the error
	return samplesRead, err
}

// ------------------------------------------------------------------------- //
// Helpers
// ------------------------------------------------------------------------- //

// TODO: These implementations are very unoptimized. They are converting the
//   input to a float64 using the dequantization process and then quantizing
//   the results as int16. In many cases, it should be possible to get the
//   same result without using float64 as an intermediate representation.

func uint8ToInt16(x uint8) int16 {
	m := [2]float64{255.0 / 32512.0, 1.0 / 127.0}
	b := [2]float64{-1.0, -128.0 / 127}
	idx := (x & 0x80) >> 7
	dequantized := m[idx]*float64(x) + b[idx]
	return int16((dequantized * 32767.5) - 0.5)
}

func int24ToInt16(x int32) int16 {
	const (
		minInt24 = -1 << 23
		maxInt24 = 1<<23 - 1
	)

	sign := (x & minInt24) >> 23
	divisor := float64(maxInt24) - float64(sign)
	dequantized := float64(x) / divisor
	return int16((dequantized * 32767.5) - 0.5)
}

func int32ToInt16(x int32) int16 {
	sign := (x & math.MinInt32) >> 31
	divisor := float64(math.MaxInt32) - float64(sign)
	dequantized := float64(x) / divisor
	return int16((dequantized * 32767.5) - 0.5)
}

func float32ToInt16(x float32) int16 {
	return int16((x * 32767.5) - 0.5)
}

func float64ToInt16(x float64) int16 {
	return int16((x * 32767.5) - 0.5)
}
