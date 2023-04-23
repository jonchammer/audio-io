// Package wave contains types and functions that facilitate working with
// .wave files.
package wave

// References
//   - https://wavefilegem.com/how_wave_files_work.html
//   - http://www-mmsp.ece.mcgill.ca/Documents/AudioFormats/WAVE/WAVE.html

// ------------------------------------------------------------------------- //
// FormatCode
// ------------------------------------------------------------------------- //

// FormatCode is an enum defined by the Wave specification that dictates how
// audio samples are to be interpreted.
type FormatCode uint16

const (
	FormatCodePCM        FormatCode = 0x0001
	FormatCodeIEEEFloat  FormatCode = 0x0003
	FormatCodeExtensible FormatCode = 0xFFFE
)

// ------------------------------------------------------------------------- //
// SampleType
// ------------------------------------------------------------------------- //

// SampleType represents the type of audio data that can be accepted by a
// particular Writer.
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
