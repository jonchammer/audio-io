// Package wave contains types and functions that facilitate working with
// Wave (.wav) files.
package wave

import (
	"fmt"

	"github.com/jonchammer/audio-io/core"
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

// effectiveFormatCode returns the FormatCode that matches 's'
func effectiveFormatCode(s core.SampleType) FormatCode {
	switch s {
	case core.SampleTypeFloat32, core.SampleTypeFloat64:
		return FormatCodeIEEEFloat
	default:
		return FormatCodePCM
	}
}
