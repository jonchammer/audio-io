package core

import (
	"fmt"
)

// ------------------------------------------------------------------------- //
// SampleType
// ------------------------------------------------------------------------- //

// SampleType represents the type of audio data that can be accepted by a
// particular Writer or the type of data that can be extracted from a Reader.
type SampleType int

const (
	SampleTypeUint8 SampleType = iota + 1
	SampleTypeInt16
	SampleTypeInt24
	SampleTypeInt32
	SampleTypeFloat32
	SampleTypeFloat64
)

// IsValid returns true if 's' represents a valid SampleType
func (s SampleType) IsValid() bool {
	return s >= SampleTypeUint8 && s <= SampleTypeFloat64
}

// Size returns the size of the sample, measured in bytes.
func (s SampleType) Size() int {
	switch s {
	case SampleTypeUint8:
		return 1
	case SampleTypeInt16:
		return 2
	case SampleTypeInt24:
		return 3
	case SampleTypeInt32:
		return 4
	case SampleTypeFloat32:
		return 4
	default:
		return 8
	}
}

func (s SampleType) String() string {
	switch s {
	case SampleTypeUint8:
		return "Uint8"
	case SampleTypeInt16:
		return "Int16"
	case SampleTypeInt24:
		return "Int24"
	case SampleTypeInt32:
		return "Int32"
	case SampleTypeFloat32:
		return "Float32"
	case SampleTypeFloat64:
		return "Float64"
	default:
		return fmt.Sprintf("SampleType(%d)", s)
	}
}
