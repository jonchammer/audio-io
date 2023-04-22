package core

import (
	"math"
)

// DequantizeUint8 maps input values in the range [0, 255] to the range
// [-1.0, 1.0], with input 128 mapping to 0.0.
func DequantizeUint8(input []uint8) []float64 {

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

	res := make([]float64, len(input))
	for i := 0; i < len(input); i++ {
		idx := (input[i] & 0x80) >> 7
		res[i] = m[idx]*(float64(input[i])) + b[idx]
	}
	return res
}

// DequantizeInt16 maps input values in the range [-32768, 32767] to the range
// [-1.0, 1.0], with input 0 mapping to 0.0.
func DequantizeInt16(input []int16) []float64 {

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

	res := make([]float64, len(input))
	for i := 0; i < len(input); i++ {
		sign := (input[i] & math.MinInt16) >> 15
		divisor := float64(math.MaxInt16) - float64(sign)
		res[i] = float64(input[i]) / divisor
	}
	return res
}

// DequantizeInt24 maps input values in the range [-8388608, 8388607] to the
// range [-1.0, 1.0], with input 0 mapping to 0.0.
func DequantizeInt24(input []int32) []float64 {

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

	const (
		minInt24 = -1 << 23
		maxInt24 = 1<<23 - 1
	)

	res := make([]float64, len(input))
	for i := 0; i < len(input); i++ {
		sign := (input[i] & minInt24) >> 23
		divisor := float64(maxInt24) - float64(sign)
		res[i] = float64(input[i]) / divisor
	}
	return res
}

// DequantizeInt32 maps input values in the range [-2147483648, 2147483647] to
// the range [-1.0, 1.0], with input 0 mapping to 0.0.
func DequantizeInt32(input []int32) []float64 {

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

	res := make([]float64, len(input))
	for i := 0; i < len(input); i++ {
		sign := (input[i] & math.MinInt32) >> 31
		divisor := float64(math.MaxInt32) - float64(sign)
		res[i] = float64(input[i]) / divisor
	}
	return res
}

// DequantizeFloat32 casts each input value from a float32 to a float64.
func DequantizeFloat32(input []float32) []float64 {
	res := make([]float64, len(input))
	for i := 0; i < len(input); i++ {
		res[i] = float64(input[i])
	}
	return res
}
