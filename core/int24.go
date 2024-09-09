package core

import (
	"strconv"
)

const (
	MinInt24 = -1 << 23
	MaxInt24 = 1<<23 - 1
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
