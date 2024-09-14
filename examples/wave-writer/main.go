// wave-writer demonstrates how to use the audio-io 'core' and 'wave' packages
// to generate a simple wave file containing a single channel audio signal.
//
// This example assumes that all audio data can be kept in memory at once. See
// the 'stream-writer' example for a slightly more sophisticated use case (where
// audio data is generated in blocks that can be written as they're generated).
package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"time"

	"github.com/jonchammer/audio-io/core"
	"github.com/jonchammer/audio-io/wave"
)

func main() {

	const (
		duration      = 2 * time.Second
		sampleRate    = 44100
		sineFrequency = 440.0 // The note 'A'
		sineGainDb    = -12.0 // Gain of the signal, measured in decibels
	)

	// 1. Generate a single channel sine wave for testing. We'll generate the
	// data using high-precision float64 samples, but we'll quantize the data
	// before we write it to disk.
	samples := generateSineWave(duration, sampleRate, sineFrequency, sineGainDb)

	// 2. Open a file to store the resulting .wav file.
	f, err := os.Create("example.wav")
	if err != nil {
		failF(err)
	}
	defer func() {
		_ = f.Close()
	}()

	// 3. Save the audio data as a wave file
	err = saveAsWave(samples, f, core.SampleTypeInt24, sampleRate)
	if err != nil {
		failF(err)
	}
}

func generateSineWave(
	duration time.Duration,
	sampleRate int,
	frequency float64,
	gainDb float64,
) []float64 {

	const Tau = 2 * math.Pi

	factor := math.Pow(10.0, gainDb/20.0)
	totalSamples := int(duration.Seconds() * float64(sampleRate))
	data := make([]float64, totalSamples)
	for i := 0; i < totalSamples; i++ {
		data[i] = factor * math.Sin(Tau*frequency*(float64(i)/float64(sampleRate)))
	}
	return data
}

func saveAsWave(
	data []float64,
	out io.WriteSeeker,
	format core.SampleType,
	sampleRate int,
) error {

	// Create a writer and set the properties we care about
	w, err := wave.NewWriter(out, format, uint32(sampleRate))
	if err != nil {
		return err
	}

	// Quantize the audio data and write it to our wave writer
	switch format {
	case core.SampleTypeUint8:
		samples := make([]uint8, len(data))
		core.ConvertFloat64ToUint8(samples, data)
		err = w.WriteUint8(samples)
	case core.SampleTypeInt16:
		samples := make([]int16, len(data))
		core.ConvertFloat64ToInt16(samples, data)
		err = w.WriteInt16(samples)
	case core.SampleTypeInt24:
		samples := make([]core.Int24, len(data))
		core.ConvertFloat64ToInt24(samples, data)
		err = w.WriteInt24(samples)
	case core.SampleTypeInt32:
		samples := make([]int32, len(data))
		core.ConvertFloat64ToInt32(samples, data)
		err = w.WriteInt32(samples)
	case core.SampleTypeFloat32:
		samples := make([]float32, len(data))
		core.ConvertFloat64ToFloat32(samples, data)
		err = w.WriteFloat32(samples)
	default:
		err = w.WriteFloat64(data)
	}
	if err != nil {
		return err
	}

	// The writer must be flushed to ensure that all metadata is up-to-date.
	err = w.Flush()
	if err != nil {
		return err
	}

	return err
}

func failF(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)
}
