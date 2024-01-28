// stream-reader demonstrates how to use the audio-io 'core' and 'wave'
// packages to stream (and process) individual blocks of audio data from a wave
// file.
//
// This example reads blocks of audio data, normalizes those samples
// (converting to float64), and then calculates some basic statistics for each
// block. The block-level statistics are used to generate a crude ASCII
// 'oscillogram' (a representation of the amplitude of the signal over time).
//
//	(+) 1.0 + ---------------------------------------- +
//	        | █                █                 █     |
//	        | █                █                 █     |
//	        | █   █   █    █   █    █   █   █    █     |
//	        | ██ ██   ███  █   ███  █   █   █    █     |
//	    0.0 + ████████████████████████████████████████ +
//	        | ██  █   █    █   █    █   █   █    █     |
//	        | █                                  █     |
//	        |                                          |
//	        |                                          |
//	(-) 1.0 + ---------------------------------------- +
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jonchammer/audio-io/core"
	"github.com/jonchammer/audio-io/wave"
)

func main() {

	// We'll divide the audio file into blocks of samples. This constant
	// determines how many blocks we'll end up with. It will also represent the
	// width of the graph (measured in characters)
	const maxBlocks = 80

	// Determines the height of the graph (measured in terminal rows)
	const graphHeight = 9

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

	// For the sake of this example, we'll determine the number of samples per
	// block dynamically. The first 'N - 1' blocks will have 'samplesBerBlock'
	// samples each, and the last block will consist of however many samples
	// are left over.
	header, err := reader.Header()
	if err != nil {
		failF(err)
	}
	samplesBerBlock := int(header.SampleCount() / (maxBlocks - 1))

	// We'll use a channel to facilitate communication between a background
	// goroutine (responsible for reading blocks of samples) and the main
	// goroutine (responsible for calculating some statistics for each block).
	//
	// Note that individual audio blocks (not the entire audio file) are kept
	// in memory. As a result, memory usage is not expected to appreciably grow
	// while we're streaming.
	c := make(chan []float64)
	go func() {
		err = readNormalizedSamples(c, reader, samplesBerBlock)
		if err != nil {
			failF(err)
		}
	}()

	minValues := make([]float64, 0)
	maxValues := make([]float64, 0)
	for block := range c {
		minValue, maxValue := minMax(block)
		minValues = append(minValues, minValue)
		maxValues = append(maxValues, maxValue)
	}

	// Just for fun, we'll render a simple time-based visualization of the
	// audio amplitude to the console.
	fmt.Printf("Processed '%d' audio blocks\n", len(minValues))
	drawGraph(minValues, maxValues, graphHeight)
}

// readNormalizedSamples reads blocks of audio samples from 'reader',
// normalizes them, and pushes the results into 'c'. Each block of samples will
// have at most 'maxSamples' elements in it. The last block will likely have
// fewer elements. 'c' will be closed when all samples have been read.
func readNormalizedSamples(
	c chan []float64,
	reader *wave.Reader,
	maxSamples int,
) error {

	// Ensure that the channel is closed, even if we fail to read data from
	// the file.
	defer close(c)

	// NOTE: If the header has already been read, a cached version will be
	// returned.
	header, err := reader.Header()
	if err != nil {
		return err
	}

	// Ensure the sample type is known
	sampleType, err := header.SampleType()
	if err != nil {
		return err
	}

	switch sampleType {
	case core.SampleTypeUint8:
		{
			for {
				data := make([]uint8, maxSamples)
				samplesRead, err := reader.ReadUint8(data)
				if err != nil && err == io.EOF {
					break
				}

				normalized := core.DequantizeUint8(data[:samplesRead])
				c <- normalized
			}
		}

	case core.SampleTypeInt16:
		{
			for {
				data := make([]int16, maxSamples)
				samplesRead, err := reader.ReadInt16(data)
				if err != nil && err == io.EOF {
					break
				}

				normalized := core.DequantizeInt16(data[:samplesRead])
				c <- normalized
			}
		}

	case core.SampleTypeInt24:
		{
			for {
				data := make([]int32, maxSamples)
				samplesRead, err := reader.ReadInt24(data)
				if err != nil && err == io.EOF {
					break
				}

				normalized := core.DequantizeInt24(data[:samplesRead])
				c <- normalized
			}
		}

	case core.SampleTypeInt32:
		{
			for {
				data := make([]int32, maxSamples)
				samplesRead, err := reader.ReadInt32(data)
				if err != nil && err == io.EOF {
					break
				}

				normalized := core.DequantizeInt32(data[:samplesRead])
				c <- normalized
			}
		}

	case core.SampleTypeFloat32:
		{
			for {
				data := make([]float32, maxSamples)
				samplesRead, err := reader.ReadFloat32(data)
				if err != nil && err == io.EOF {
					break
				}

				normalized := core.DequantizeFloat32(data[:samplesRead])
				c <- normalized
			}
		}

	case core.SampleTypeFloat64:
		{
			for {
				data := make([]float64, maxSamples)
				samplesRead, err := reader.ReadFloat64(data)
				if err != nil && err == io.EOF {
					break
				}

				c <- data[:samplesRead]
			}
		}
	}

	return nil
}

