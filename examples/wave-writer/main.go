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
		sineGain      = 0.25  // Amplitude multiplier
	)

	// 1. Generate a single channel sine wave for testing. We'll generate the
	// data using high-precision float64 samples, but we'll quantize the data
	// before we write it to disk.
	samples := generateSineWave(duration, sampleRate, sineFrequency, sineGain)

	// 2. Open a file to store the resulting .wav file.
	f, err := os.Create("example.wav")
	if err != nil {
		failF(err)
	}
	defer func() {
		_ = f.Close()
	}()

	// 3. Save the audio data as a wave file
	err = saveAsWave(samples, f, wave.SampleTypeInt24, sampleRate)
	if err != nil {
		failF(err)
	}
}

func generateSineWave(
	duration time.Duration,
	sampleRate int,
	frequency float64,
	gain float64,
) []float64 {

	const Tau = 2 * math.Pi

	totalSamples := int(duration.Seconds() * float64(sampleRate))
	data := make([]float64, totalSamples)
	for i := 0; i < totalSamples; i++ {
		data[i] = gain * math.Sin(Tau*frequency*(float64(i)/float64(sampleRate)))
	}
	return data
}

func saveAsWave(
	data []float64,
	out io.WriteSeeker,
	format wave.SampleType,
	sampleRate int,
) error {

	// Create a writer and set the properties we care about
	w, err := wave.NewWriter(out, format, uint32(sampleRate))
	if err != nil {
		return err
	}

	// Quantize the audio data and write it to our wave writer
	switch format {
	case wave.SampleTypeUint8:
		err = w.WriteUint8(core.QuantizeToUint8(data))
	case wave.SampleTypeInt16:
		err = w.WriteInt16(core.QuantizeToInt16(data))
	case wave.SampleTypeInt24:
		err = w.WriteInt24(core.QuantizeToInt24(data))
	case wave.SampleTypeInt32:
		err = w.WriteInt32(core.QuantizeToInt32(data))
	case wave.SampleTypeFloat32:
		err = w.WriteFloat32(core.QuantizeToFloat32(data))
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
