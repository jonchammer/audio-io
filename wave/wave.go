// Package wave contains types and functions that facilitate working with
// Wave (.wav) files.
package wave

import (
	"fmt"
)

// References
//   - https://wavefilegem.com/how_wave_files_work.html
//   - http://www-mmsp.ece.mcgill.ca/Documents/AudioFormats/WAVE/WAVE.html
//   - https://www.recordingblogs.com/wiki/wave-file-format

// ------------------------------------------------------------------------- //
// FormatCode
// ------------------------------------------------------------------------- //

// FormatCode is an enum defined by the Wave specification that dictates how
// the audio data in a wave file is to be interpreted.
type FormatCode uint16

const (
	FormatCodePCM        FormatCode = 0x0001
	FormatCodeIEEEFloat  FormatCode = 0x0003
	FormatCodeExtensible FormatCode = 0xFFFE
)

// IsValid returns true if 'f' represents a valid FormatCode
func (f FormatCode) IsValid() bool {
	return f == FormatCodePCM || f == FormatCodeIEEEFloat || f == FormatCodeExtensible
}

func (f FormatCode) String() string {
	switch f {
	case FormatCodePCM:
		return "PCM"
	case FormatCodeIEEEFloat:
		return "IEEE Float"
	case FormatCodeExtensible:
		return "Extensible"
	default:
		return fmt.Sprintf("FormatCode(%d)", f)
	}
}

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

// EffectiveFormatCode returns the FormatCode that matches 's'
func (s SampleType) EffectiveFormatCode() FormatCode {
	switch s {
	case SampleTypeFloat32, SampleTypeFloat64:
		return FormatCodeIEEEFloat
	default:
		return FormatCodePCM
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