// minMax computes the min and max values for the given slice.
func minMax(samples []float64) (float64, float64) {
	minValue := samples[0]
	maxValue := samples[0]
	for _, s := range samples {
		if s < minValue {
			minValue = s
		}
		if s > maxValue {
			maxValue = s
		}
	}
	return minValue, maxValue
}

// drawGraph renders the provided minimum and maximum values as a simple ASCII
// chart, like the one given below. The x-axis represents individual blocks of
// audio (proportional to time), the y-axis represents the max/min amplitudes
// encountered in that block. The height of the graph (in rows) can be provided
// as input to increase the vertical resolution.
//
// NOTE: There are better ways of visualizing audio data, but this method works
// well enough for this demo.
//
// Example output (height = 9):
//
//	+1.0 + ---------------------------------------- +
//	     | █                █                 █     |
//	     | █                █                 █     |
//	     | █   █   █    █   █    █   █   █    █     |
//	     | ██ ██   ███  █   ███  █   █   █    █     |
//	 0.0 + ████████████████████████████████████████ +
//	     | ██  █   █    █   █    █   █   █    █     |
//	     | █                                  █     |
//	     |                                          |
//	     |                                          |
//	-1.0 + ---------------------------------------- +
func drawGraph(minValues []float64, maxValues []float64, height int) {

	width := len(minValues)

	// This 2D array will be our graph's "canvas". For simplicity when
	// rendering, we'll assume that (0, 0) is the top left corner of the graph,
	// meaning that the 'y' axis points downward and the 'x' axis points to
	// the right.
	cells := make([][]rune, height)
	for y := 0; y < height; y++ {
		cells[y] = make([]rune, width)
		for x := 0; x < width; x++ {
			cells[y][x] = ' '
		}
	}

	// Linearly map each element of the min/max arrays to a row number and then
	// fill in all cells between those rows to form a solid vertical bar. (This
	// is an extremely primitive way of rendering signal amplitudes, but it
	// will work well enough for this example).
	for x := 0; x < width; x++ {
		maxRow := int((float64(height) * (1 - maxValues[x])) / 2)
		minRow := int((float64(height) * (1 - minValues[x])) / 2)
		for y := maxRow; y <= minRow; y++ {
			cells[y][x] = '█'
		}
	}

	// Print the graph to the console
	fmt.Println("+1.0 +", strings.Repeat("-", width), "+")
	for y := 0; y < height; y++ {
		leftEdge := "     |"
		rightEdge := "|"
		if y == height/2 {
			leftEdge = " 0.0 +"
			rightEdge = "+"
		}

		fmt.Println(leftEdge, string(cells[y]), rightEdge)
	}
	fmt.Println("-1.0 +", strings.Repeat("-", width), "+")
}

func failF(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)
}
