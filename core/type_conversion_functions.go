package core

import (
	"math"
)

// ------------------------------------------------------------------------- //
// X -> uint8 Converters
// ------------------------------------------------------------------------- //

func ConvertUint8ToUint8(out []uint8, in []uint8) {
	_ = copy(out, in)
}

func ConvertInt16ToUint8(out []uint8, in []int16) {
	// panic("int16 -> uint8 not implemented")
}

func ConvertInt24ToUint8(out []uint8, in []Int24) {

	// TODO: Revisit this. The implementation can probably be improved.

	for i := 0; i < len(in); i++ {
		x := in[i].AsInt32()
		sign := (x & MinInt24) >> 23
		divisor := float64(MaxInt24) - float64(sign)
		x2 := float64(x) / divisor
		out[i] = uint8((x2 * 127.5) + 128.0)
	}
}

func ConvertInt32ToUint8(out []uint8, in []int32) {
	// panic("int32 -> uint8 not implemented")
}

func ConvertFloat32ToUint8(out []uint8, in []float32) {

	// [-1, 1] -> [0, 2] -> [0, 1] -> [0, 255]
	//   -> (x + 1) * 127.5 + 0.5
	//   -> (x * 127.5) + 128.0
	for i := 0; i < len(in); i++ {
		out[i] = uint8((in[i] * 127.5) + 128.0)
	}
}

func ConvertFloat64ToUint8(out []uint8, in []float64) {

	// [-1, 1] -> [0, 2] -> [0, 1] -> [0, 255]
	//   -> (x + 1) * 127.5 + 0.5
	//   -> (x * 127.5) + 128.0
	for i := 0; i < len(in); i++ {
		out[i] = uint8((in[i] * 127.5) + 128.0)
	}
}

// ------------------------------------------------------------------------- //
// X -> int16 Converters
// ------------------------------------------------------------------------- //

func ConvertUint8ToInt16(out []int16, in []uint8) {
	// TODO: Revisit this. The implementation can probably be improved.

	m := [2]float64{255.0 / 32512.0, 1.0 / 127.0}
	b := [2]float64{-1.0, -128.0 / 127}

	for i := range in {
		idx := (in[i] & 0x80) >> 7
		dequantized := m[idx]*float64(in[i]) + b[idx]
		out[i] = int16((dequantized * 32767.5) - 0.5)
	}
}

func ConvertInt16ToInt16(out []int16, in []int16) {
	_ = copy(out, in)
}

func ConvertInt24ToInt16(out []int16, in []Int24) {

	// TODO: Revisit this. The implementation can probably be improved.

	for i := range in {
		x := in[i].AsInt32()
		sign := (x & MinInt24) >> 23
		divisor := float64(MaxInt24) - float64(sign)
		dequantized := float64(x) / divisor
		out[i] = int16((dequantized * 32767.5) - 0.5)
	}
}

func ConvertInt32ToInt16(out []int16, in []int32) {

	// TODO: Revisit this. The implementation can probably be improved.

	for i := range in {
		sign := (in[i] & math.MinInt32) >> 31
		divisor := float64(math.MaxInt32) - float64(sign)
		dequantized := float64(in[i]) / divisor
		out[i] = int16((dequantized * 32767.5) - 0.5)
	}
}

func ConvertFloat32ToInt16(out []int16, in []float32) {

	// [-1, 1] -> [0, 2] -> [0, 65535] -> [-32768, 32767]
	//   -> (x + 1) * 32767.5 - 32768.0
	//   -> (x * 32767.5) + 32767.5 - 32768.0
	//   -> (x * 32767.5) - 0.5
	//
	// Note that 0.0 -> -0.5, but truncation towards 0 ensures that the result
	// is actually 0, as intended.
	for i := range in {
		out[i] = int16((in[i] * 32767.5) - 0.5)
	}
}

func ConvertFloat64ToInt16(out []int16, in []float64) {

	// [-1, 1] -> [0, 2] -> [0, 65535] -> [-32768, 32767]
	//   -> (x + 1) * 32767.5 - 32768.0
	//   -> (x * 32767.5) + 32767.5 - 32768.0
	//   -> (x * 32767.5) - 0.5
	//
	// Note that 0.0 -> -0.5, but truncation towards 0 ensures that the result
	// is actually 0, as intended.
	for i := range in {
		out[i] = int16((in[i] * 32767.5) - 0.5)
	}
}

// ------------------------------------------------------------------------- //
// X -> int24 Converters
// ------------------------------------------------------------------------- //

func ConvertUint8ToInt24(out []Int24, in []uint8) {
	// panic("uint8 -> int24 not implemented")
}

func ConvertInt16ToInt24(out []Int24, in []int16) {
	// panic("int16 -> int24 not implemented")
}

