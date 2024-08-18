package core

import (
	"errors"
	"fmt"
	"io"
)

// A SampleTypeConverter is an io.Reader that serves to translate audio samples
// from any source SampleType to any dest SampleType at runtime. For example, a
// particular library may only support int16 samples, but the raw data could
// use float64 samples. A SampleTypeConverter will translate each sample from
// float64 to int16 dynamically as data is read.
//
// Note that when Converting from a larger bit depth to a smaller one, some
// information will necessarily be lost.
//
// The zero value for SampleTypeConverter is invalid. Instances should be
// created using NewSampleTypeConverter instead.
type SampleTypeConverter struct {
	reader     *AlignedBufferedReader
	sourceType SampleType
	destType   SampleType
}

func NewSampleTypeConverter(
	r io.Reader,
	sourceType SampleType,
	destType SampleType,
) (*SampleTypeConverter, error) {

	if !sourceType.IsValid() {
		return nil, fmt.Errorf("source type '%s' is invalid", sourceType)
	}
	if !destType.IsValid() {
		return nil, fmt.Errorf("dest type '%s' is invalid", destType)
	}

	return &SampleTypeConverter{
		reader:     NewAlignedBufferedReader(r, sourceType.Size()),
		sourceType: sourceType,
		destType:   destType,
	}, nil
}

func (d *SampleTypeConverter) Read(out []byte) (int, error) {
	switch d.destType {
	case SampleTypeUint8:
		samples, err := d.ReadUint8(out)
		return samples * 1, err
	case SampleTypeInt16:
		samples, err := d.ReadInt16(AliasAs[int16](out))
		return samples * 2, err
	case SampleTypeInt24:
		samples, err := d.ReadInt24(AliasAs[Int24](out))
		return samples * 3, err
	case SampleTypeInt32:
		samples, err := d.ReadInt32(AliasAs[int32](out))
		return samples * 4, err
	case SampleTypeFloat32:
		samples, err := d.ReadFloat32(AliasAs[float32](out))
		return samples * 4, err
	default:
		samples, err := d.ReadFloat64(AliasAs[float64](out))
		return samples * 8, err
	}
}

func (d *SampleTypeConverter) ReadUint8(out []uint8) (int, error) {

	if d.destType != SampleTypeUint8 {
		return 0, errors.New("expected dest type to be uint8")
	}

	// Read as many as len(out) samples from the input. The length of 'inRaw'
	// will always be an integer multiple of the 'sampleSize'.
	sampleSize := d.sourceType.Size()
	inRaw, err := d.reader.ReadBuffer(sampleSize * len(out))
	samples := len(inRaw) / sampleSize

	// The data in 'inRaw' will be interpreted according to the source type and
	// then converted to the appropriate dest type. The converted data will be
	// written into 'out'.
	switch d.sourceType {
	case SampleTypeUint8:
		in := AliasAs[uint8](inRaw)
		ConvertUint8ToUint8(out, in)
	case SampleTypeInt16:
		in := AliasAs[int16](inRaw)
		ConvertInt16ToUint8(out, in)
	case SampleTypeInt24:
		in := AliasAs[Int24](inRaw)
		ConvertInt24ToUint8(out, in)
	case SampleTypeInt32:
		in := AliasAs[int32](inRaw)
		ConvertInt32ToUint8(out, in)
	case SampleTypeFloat32:
		in := AliasAs[float32](inRaw)
		ConvertFloat32ToUint8(out, in)
	default:
		in := AliasAs[float64](inRaw)
		ConvertFloat64ToUint8(out, in)
	}

	return samples, err
}

func (d *SampleTypeConverter) ReadInt16(out []int16) (int, error) {

	if d.destType != SampleTypeInt16 {
		return 0, errors.New("expected dest type to be int16")
	}

	// Read as many as len(out) samples from the input. The length of 'inRaw'
	// will always be an integer multiple of the 'sampleSize'.
	sampleSize := d.sourceType.Size()
	inRaw, err := d.reader.ReadBuffer(sampleSize * len(out))
	samples := len(inRaw) / sampleSize

	// The data in 'inRaw' will be interpreted according to the source type and
	// then converted to the appropriate dest type. The converted data will be
	// written into 'out'.
	switch d.sourceType {
	case SampleTypeUint8:
		in := AliasAs[uint8](inRaw)
		ConvertUint8ToInt16(out, in)
	case SampleTypeInt16:
		in := AliasAs[int16](inRaw)
		ConvertInt16ToInt16(out, in)
	case SampleTypeInt24:
		in := AliasAs[Int24](inRaw)
		ConvertInt24ToInt16(out, in)
	case SampleTypeInt32:
		in := AliasAs[int32](inRaw)
		ConvertInt32ToInt16(out, in)
	case SampleTypeFloat32:
		in := AliasAs[float32](inRaw)
		ConvertFloat32ToInt16(out, in)
	default:
		in := AliasAs[float64](inRaw)
		ConvertFloat64ToInt16(out, in)
	}

	return samples, err
}

