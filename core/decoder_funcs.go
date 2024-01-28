package core

import (
	"encoding/binary"
)

type key struct {
	sourceType SampleType
	destType   SampleType
}

type decoderFunc func(AudioReader, []byte) (int, error)

var decoderFuncs = map[key]decoderFunc{
	{SampleTypeUint8, SampleTypeUint8}: readUint8ToUint8,
	{SampleTypeUint8, SampleTypeInt16}: readUint8ToInt16,
	// {SampleTypeUint8, SampleTypeInt24}:   readUint8ToInt24,
	{SampleTypeUint8, SampleTypeInt32}: readUint8ToInt32,
	// {SampleTypeUint8, SampleTypeFloat32}: readUint8ToFloat32,
	// {SampleTypeUint8, SampleTypeFloat64}: readUint8ToFloat64,

	// {SampleTypeInt16, SampleTypeUint8}:   readInt16ToUint8,
	{SampleTypeInt16, SampleTypeInt16}: readInt16ToInt16,
	// {SampleTypeInt16, SampleTypeInt24}:   readInt16ToInt24,
	// {SampleTypeInt16, SampleTypeInt32}:   readInt16ToInt32,
	// {SampleTypeInt16, SampleTypeFloat32}: readInt16ToFloat32,
	// {SampleTypeInt16, SampleTypeFloat64}: readInt16ToFloat64,

	{SampleTypeInt24, SampleTypeUint8}: readInt24ToUint8,
	{SampleTypeInt24, SampleTypeInt16}: readInt24ToInt16,
	// {SampleTypeInt24, SampleTypeInt24}:   readInt24ToInt24,
	{SampleTypeInt24, SampleTypeInt32}: readInt24ToInt32,
	// {SampleTypeInt24, SampleTypeFloat32}: readInt24ToFloat32,
	// {SampleTypeInt24, SampleTypeFloat64}: readInt24ToFloat64,

	// {SampleTypeInt32, SampleTypeUint8}:   readInt32ToUint8,
	// {SampleTypeInt32, SampleTypeInt16}:   readInt32ToInt16,
	// {SampleTypeInt32, SampleTypeInt24}:   readInt32ToInt24,
	// {SampleTypeInt32, SampleTypeInt32}:   readInt32ToInt32,
	// {SampleTypeInt32, SampleTypeFloat32}: readInt32ToFloat32,
	// {SampleTypeInt32, SampleTypeFloat64}: readInt32ToFloat64,

	// {SampleTypeFloat32, SampleTypeUint8}:   readFloat32ToUint8,
	// {SampleTypeFloat32, SampleTypeInt16}:   readFloat32ToInt16,
	// {SampleTypeFloat32, SampleTypeInt24}:   readFloat32ToInt24,
	// {SampleTypeFloat32, SampleTypeInt32}:   readFloat32ToInt32,
	// {SampleTypeFloat32, SampleTypeFloat32}: readFloat32ToFloat32,
	// {SampleTypeFloat32, SampleTypeFloat64}: readFloat32ToFloat64,

	// {SampleTypeFloat64, SampleTypeUint8}:   readFloat64ToUint8,
	// {SampleTypeFloat64, SampleTypeInt16}:   readFloat64ToInt16,
	// {SampleTypeFloat64, SampleTypeInt24}:   readFloat64ToInt24,
	// {SampleTypeFloat64, SampleTypeInt32}:   readFloat64ToInt32,
	// {SampleTypeFloat64, SampleTypeFloat32}: readFloat64ToFloat32,
	// {SampleTypeFloat64, SampleTypeFloat64}: readFloat64ToFloat64,
}

// ------------------------------------------------------------------------- //
// Uint8
// ------------------------------------------------------------------------- //

func readUint8ToUint8(r AudioReader, p []byte) (int, error) {
	return r.ReadUint8(p)
}