func ConvertInt24ToInt24(out []Int24, in []Int24) {
	copy(out, in)
}

func ConvertInt32ToInt24(out []Int24, in []int32) {
	// panic("int32 -> int24 not implemented")
}

func ConvertFloat32ToInt24(out []Int24, in []float32) {

	// [-1, 1] -> [0, 2] -> [0, 16777215] -> [-8388608, 8388607]
	//   -> (x + 1) * 8388607.5 - 8388608.0
	//   -> (x * 8388607.5) + 8388607.5 - 8388608.0
	//   -> (x * 8388607.5) - 0.5
	//
	// Note that 0.0 -> -0.5, but truncation towards 0 ensures that the result
	// is actually 0, as intended.
	for i := 0; i < len(in); i++ {
		out[i] = Int24FromInt32(int32((in[i] * 8388607.5) - 0.5))
	}
}

func ConvertFloat64ToInt24(out []Int24, in []float64) {

	// [-1, 1] -> [0, 2] -> [0, 16777215] -> [-8388608, 8388607]
	//   -> (x + 1) * 8388607.5 - 8388608.0
	//   -> (x * 8388607.5) + 8388607.5 - 8388608.0
	//   -> (x * 8388607.5) - 0.5
	//
	// Note that 0.0 -> -0.5, but truncation towards 0 ensures that the result
	// is actually 0, as intended.
	for i := 0; i < len(in); i++ {
		out[i] = Int24FromInt32(int32((in[i] * 8388607.5) - 0.5))
	}
}

// ------------------------------------------------------------------------- //
// X -> int32 Converters
// ------------------------------------------------------------------------- //

func ConvertUint8ToInt32(out []int32, in []uint8) {
	// panic("uint8 -> int32 not implemented")
}

func ConvertInt16ToInt32(out []int32, in []int16) {
	// panic("int16 -> int32 not implemented")
}

func ConvertInt24ToInt32(out []int32, in []Int24) {
	// panic("int24 -> int32 not implemented")
}

func ConvertInt32ToInt32(out []int32, in []int32) {
	copy(out, in)
}

func ConvertFloat32ToInt32(out []int32, in []float32) {

	// [-1, 1] -> [0, 2] -> [0, 4294967295] -> [-2147483648, 2147483647]
	//   -> (x + 1) * 2147483647.5 - 2147483648.0
	//   -> (x * 2147483647.5) + 2147483647.5 - 2147483648.0
	//   -> (x * 2147483647.5) - 0.5
	//
	// Note that 0.0 -> -0.5, but truncation towards 0 ensures that the result
	// is actually 0, as intended.
	for i := 0; i < len(in); i++ {
		out[i] = int32((in[i] * 2147483647.5) - 0.5)
	}
}

func ConvertFloat64ToInt32(out []int32, in []float64) {

	// [-1, 1] -> [0, 2] -> [0, 4294967295] -> [-2147483648, 2147483647]
	//   -> (x + 1) * 2147483647.5 - 2147483648.0
	//   -> (x * 2147483647.5) + 2147483647.5 - 2147483648.0
	//   -> (x * 2147483647.5) - 0.5
	//
	// Note that 0.0 -> -0.5, but truncation towards 0 ensures that the result
	// is actually 0, as intended.
	for i := 0; i < len(in); i++ {
		out[i] = int32((in[i] * 2147483647.5) - 0.5)
	}
}

// ------------------------------------------------------------------------- //
// X -> float32 Converters
// ------------------------------------------------------------------------- //

func ConvertUint8ToFloat32(out []float32, in []uint8) {

	// Model parameters
	m := [2]float32{255.0 / 32512.0, 1.0 / 127.0}
	b := [2]float32{-1.0, -128.0 / 127}

	for i := 0; i < len(in); i++ {
		idx := (in[i] & 0x80) >> 7
		out[i] = m[idx]*(float32(in[i])) + b[idx]
	}
}

func ConvertInt16ToFloat32(out []float32, in []int16) {
	for i := 0; i < len(in); i++ {
		sign := (in[i] & math.MinInt16) >> 15
		divisor := float32(math.MaxInt16) - float32(sign)
		out[i] = float32(in[i]) / divisor
	}
}

func ConvertInt24ToFloat32(out []float32, in []Int24) {

	for i := 0; i < len(in); i++ {
		x := in[i].AsInt32()
		sign := (x & MinInt24) >> 23
		divisor := float32(MaxInt24) - float32(sign)
		out[i] = float32(x) / divisor
	}
}

func ConvertInt32ToFloat32(out []float32, in []int32) {
	for i := 0; i < len(in); i++ {
		sign := (in[i] & math.MinInt32) >> 31
		divisor := float32(math.MaxInt32) - float32(sign)
		out[i] = float32(in[i]) / divisor
	}
}

