package core

import (
	"fmt"
)

// ------------------------------------------------------------------------- //
// SampleType
// ------------------------------------------------------------------------- //

// SampleType is an enum that describes various ways of storing audio samples
// in memory. A SampleType determines bit depth (e.g. 24-bit) as well as
// encoding (e.g. PCM vs. IEEE float).
type SampleType int

const (

	// SampleTypeUint8 stores audio samples as unsigned 8-bit integers in the
	// range [-128 - 127]
	SampleTypeUint8 SampleType = iota + 1

	// SampleTypeInt16 stores audio samples as signed 16-bit integers in the
	// range [-65,536 - 65,535]
	SampleTypeInt16

	// SampleTypeInt24 stores audio samples as signed 24-bit integers in the
	// range [-8,388,608 - 8,388,607]
	SampleTypeInt24

	// SampleTypeInt32 stores audio samples as signed 32-bit integers in the
	// range [-2,147,483,648 - 2,147,483,647]
	SampleTypeInt32

	// SampleTypeFloat32 stores audio samples as 32-bit IEEE float values in
	// the range [0 - 1]
	SampleTypeFloat32

	// SampleTypeFloat64 stores audio samples as 64-bit IEEE float values in
	// the range [0 - 1]
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
