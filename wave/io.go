package wave

import (
	"io"
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