func ConvertFloat32ToFloat32(out []float32, in []float32) {
	copy(out, in)
}

func ConvertFloat64ToFloat32(out []float32, in []float64) {
	for i := 0; i < len(in); i++ {
		out[i] = float32(in[i])
	}
}

// ------------------------------------------------------------------------- //
// X -> float64 Converters
// ------------------------------------------------------------------------- //

func ConvertUint8ToFloat64(out []float64, in []uint8) {

	// This transformation cannot actually be performed linearly with the
	// constraint that 128 maps to 0.0. The basic mapping:
	//
	//   y = 2 * (x / 255) - 1
	//     = x * (2 / 255) - 1
	//
	// produces a value of y = 1.0 / 255.0 ~= 0.00392 for x = 128, rather than
	// 0.0, as expected.
	//
	// To counteract that, we'll use a "leaky rectifier" instead - that is, two
	// linear models that are applied in different parts of the domain. The
	// first model will include the points p0 = (0, -1) and p1 = (127, -1/256).
	// The second model will include the points p2 = (128, 0) and (255, 1).
	// After we work through the calculations, we arrive at the slope and bias
	// values below for each of the two models.
	//
	// As an implementation note, we can avoid the branch that switches between
	// the two models by checking the top bit of the input. A value of 0
	// implies that x is in the domain [0, 127]. A value of 1 implies that x
	// is in the domain [128, 255]. We'll use that top bit as an index into
	// arrays that hold the parameters for each model.

	// Model parameters
	m := [2]float64{255.0 / 32512.0, 1.0 / 127.0}
	b := [2]float64{-1.0, -128.0 / 127}

	for i := 0; i < len(in); i++ {
		idx := (in[i] & 0x80) >> 7
		out[i] = m[idx]*(float64(in[i])) + b[idx]
	}
}

func ConvertInt16ToFloat64(out []float64, in []int16) {

	// In order to guarantee the most accurate results, we'll start with two
	// different formulae: one for negative inputs and one for positive ones.
	//
	// Negative x:
	//   divisor = 32767 + 1 == 32767 - (-1)
	// Positive x:
	//   divisor = 32767 + 0 == 32767 - 0
	// result = x / divisor
	//
	// Practically, we can avoid the branch when calculating the divisor by
	// extracting the sign bit from the input and adding it to the divisor.
	// The actual implementation subtracts the sign bit (rather than adding)
	// because the sign bit will be interpreted as int32(-1) the way we
	// calculate it.

	for i := 0; i < len(in); i++ {
		sign := (in[i] & math.MinInt16) >> 15
		divisor := float64(math.MaxInt16) - float64(sign)
		out[i] = float64(in[i]) / divisor
	}
}

func ConvertInt24ToFloat64(out []float64, in []Int24) {

	// In order to guarantee the most accurate results, we'll start with two
	// different formulae: one for negative inputs and one for positive ones.
	//
	// Negative x:
	//   divisor = 8388607 + 1 == 8388607 - (-1)
	// Positive x:
	//   divisor = 8388607 + 0 == 8388607 - 0
	// result = x / divisor
	//
	// Practically, we can avoid the branch when calculating the divisor by
	// extracting the sign bit from the input and adding it to the divisor.
	// The actual implementation subtracts the sign bit (rather than adding)
	// because the sign bit will be interpreted as int32(-1) the way we
	// calculate it.

	for i := 0; i < len(in); i++ {
		x := in[i].AsInt32()
		sign := (x & MinInt24) >> 23
		divisor := float64(MaxInt24) - float64(sign)
		out[i] = float64(x) / divisor
	}
}

func ConvertInt32ToFloat64(out []float64, in []int32) {

	// In order to guarantee the most accurate results, we'll start with two
	// different formulae: one for negative inputs and one for positive ones.
	//
	// Negative x:
	//   divisor = 2147483647 + 1 == 2147483647 - (-1)
	// Positive x:
	//   divisor = 2147483647 + 0 == 2147483647 - 0
	// result = x / divisor
	//
	// Practically, we can avoid the branch when calculating the divisor by
	// extracting the sign bit from the input and adding it to the divisor.
	// The actual implementation subtracts the sign bit (rather than adding)
	// because the sign bit will be interpreted as int32(-1) the way we
	// calculate it.

	for i := 0; i < len(in); i++ {
		sign := (in[i] & math.MinInt32) >> 31
		divisor := float64(math.MaxInt32) - float64(sign)
		out[i] = float64(in[i]) / divisor
	}
}

func ConvertFloat32ToFloat64(out []float64, in []float32) {
	for i := 0; i < len(in); i++ {
		out[i] = float64(in[i])
	}
}

func ConvertFloat64ToFloat64(out []float64, in []float64) {
	copy(out, in)
}
