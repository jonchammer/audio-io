package core

import (
	"unsafe"
)

// AliasAs interprets the given []byte as a slice of a different primitive type
// (e.g. []int32) without copying the data. As such, any changes made to the
// aliased slice will be reflected in the original source.
//
// NOTE: Users of AliasAs should be aware of both the byte ordering of 'src'
// and the native byte ordering of their host machine. This implementation
// makes no attempt to correct for discrepancies between the two. Incorrect
// assumptions about byte ordering may result in functional discrepancies.
func AliasAs[T uint8 | int16 | Int24 | int32 | int64 | float32 | float64](
	src []byte,
) []T {

	// Determine how many elements should be in the resulting slice. E.g. if
	// 'src' is 24 bytes, the length depends on T:
	//     uint8 (1 byte)  == 24 elements
	//     int16 (2 bytes) == 12 elements
	//     int32 (4 bytes) ==  6 elements
	//     int64 (8 bytes) ==  3 elements
	//   float32 (4 bytes) ==  6 elements
	//   float64 (8 bytes) ==  3 elements
	var t T
	aliasedLength := len(src) / int(unsafe.Sizeof(t))

	// Construct a new slice that uses the same data but has the correct type
	return unsafe.Slice(
		(*T)(unsafe.Pointer(unsafe.SliceData(src))),
		aliasedLength,
	)
}
