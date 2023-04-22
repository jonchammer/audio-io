# audio-io

[![Release](https://img.shields.io/github/release/jonchammer/audio-io.svg?style=flat-square)](https://github.com/jonchammer/audio-io/releases)

`audio-io` is a collection of I/O utilities written in Go that simplify the
process of working with audio data. Key features include:
  * A `.wav` file writer that supports:
    - PCM `uint8`, `int16`, `int24`, and `int32` formats
    - IEEE float `float32` and `float64` formats
    - Arbitrary number of audio channels
    - Arbitrary sample rates
    - Memory-efficient streaming of audio data to disk (e.g. suitable for 
      real-time audio generation)
  * Quantizers/dequantizers
    - Suitable for conversions between the `uint8`, `int16`, `int24`, `int32`, 
      `float32`, and `float64` audio formats
  * Interleavers/deinterleavers
    - Used to simplify the process of working with multi-channel audio files

## Examples
Several examples that demonstrate how to use this library are included in the
`examples` folder. 