func (d *SampleTypeConverter) ReadInt24(out []Int24) (int, error) {

	if d.destType != SampleTypeInt24 {
		return 0, errors.New("expected dest type to be int24")
	}

	// Read as many as len(out) samples from the input. The length of 'inRaw'
	// will always be an integer multiple of the 'sampleSize'.
	sampleSize := d.sourceType.Size()
	inRaw, err := d.reader.ReadBuffer(sampleSize * len(out))
	samples := len(inRaw) / sampleSize

	// The data in 'inRaw' will be interpreted according to the source type and
	// then converted to the appropriate dest type. The converted data will be
	// written into 'out'.
	switch d.sourceType {
	case SampleTypeUint8:
		in := AliasAs[uint8](inRaw)
		ConvertUint8ToInt24(out, in)
	case SampleTypeInt16:
		in := AliasAs[int16](inRaw)
		ConvertInt16ToInt24(out, in)
	case SampleTypeInt24:
		in := AliasAs[Int24](inRaw)
		ConvertInt24ToInt24(out, in)
	case SampleTypeInt32:
		in := AliasAs[int32](inRaw)
		ConvertInt32ToInt24(out, in)
	case SampleTypeFloat32:
		in := AliasAs[float32](inRaw)
		ConvertFloat32ToInt24(out, in)
	default:
		in := AliasAs[float64](inRaw)
		ConvertFloat64ToInt24(out, in)
	}

	return samples, err
}

func (d *SampleTypeConverter) ReadInt32(out []int32) (int, error) {

	if d.destType != SampleTypeInt32 {
		return 0, errors.New("expected dest type to be int32")
	}

	// Read as many as len(out) samples from the input. The length of 'inRaw'
	// will always be an integer multiple of the 'sampleSize'.
	sampleSize := d.sourceType.Size()
	inRaw, err := d.reader.ReadBuffer(sampleSize * len(out))
	samples := len(inRaw) / sampleSize

	// The data in 'inRaw' will be interpreted according to the source type and
	// then converted to the appropriate dest type. The converted data will be
	// written into 'out'.
	switch d.sourceType {
	case SampleTypeUint8:
		in := AliasAs[uint8](inRaw)
		ConvertUint8ToInt32(out, in)
	case SampleTypeInt16:
		in := AliasAs[int16](inRaw)
		ConvertInt16ToInt32(out, in)
	case SampleTypeInt24:
		in := AliasAs[Int24](inRaw)
		ConvertInt24ToInt32(out, in)
	case SampleTypeInt32:
		in := AliasAs[int32](inRaw)
		ConvertInt32ToInt32(out, in)
	case SampleTypeFloat32:
		in := AliasAs[float32](inRaw)
		ConvertFloat32ToInt32(out, in)
	default:
		in := AliasAs[float64](inRaw)
		ConvertFloat64ToInt32(out, in)
	}

	return samples, err
}

func (d *SampleTypeConverter) ReadFloat32(out []float32) (int, error) {

	if d.destType != SampleTypeFloat32 {
		return 0, errors.New("expected dest type to be float32")
	}

	// Read as many as len(out) samples from the input. The length of 'inRaw'
	// will always be an integer multiple of the 'sampleSize'.
	sampleSize := d.sourceType.Size()
	inRaw, err := d.reader.ReadBuffer(sampleSize * len(out))
	samples := len(inRaw) / sampleSize

	// The data in 'inRaw' will be interpreted according to the source type and
	// then converted to the appropriate dest type. The converted data will be
	// written into 'out'.
	switch d.sourceType {
	case SampleTypeUint8:
		in := AliasAs[uint8](inRaw)
		ConvertUint8ToFloat32(out, in)
	case SampleTypeInt16:
		in := AliasAs[int16](inRaw)
		ConvertInt16ToFloat32(out, in)
	case SampleTypeInt24:
		in := AliasAs[Int24](inRaw)
		ConvertInt24ToFloat32(out, in)
	case SampleTypeInt32:
		in := AliasAs[int32](inRaw)
		ConvertInt32ToFloat32(out, in)
	case SampleTypeFloat32:
		in := AliasAs[float32](inRaw)
		ConvertFloat32ToFloat32(out, in)
	default:
		in := AliasAs[float64](inRaw)
		ConvertFloat64ToFloat32(out, in)
	}

	return samples, err
}

func (d *SampleTypeConverter) ReadFloat64(out []float64) (int, error) {

	if d.destType != SampleTypeFloat64 {
		return 0, errors.New("expected dest type to be float64")
	}

	// Read as many as len(out) samples from the input. The length of 'inRaw'
	// will always be an integer multiple of the 'sampleSize'.
	sampleSize := d.sourceType.Size()
	inRaw, err := d.reader.ReadBuffer(sampleSize * len(out))
	samples := len(inRaw) / sampleSize

	// The data in 'inRaw' will be interpreted according to the source type and
	// then converted to the appropriate dest type. The converted data will be
	// written into 'out'.
	switch d.sourceType {
	case SampleTypeUint8:
		in := AliasAs[uint8](inRaw)
		ConvertUint8ToFloat64(out, in)
	case SampleTypeInt16:
		in := AliasAs[int16](inRaw)
		ConvertInt16ToFloat64(out, in)
	case SampleTypeInt24:
		in := AliasAs[Int24](inRaw)
		ConvertInt24ToFloat64(out, in)
	case SampleTypeInt32:
		in := AliasAs[int32](inRaw)
		ConvertInt32ToFloat64(out, in)
	case SampleTypeFloat32:
		in := AliasAs[float32](inRaw)
		ConvertFloat32ToFloat64(out, in)
	default:
		in := AliasAs[float64](inRaw)
		ConvertFloat64ToFloat64(out, in)
	}

	return samples, err
}
