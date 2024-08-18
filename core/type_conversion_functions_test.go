package core

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestConversionFunctions(t *testing.T) {
	sampleTypes := []SampleType{
		SampleTypeUint8,
		SampleTypeInt16,
		SampleTypeInt24,
		SampleTypeInt32,
		SampleTypeFloat32,
		SampleTypeFloat64,
	}
	for _, source := range sampleTypes {
		for _, dest := range sampleTypes {
			name := fmt.Sprintf("Convert%sTo%s", source, dest)
			t.Run(name, func(t *testing.T) {
				testSingleConversionFunction(t, source, dest)
			})
		}
	}
}

func testSingleConversionFunction(
	t *testing.T, sourceType SampleType, destType SampleType,
) {

	// Min, 0, and Max values for each sample type
	testVals := map[SampleType]any{
		SampleTypeUint8: []uint8{0, 128, 255},
		SampleTypeInt16: []int16{-32768, 0, 32767},
		SampleTypeInt24: []Int24{
			Int24FromInt32(-8388608),
			Int24FromInt32(0),
			Int24FromInt32(8388607),
		},
		SampleTypeInt32:   []int32{-2147483648, 0, 2147483647},
		SampleTypeFloat32: []float32{-1.0, 0.0, 1.0},
		SampleTypeFloat64: []float64{-1.0, 0.0, 1.0},
	}

	// reflect.Type values for each sample type
	outTypes := map[SampleType]reflect.Type{
		SampleTypeUint8:   reflect.TypeFor[[]uint8](),
		SampleTypeInt16:   reflect.TypeFor[[]int16](),
		SampleTypeInt24:   reflect.TypeFor[[]Int24](),
		SampleTypeInt32:   reflect.TypeFor[[]int32](),
		SampleTypeFloat32: reflect.TypeFor[[]float32](),
		SampleTypeFloat64: reflect.TypeFor[[]float64](),
	}

	// All functions we're testing
	conversionFuncs := map[string]any{
		"ConvertUint8ToUint8":   ConvertUint8ToUint8,
		"ConvertInt16ToUint8":   ConvertInt16ToUint8,
		"ConvertInt24ToUint8":   ConvertInt24ToUint8,
		"ConvertInt32ToUint8":   ConvertInt32ToUint8,
		"ConvertFloat32ToUint8": ConvertFloat32ToUint8,
		"ConvertFloat64ToUint8": ConvertFloat64ToUint8,

		"ConvertUint8ToInt16":   ConvertUint8ToInt16,
		"ConvertInt16ToInt16":   ConvertInt16ToInt16,
		"ConvertInt24ToInt16":   ConvertInt24ToInt16,
		"ConvertInt32ToInt16":   ConvertInt32ToInt16,
		"ConvertFloat32ToInt16": ConvertFloat32ToInt16,
		"ConvertFloat64ToInt16": ConvertFloat64ToInt16,

		"ConvertUint8ToInt24":   ConvertUint8ToInt24,
		"ConvertInt16ToInt24":   ConvertInt16ToInt24,
		"ConvertInt24ToInt24":   ConvertInt24ToInt24,
		"ConvertInt32ToInt24":   ConvertInt32ToInt24,
		"ConvertFloat32ToInt24": ConvertFloat32ToInt24,
		"ConvertFloat64ToInt24": ConvertFloat64ToInt24,

		"ConvertUint8ToInt32":   ConvertUint8ToInt32,
		"ConvertInt16ToInt32":   ConvertInt16ToInt32,
		"ConvertInt24ToInt32":   ConvertInt24ToInt32,
		"ConvertInt32ToInt32":   ConvertInt32ToInt32,
		"ConvertFloat32ToInt32": ConvertFloat32ToInt32,
		"ConvertFloat64ToInt32": ConvertFloat64ToInt32,

		"ConvertUint8ToFloat32":   ConvertUint8ToFloat32,
		"ConvertInt16ToFloat32":   ConvertInt16ToFloat32,
		"ConvertInt24ToFloat32":   ConvertInt24ToFloat32,
		"ConvertInt32ToFloat32":   ConvertInt32ToFloat32,
		"ConvertFloat32ToFloat32": ConvertFloat32ToFloat32,
		"ConvertFloat64ToFloat32": ConvertFloat64ToFloat32,

		"ConvertUint8ToFloat64":   ConvertUint8ToFloat64,
		"ConvertInt16ToFloat64":   ConvertInt16ToFloat64,
		"ConvertInt24ToFloat64":   ConvertInt24ToFloat64,
		"ConvertInt32ToFloat64":   ConvertInt32ToFloat64,
		"ConvertFloat32ToFloat64": ConvertFloat32ToFloat64,
		"ConvertFloat64ToFloat64": ConvertFloat64ToFloat64,
	}

	// Look up the 'in' and 'expected' slices based on the source and dest
	// sample types
	in := testVals[sourceType]
	expected := testVals[destType]

	// Create reflect.Value objects for the function inputs (the 'in' and 'out'
	// slices). The out slice doesn't yet exist, so we'll create it dynamically
	// based on the dest sample type.
	inVal := reflect.ValueOf(in)
	outVal := reflect.MakeSlice(outTypes[destType], inVal.Len(), inVal.Cap())
	out := outVal.Interface()

	// Use the 'reflect' library to call the proper type conversion function
	// by name
	fnName := fmt.Sprintf("Convert%sTo%s", sourceType, destType)
	f := reflect.ValueOf(conversionFuncs[fnName])
	_ = f.Call([]reflect.Value{outVal, inVal})

	// Test assertions
	require.Equal(t, expected, out)
}
