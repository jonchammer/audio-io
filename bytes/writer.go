package bytes

import (
	"errors"
	"io"
)

var (
	// ErrTooLarge is passed to panic if memory cannot be allocated to store
	// data in a buffer.
	ErrTooLarge = errors.New("bytes.Writer: too large")
)

// A Writer is a variable-sized buffer of bytes that supports both the
// io.Writer and io.Seeker interfaces. In other words, it supports
// random-access, in-memory writes. The zero value for Writer is an empty
// buffer ready to use.
type Writer struct {
	buf []byte
	idx int
}

// Cap returns the capacity of the buffer's underlying byte slice, that is, the
// total space currently allocated for the buffer's data.
func (w *Writer) Cap() int {
	return cap(w.buf)
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to the
// buffer without another allocation.
//
// If n is negative, Grow will panic.
// If the buffer can't be grown properly, Grow will panic with ErrTooLarge.
func (w *Writer) Grow(n int) {
	if n < 0 {
		panic("bytes.Writer.Grow: negative count")
	}
	w.reserve(len(w.buf) + n)
}

// Len returns the current length of the buffer's underlying byte slice.
// Because Writer supports random-access, this is the only correct way to
// determine how many bytes are currently present in the buffer.
func (w *Writer) Len() int {
	return len(w.buf)
}

// Bytes returns a reference to the buffer's underlying byte slice. The
// returned slice is only valid for use until the next buffer modification.
func (w *Writer) Bytes() []byte {
	return w.buf
}

// Write appends the contents of p to the buffer, growing the buffer as
// needed. The return value n is the length of p; err is always nil. If the
// buffer becomes too large, Write will panic with ErrTooLarge.
//
// Note that the total number of bytes written will not necessarily align with
// the buffer's length. If Seek is called, some elements of the buffer may be
// overwritten over time.
func (w *Writer) Write(p []byte) (n int, err error) {

	m := len(p)
	minLength := w.idx + m

	// Increase slice capacity if needed
	w.reserve(minLength)

	// Increase slice length if necessary
	if minLength > len(w.buf) {
		w.buf = w.buf[:minLength]
	}

	// Copy elements from p into w.buf
	copy(w.buf[w.idx:], p)

	w.idx += m
	return m, nil
}

// Seek adjusts the write head for the buffer according to the io.Seeker
// contract. The resulting write head, relative to the start of the buffer,
// will be returned upon success. Seek can fail if the resolved index for the
// write head is out of bounds for the current buffer.
func (w *Writer) Seek(offset int64, whence int) (int64, error) {

	var proposedIdx int
	switch whence {
	case io.SeekStart:
		proposedIdx = int(offset)
	case io.SeekCurrent:
		proposedIdx = w.idx + int(offset)
	case io.SeekEnd:
		proposedIdx = len(w.buf) + int(offset)
	}

	if proposedIdx < 0 || proposedIdx > len(w.buf) {
		return 0, errors.New("resolved index for Seek() is out of bounds")
	}
	w.idx = proposedIdx

	return int64(w.idx), nil
}

func (w *Writer) reserve(n int) {

	// Deal with allocations that are too large
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()

	// If we have sufficient capacity already, there's nothing more to do.
	if n <= cap(w.buf) {
		return
	}

	// Otherwise, we'll need to reallocate and copy the elements over
	w.buf = append(make([]byte, 0, n), w.buf...)
}
