package core

type AudioReader interface {
	ReadUint8([]uint8) (int, error)
	ReadInt16([]int16) (int, error)
	ReadInt24([]int32) (int, error)
	ReadInt32([]int32) (int, error)
	ReadFloat32([]float32) (int, error)
	ReadFloat64([]float64) (int, error)
}

// A Decoder is an io.Reader that serves to translate audio samples from one
// type to another at runtime. This allows for interoperability with external
// libraries/systems that don't support the native (source) type. For example,
// some audio playback libraries do not have native support for Int24 samples.
// A Decoder can be used to transform those samples into Int16 in real time.
type Decoder struct {
	baseReader AudioReader
	fn         decoderFunc
}

func NewDecoder(
	r AudioReader,
	sourceType SampleType,
	destType SampleType,
) *Decoder {
	return &Decoder{
		baseReader: r,
		fn:         decoderFuncs[key{sourceType, destType}],
	}
}

func (d *Decoder) Read(p []byte) (int, error) {
	return d.fn(d.baseReader, p)
}
