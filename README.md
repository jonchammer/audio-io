![GitHub Workflow Status (with branch)](https://img.shields.io/github/actions/workflow/status/jonchammer/audio-io/test.yml?branch=main&style=flat-square)
![License](https://img.shields.io/github/license/jonchammer/audio-io?style=flat-square)
[![Release](https://img.shields.io/github/release/jonchammer/audio-io.svg?style=flat-square)](https://github.com/jonchammer/audio-io/releases)
[![GoDoc](https://pkg.go.dev/badge/github.com/jonchammer/audio-io?status.svg)](https://pkg.go.dev/github.com/jonchammer/audio-io?tab=doc)

# Audio I/O

`audio-io` is a collection of I/O utilities written in Go that simplify the
process of working with audio data. Key features include:
  * A `.wav` file writer that supports:
    - PCM `uint8`, `int16`, `int24`, and `int32` formats
    - IEEE float `float32` and `float64` formats
    - Arbitrary number of audio channels
    - Arbitrary frame (or sample) rates
    - Memory-efficient streaming of audio data to disk (e.g. suitable for 
      real-time audio generation)
  * Quantizers/dequantizers
    - Suitable for conversions between the `uint8`, `int16`, `int24`, `int32`, 
      `float32`, and `float64` audio formats
  * Interleavers/deinterleavers
    - Used to simplify the process of working with multi-channel audio files

## Writing wave files
The `wave.Writer` type can be used to generate .wav files from a set of audio 
samples. It wraps an existing `io.WriteSeeker` such as an `io.File` or a 
`bytes.Writer` and handles all metadata required by the wave specification. 
(Note that `bytes.Writer` is an in-memory buffer provided in this repository).

```go
package main

import (
	"os"
	"github.com/jonchammer/audio-io/wave"
)

func main() {
	// NOTE: This example ignores error handling for the sake of brevity.
	
	// Create a file to store our formatted audio data
	f, _ := os.Create("example.wav")
	defer func() {
		_ = f.Close()
	}()
	
	// Create a wave.Writer that wraps 'f'. In this example, we're assuming
	// that samples will be in the 16-bit PCM format, that our sample rate is
	// 44100 samples/second, and that each sample consists of 2 blocks.
	w, _ := wave.NewWriter(
		f, wave.SampleTypeInt16, uint32(44100), wave.WithChannelCount(2),
	)
	defer func() {
		_ = w.Flush()
	}()
	
	// Write your audio samples to the writer. Audio data is assumed to be
	// organized into 'samples', where each sample consists of multiple blocks. 
	// The number of blocks in a sample corresponds to the number of channels.
	var samples []int16 = ...
	_ = w.WriteInt16(samples)
}
```

`wave.NewWriter` requires that the caller provide a base writer, the expected
data type for the audio samples, and the sample rate. `wave.Writer` supports 
the `uint8`, `int16`, `int24`, and `int32` PCM formats as well as the `float32` 
and `float64` IEEE formats. 

Other optional writer properties can be set via functional arguments. In the 
example above, the channel count has been set, overriding the default value 
of 1.

### Working with multiple channels
In this library, each audio **frame** consists of 1 or more **samples**, with 
one sample per audio channel. A sample is represented as a single number with a
Go type of `uint8`, `int16`, `int32`, `float32`, or `float64`. 

The API assumes that all *frames* are contiguous, meaning that audio data is 
organized by *sample*, rather than by *channel*.

The table below visually demonstrates how samples are laid out in a slice:

```text
+--------+--------+----------------------------+
| Frame |  Index |                       Value |
+--------+--------+----------------------------+
|        |      0 |       Channel 0 - Sample 0 | 
|        |      1 |       Channel 1 - Sample 0 | 
|   0    |      2 |       Channel 2 - Sample 0 | 
|                     ...                      |
|        |  N - 1 | Channel (N - 1) - Sample 0 |
+--------+--------+----------------------------+
|        |      N |       Channel 0 - Sample 1 | 
|        |  N + 1 |       Channel 1 - Sample 1 | 
|   1    |  N + 2 |       Channel 2 - Sample 1 | 
|                     ...                      |
|        | 2N - 1 | Channel (N - 1) - Sample 1 |
+--------+--------+----------------------------+
|                     ...                      |
```

Applications that generate multi-channel audio will sometimes choose to keep
other conventions (e.g. organizing data first by channel and then by sample, or 
using a 2-dimensional array to store samples). It is the responsibility of the 
caller to ensure that their audio data is in the format expected by this 
library.

**NOTE**: The `core` package includes some generic utilities for interleaving / 
deinterleaving arbitrary slices (`InterleaveSlices` and `DeinterleaveSlices`).
These may prove useful if your application uses the "2D array" approach to 
block organization.

### Streaming
The `wave.Writer` API was designed to easily support efficient streaming of 
data. Each call to `WriteXXX` **appends** data to the base `io.WriteSeeker` 
that was provided when the writer was created, meaning that those APIs can be 
called whenever a new buffer of samples is available. See the `stream-writer` 
example in the `examples` folder for more details.

### Quantization
The `core` package includes quantizers and dequantizers that allow for 
conversions between audio sample types (e.g. `int32` and `float32`). PCM types,
`uint8`, `int16`, `int24`, and `int32` are assumed to use the full range of 
those types, while IEEE types (`float32` and `float64`) are assumed to use the
[-1.0, 1.0] range.

A quick note on `int24`: most programming languages (including Go) don't 
provide a native 24-bit integer type, so we usually use `int32` as a container 
type with the understanding that values are expected to fall in the range 
[-8388608, 8388607]. There is special logic where needed in the library to pack 
and unpack 24-bit integers as appropriate.

## Examples
Several complete examples that demonstrate how to use this library are included 
in the `examples` folder. 

## Developer Information

Execute test suite manually:
```sh
go test ./...
```

Examine generated .wav files (MacOS): 
```sh
afinfo example.wav
```

Play generated .wav files (MacOS):
```sh
afplay example.wav
```
