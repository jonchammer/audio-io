package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/jonchammer/audio-io/core"
	"github.com/jonchammer/audio-io/wave"
)

func main() {

	const (
		duration      = 5 * time.Second
		sampleRate    = 44100
		sineFrequency = 440.0 // The note 'A'
		sineGain      = 0.25  // Amplitude multiplier
	)

	// 1. Open a file to store the resulting .wav file.
	f, err := os.Create("example.wav")
	if err != nil {
		failF(err)
	}
	defer func() {
		_ = f.Close()
	}()

	// 2. Create a wave writer and set the properties we care about
	w, err := wave.NewWriter(
		f, wave.SampleTypeInt24, uint32(sampleRate),
	)
	if err != nil {
		failF(err)
	}
	defer func() {
		err = w.Flush()
		if err != nil {
			failF(err)
		}
	}()

	// 3. We'll use a channel to facilitate communication between a background
	// goroutine (responsible for generating blocks of audio data) and the
	// main goroutine (responsible for transforming those blocks and writing
	// them to disk).
	//
	// Note that individual audio blocks (not the entire audio file) are kept
	// in memory. As a result, memory usage is not expected to appreciably grow
	// while we're streaming.
	c := make(chan []float64)
	go generateSineWave(c, duration, sampleRate, sineFrequency, sineGain)
	for block := range c {
		err = w.WriteInterleavedInt24(core.QuantizeToInt24(block))
		if err != nil {
			failF(err)
		}
	}
}

func generateSineWave(
	c chan []float64,
	duration time.Duration,
	sampleRate int,
	frequency float64,
	gain float64,
) {

	const (
		Tau       = 2 * math.Pi
		BlockSize = 512
	)

	// Work out how many total samples we need based on the duration and the
	// sample rate.
	totalSamples := int(duration.Seconds() * float64(sampleRate))

	// We'll divide the samples into blocks to simulate a more real-world
	// scenario. All blocks but the last will be full sized. The last block
	// will contain any remaining samples.
	blockCount := totalSamples / BlockSize
	lastBlockSamples := totalSamples % BlockSize

	var x float64
	samplingInterval := 1.0 / float64(sampleRate)

	// Publish each full sized block to the channel
	for b := 0; b < blockCount; b++ {
		buffer := make([]float64, BlockSize)
		for i := 0; i < BlockSize; i++ {
			buffer[i] = gain * math.Sin(Tau*frequency*x)
			x += samplingInterval
		}
		c <- buffer
	}

	// Publish the last block
	buffer := make([]float64, lastBlockSamples)
	for i := 0; i < lastBlockSamples; i++ {
		buffer[i] = gain * math.Sin(Tau*frequency*x)
		x += samplingInterval
	}
	c <- buffer

	// Close the channel to signify that we're done producing audio data
	close(c)
}

func failF(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)
}
