package main

import (
	"fmt"
	"math"
	"os"

	"github.com/jonchammer/audio-io/core"
	"github.com/jonchammer/audio-io/wave"
)

func main() {

	// Open the provided file for reading
	file, err := os.Open("voice.wav")
	if err != nil {
		failF(err)
	}
	defer func() {
		_ = file.Close()
	}()

	// Create a reader using the file
	reader := wave.NewReader(file)

	// Print some useful metadata about the wave file to the console
	err = printWaveMetadata(reader)
	if err != nil {
		failF(err)
	}

	// Read the audio data as 'normalized' (or dequantized) float64 samples
	normalizedAudioData, err := readNormalizedAudioData(reader)
	if err != nil {
		failF(err)
	}
	fmt.Printf("Successfully read '%d' audio samples\n", len(normalizedAudioData))
	fmt.Println()

	// Compute some example statistics on the normalized audio data
	min, max := minMax(normalizedAudioData)
	factor := math.Max(math.Abs(min), math.Abs(max))
	gainDb := 20 * math.Log10(factor)
	fmt.Printf("Min value (normalized): %f\n", min)
	fmt.Printf("Max value (normalized): %+f\n", max)
	fmt.Printf("Estimated gain (dB):    %f\n", gainDb)
}

func printWaveMetadata(r *wave.Reader) error {

	// The header contains most relevant metadata, including details on how the
	// samples should be interpreted.
	header, err := r.Header()
	if err != nil {
		return err
	}

	// Well-formed wave files should have a known sample type, corresponding
	// to one of the wave.SampleTypeXXX constants.
	sampleType, _ := header.SampleType()

	// Print some relevant metadata for the file.
	fmt.Printf("Successfully read wave header\n")
	fmt.Printf("File Size (bytes):      %d\n", header.ReportedFileSizeBytes)
	fmt.Printf("Sample Type:            %s\n", sampleType)
	fmt.Printf("Bytes per Sample:       %d\n", sampleType.Size())
	fmt.Printf("Frame Rate:             %d\n", header.FrameRate())
	fmt.Printf("Bit Rate (bits/second): %d\n", header.BitRate())
	fmt.Printf("Channel Count:          %d\n", header.ChannelCount())
	fmt.Printf("Frame Count:            %d\n", header.FrameCount())
	fmt.Printf("Sample Count:           %d\n", header.SampleCount())
	fmt.Printf("Play Time:              %v\n", header.PlayTime())

	// We can optionally perform some basic integrity checks on the header to
	// validate that it is logically consistent. A failed validation check
	// doesn't necessary mean that the audio data can't be recovered, but it
	// does mean that the data presented above may not be trustworthy.
	msg := "PASSED"
	if err = header.Validate(); err != nil {
		msg = "FAILED - " + err.Error()
	}
	fmt.Printf("Validation:             %s\n", msg)
	fmt.Println()

	// This library knows how to interpret several common WAVE chunks, but if
	// a particular chunk isn't recognized, the user has an option to deal with
	// it themselves.
	if len(header.AdditionalChunks) != 0 {
		fmt.Println("Additional Chunks:")
		for i, c := range header.AdditionalChunks {
			fmt.Printf("%d - ID:   '%s'\n", i, string(c.ID[:]))
			fmt.Printf("    Size: %d\n", c.Size)
		}
		fmt.Println()
	}

	return nil
}

// readNormalizedAudioData demonstrates how to determine the correct data type
// for the audio samples based on the wave file metadata. In this example, we
// convert the samples to float64 and dequantize them to the range [-1, 1] so
// that we can process audio data in a fairly standard way.
func readNormalizedAudioData(r *wave.Reader) ([]float64, error) {

	// NOTE: If the header has already been read, a cached version will be
	// returned.
	header, err := r.Header()
	if err != nil {
		return nil, err
	}

	// Ensure the sample type is known
	sampleType, err := header.SampleType()
	if err != nil {
		return nil, err
	}

	var dequantizedAudioSamples []float64

	switch sampleType {
	case wave.SampleTypeUint8:
		{
			data := make([]uint8, header.SampleCount())
			_, err = r.ReadUint8(data)
			if err != nil {
				return nil, err
			}
			dequantizedAudioSamples = core.DequantizeUint8(data)
		}

	case wave.SampleTypeInt16:
		{
			data := make([]int16, header.SampleCount())
			_, err = r.ReadInt16(data)
			if err != nil {
				return nil, err
			}
			dequantizedAudioSamples = core.DequantizeInt16(data)
		}

	case wave.SampleTypeInt24:
		{
			data := make([]int32, header.SampleCount())
			_, err = r.ReadInt24(data)
			if err != nil {
				return nil, err
			}
			dequantizedAudioSamples = core.DequantizeInt24(data)
		}

	case wave.SampleTypeInt32:
		{
			data := make([]int32, header.SampleCount())
			_, err = r.ReadInt32(data)
			if err != nil {
				return nil, err
			}
			dequantizedAudioSamples = core.DequantizeInt32(data)
		}

	case wave.SampleTypeFloat32:
		{
			data := make([]float32, header.SampleCount())
			_, err = r.ReadFloat32(data)
			if err != nil {
				return nil, err
			}
			dequantizedAudioSamples = core.DequantizeFloat32(data)
		}

	case wave.SampleTypeFloat64:
		{
			dequantizedAudioSamples = make([]float64, header.SampleCount())
			_, err = r.ReadFloat64(dequantizedAudioSamples)
			if err != nil {
				return nil, err
			}
		}
	}

	return dequantizedAudioSamples, nil
}

// minMax computes the min and max values for the given slice.
func minMax(samples []float64) (float64, float64) {
	min := samples[0]
	max := samples[0]
	for _, s := range samples {
		if s < min {
			min = s
		}
		if s > max {
			max = s
		}
	}
	return min, max
}

func failF(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)
}