func readUint8ToInt16(r AudioReader, p []byte) (int, error) {

	// 'p' will be interpreted as int16
	n := 2
	maxSamples := len(p) / n

	data := make([]uint8, maxSamples)
	samplesRead, err := r.ReadUint8(data)
	if err != nil {
		return 0, err
	}

	int16Data := QuantizeToInt16(DequantizeUint8(data))
	for i, element := range int16Data {
		binary.LittleEndian.PutUint16(p[n*i:], (uint16)(element))
	}

	return samplesRead * n, nil
}

// func readUint8ToInt24(r AudioReader, p []byte) (int, error) {
//
// 	// 'p' will be interpreted as int24
// 	n := 3
// 	maxSamples := len(p) / n
//
// 	data := make([]uint8, maxSamples)
// 	samplesRead, err := r.ReadUint8(data)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	floatData := DequantizeUint8(data)
// 	int24Data := QuantizeToInt24(floatData)
// 	for i, element := range int24Data {
// 		binary.LittleEndian.PutUint16(p[n*i:], (uint16)(element))
// 	}
//
// 	return samplesRead * n, nil
// }

func readUint8ToInt32(r AudioReader, p []byte) (int, error) {

	// 'p' will be interpreted as int32
	n := 4
	maxSamples := len(p) / n

	data := make([]uint8, maxSamples)
	samplesRead, err := r.ReadUint8(data)
	if err != nil {
		return 0, err
	}

	int32Data := QuantizeToInt32(DequantizeUint8(data))
	for i, element := range int32Data {
		binary.LittleEndian.PutUint32(p[n*i:], (uint32)(element))
	}

	return samplesRead * n, nil
}

// ------------------------------------------------------------------------- //
// Int16
// ------------------------------------------------------------------------- //

func readInt16ToInt16(r AudioReader, p []byte) (int, error) {

	// 'p' will be interpreted as int16
	n := 2
	maxSamples := len(p) / n

	data := make([]int16, maxSamples)
	samplesRead, err := r.ReadInt16(data)
	if err != nil {
		return 0, err
	}

	for i, element := range data {
		binary.LittleEndian.PutUint16(p[n*i:], (uint16)(element))
	}

	return samplesRead * n, nil
}

// ------------------------------------------------------------------------- //
// Int24
// ------------------------------------------------------------------------- //

func readInt24ToUint8(r AudioReader, p []byte) (int, error) {

	maxSamples := len(p)
	data := make([]int32, maxSamples)
	samplesRead, err := r.ReadInt24(data)
	if err != nil {
		return 0, err
	}

	// TODO: QuantizeToUint8 should be able to directly write to 'p'
	uint8Data := QuantizeToUint8(DequantizeInt24(data))
	for i, element := range uint8Data {
		p[i] = element
	}

	return samplesRead, nil
}

func readInt24ToInt16(r AudioReader, p []byte) (int, error) {

	// 'p' will be interpreted as int16
	n := 2
	maxSamples := len(p) / n

	data := make([]int32, maxSamples)
	samplesRead, err := r.ReadInt24(data)
	if err != nil {
		return 0, err
	}

	int16Data := QuantizeToInt16(DequantizeInt24(data))
	for i, element := range int16Data {
		binary.LittleEndian.PutUint16(p[n*i:], (uint16)(element))
	}

	return samplesRead * n, nil
}

func readInt24ToInt32(r AudioReader, p []byte) (int, error) {

	// 'p' will be interpreted as int32
	n := 4
	maxSamples := len(p) / n

	data := make([]int32, maxSamples)
	samplesRead, err := r.ReadInt24(data)
	if err != nil {
		return 0, err
	}

	int32Data := QuantizeToInt32(DequantizeInt24(data))
	for i, element := range int32Data {
		binary.LittleEndian.PutUint32(p[n*i:], (uint32)(element))
	}

	return samplesRead * n, nil
}

// ------------------------------------------------------------------------- //
// Int32
// ------------------------------------------------------------------------- //

// ------------------------------------------------------------------------- //
// Float32
// ------------------------------------------------------------------------- //

// ------------------------------------------------------------------------- //
// Float64
// ------------------------------------------------------------------------- //
