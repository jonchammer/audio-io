package core

import (
	"encoding/binary"
	"errors"
	"io"
)

type AlignedBufferedReader struct {
	source         io.Reader
	alignAs        int
	buffer         []byte
	leftoverBuffer []byte
}

// NewAlignedBufferedReader returns a new AlignedBufferedReader instance that
// pulls data from 'source'.
//
// The 'alignAs' argument controls:
//  1. How many bytes will be returned. The buffer size will always be an
//     integer multiple of 'alignAs'.
//  2. How the bytes in the buffer will be ordered. On little-endian machines,
//     bytes will be returned in the same order that they were read, but on
//     big-endian machines, the bytes will be swizzled according to the value
//     of 'alignAs'. (E.g. If 'alignAs' is 4, every sequence of 4 bytes will
//     be reversed in the buffer compared to the original source).
func NewAlignedBufferedReader(
	source io.Reader,
	alignAs int,
) *AlignedBufferedReader {
	return &AlignedBufferedReader{
		source:         source,
		alignAs:        alignAs,
		buffer:         nil,
		leftoverBuffer: nil,
	}
}

// ReadBuffer reads a chunk of bytes from the source reader into an in-memory
// buffer and returns that buffer to the caller. The 'maxBytes' argument
// controls the maximum size of the buffer.
//
// Bytes that are read but not returned as part of the buffer will be cached
// until the next call to ReadBuffer and included with the next call, ensuring
// that data is never lost.
func (r *AlignedBufferedReader) ReadBuffer(maxBytes int) ([]byte, error) {

	// (Re)allocate the internal buffer if needed
	if r.buffer == nil || len(r.buffer) < maxBytes {
		r.buffer = make([]byte, maxBytes)
	}

	// Copy any leftover bytes into the beginning of the internal buffer and
	// clear the leftover buffer (retaining capacity for future use)
	leftoverSize := copy(r.buffer, r.leftoverBuffer)
	r.leftoverBuffer = r.leftoverBuffer[:0]

	// Read up to 'maxBytes' of additional data from the source into the
	// internal buffer
	n, err := r.source.Read(r.buffer[leftoverSize:maxBytes])
	bufferSize := n + leftoverSize

	// Decide how much of the internal buffer needs to be returned to the
	// caller (based on the 'alignAs' argument), and how much needs to be
	// copied into the leftover buffer.
	resultSize := (bufferSize / r.alignAs) * r.alignAs

	// Copy any leftover bytes into the leftover buffer
	newLeftoverSize := bufferSize - resultSize
	if cap(r.leftoverBuffer) < newLeftoverSize {
		r.leftoverBuffer = make([]byte, newLeftoverSize)
	} else {
		r.leftoverBuffer = r.leftoverBuffer[:newLeftoverSize]
	}
	copy(r.leftoverBuffer, r.buffer[resultSize:bufferSize])

	// TODO: Byte swizzling to account for big-endian machines

	return r.buffer[:resultSize], err
}

func (r *AlignedBufferedReader) LeftoverSize() int {
	return len(r.leftoverBuffer)
}

func validateLittleEndian() error {

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
