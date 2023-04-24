package wave

import (
	"fmt"
	"time"
)

// A Header is a preprocessed view of the beginning of a wave file, typically
// used when reading wave files (as opposed to writing them).
type Header struct {

	// Data read from the 'fmt' chunk in the wave file
	FormatData FormatChunkData

	// Data read from the 'fact' chunk in the wave file (if present). Not all
	// wave files will have 'fact' chunks.
	FactData *FactChunkData

	// Represents the total number of bytes of audio data that can be read from
	// this wave file.
	DataBytes uint32

	// Contains any Chunks that were not explicitly handled by this library.
	AdditionalChunks []Chunk
}

// Validate performs a series of cross-calculations on this Header to ensure
// that it is internally consistent. If Validate returns nil, this Header has
// passed all checks. If Validate returns an error, that error will describe
// what integrity check failed.
func (h *Header) Validate() error {

	m := uint32(h.FormatData.BitsPerSample / 8)

	// Format code
	if !h.FormatData.FormatCode.IsValid() {
		return fmt.Errorf("format code: '%d' was not recognized", h.FormatData.FormatCode)
	}

	// Bytes per second
	expectedByteRate := h.FormatData.FrameRate * m * uint32(h.FormatData.ChannelCount)
	if h.FormatData.ByteRate != expectedByteRate {
		return fmt.Errorf(
			"byte rate: '%d' did not match expected result: '%d'",
			h.FormatData.ByteRate,
			expectedByteRate,
		)
	}

	// Block align
	expectedBlockAlign := uint16(m) * h.FormatData.ChannelCount
	if h.FormatData.BlockAlign != expectedBlockAlign {
		return fmt.Errorf(
			"block align: '%d' did not match expected result: '%d'",
			h.FormatData.BlockAlign,
			expectedBlockAlign,
		)
	}

	// Sample frames
	if h.FactData != nil {
		expectedSampleFrames := h.FrameCount()
		if h.FactData.SampleFrames != expectedSampleFrames {
			return fmt.Errorf(
				"sample frames: '%d' did not match expected result: '%d'",
				h.FactData.SampleFrames,
				expectedSampleFrames,
			)
		}
	}

	return nil
}

// SampleType returns the SampleType that should be used when reading data
// associated with this Header.
func (h *Header) SampleType() (SampleType, error) {
	fc, err := h.FormatData.EffectiveFormatCode()
	if err != nil {
		return SampleType(-1), err
	}

	if fc == FormatCodePCM {
		switch h.FormatData.BitsPerSample {
		case 8:
			return SampleTypeUint8, nil
		case 16:
			return SampleTypeInt16, nil
		case 24:
			return SampleTypeInt24, nil
		case 32:
			return SampleTypeInt32, nil
		default:
			return SampleType(-1), fmt.Errorf("unknown PCM type: '%d' bits per sample", h.FormatData.BitsPerSample)
		}
	}

	// IEEE float
	switch h.FormatData.BitsPerSample {
	case 32:
		return SampleTypeFloat32, nil
	case 64:
		return SampleTypeFloat64, nil
	default:
		return SampleType(-1), fmt.Errorf("unknown IEEE type: '%d' bits per sample", h.FormatData.BitsPerSample)
	}
}

// FrameRate returns frame rate for the wave file associated with this header,
// measured in frames/second.
func (h *Header) FrameRate() uint32 {
	return h.FormatData.FrameRate
}

// ChannelCount returns the number of channels of audio data present in the
// wave file associated with this header. The channel count also determines
// the number of blocks that must be read to return a single sample.
func (h *Header) ChannelCount() uint16 {
	return h.FormatData.ChannelCount
}

// FrameCount returns the total number of audio frames present in the wave file
// associated with this header.
func (h *Header) FrameCount() uint32 {
	return h.DataBytes / uint32(h.FormatData.BlockAlign)
}

// SampleCount returns the total number of samples present in the wave file
// associated with this header.
func (h *Header) SampleCount() uint32 {
	return h.DataBytes / uint32(h.FormatData.BitsPerSample/8)
}

// PlayTime estimates the length of the wave file associated with this header.
func (h *Header) PlayTime() time.Duration {

	// Calculate value in seconds, but convert to nanoseconds for time.Duration
	seconds := float64(h.FrameCount()) / float64(h.FormatData.FrameRate)
	return time.Duration(seconds * 1e9)
}
