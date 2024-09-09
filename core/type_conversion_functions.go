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

	// NOTE: Extract the most-significant byte of each sample and add 128
	//
	// -32768,   0, 32767 -> 0x8000, 0x0000, 0x7FFF
	//   -128,   0,   128 -> 0xFF80, 0x0000, 0x007F
	//      0, 128,   255  > 0x0000, 0x007F, 0x00FF

	for i := 0; i < len(in); i++ {
		out[i] = uint8(((in[i] >> 8) & 0xFF) + 128)
	}
}

func ConvertInt24ToUint8(out []uint8, in []Int24) {

	// NOTE: Extract the most-significant byte (Int24[2] in little-endian
	// notation) and add 128

	for i := 0; i < len(in); i++ {
		out[i] = in[i][2] + 128
	}
}

func ConvertInt32ToUint8(out []uint8, in []int32) {

	// NOTE: Extract the most-significant byte and add 128

	for i := 0; i < len(in); i++ {
		out[i] = uint8(((in[i] >> 24) & 0xFF) + 128)
	}
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

	// y = (   256 * x ) / 1   {x  < 0}
	// y = ( 25801 * x ) / 100 {x >= 0}

	// Model parameters - slope and divisor
	m := [2]int32{25801, 256}
	d := [2]int32{100, 1}

	for i := 0; i < len(in); i++ {
		normalized := int32(in[i]) - 128
		idx := (int32(normalized) & 0x80) >> 7
		out[i] = int16((normalized * m[idx]) / d[idx])
	}
}

func ConvertInt16ToInt16(out []int16, in []int16) {
	_ = copy(out, in)
}

func ConvertInt24ToInt16(out []int16, in []Int24) {
	for i := range in {
		out[i] = int16(in[i].AsInt32() / 256)
	}
}

func ConvertInt32ToInt16(out []int16, in []int32) {
	for i := range in {
		out[i] = int16(in[i] / 65536)
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

	// We'll first normalize the input samples by subtracting 128. That will
	// get us values in the range [-128, 127].
	//        0, 128,     255  > 0x0000 0000, 0x0000 007F, 0x0000 00FF
	//     -128,   0,     127  > 0xFFFF FF80, 0x0000 0000, 0x0000 007F
	// -8388608,   0, 8388607  > 0xFF80 0000, 0x0000 0000, 0x007F FFFF
	//
	// We'll use two different conversion functions - one for the negative part
	// of the domain and one for the positive part.
	//   Negative: y = ( 65536 * x ) / 1
	//   Positive: y = ( 6605203 * x ) / 100
	//
	// We'll switch between the two equations based on the sign bit of the
	// normalized input value.

	// Model parameters - slope and divisor
	m := [2]int32{6605203, 65536}
	d := [2]int32{100, 1}

	for i := 0; i < len(in); i++ {
		normalized := int32(in[i]) - 128
		idx := (normalized & 0x80) >> 7
		out[i] = Int24FromInt32((normalized * m[idx]) / d[idx])
	}
}

func ConvertInt16ToInt24(out []Int24, in []int16) {
	// y = (     256 * x ) / 1     {x  < 0}
	// y = ( 2560078 * x ) / 10000 {x >= 0}

	// Model parameters - slope and divisor
	m := [2]int64{2560078, 256}
	d := [2]int64{10000, 1}

	for i := 0; i < len(in); i++ {
		idx := (int64(in[i]) & 0x8000) >> 15
		out[i] = Int24FromInt64((int64(in[i]) * m[idx]) / d[idx])
	}
}

func ConvertInt24ToInt24(out []Int24, in []Int24) {
	copy(out, in)
}

func ConvertInt32ToInt24(out []Int24, in []int32) {
	for i := 0; i < len(in); i++ {
		out[i] = Int24FromInt32(in[i] / 256)
	}
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

	//        -128 > 0xFFFF FFFF FFFF FF80
	//           0 > 0x0000 0000 0000 0000
	//         127 > 0x0000 0000 0000 007F
	//         255 > 0x0000 0000 0000 00FF
	// -------------------------
	// -2147483648 > 0xFFFF FFFF 8000 0000
	//  2147483647 > 0x0000 0000 7FFF FFFF

	// We'll first normalize the input samples by subtracting 128. That will
	// get us values in the range [-128, 127].
	//
	// We'll then use two different conversion functions - one for the negative
	// part of the domain and one for the positive part.
	//   Negative: y = ( 16,777,216 * x ) / 1
	//   Positive: y = ( 1,690,932,006 * x ) / 100
	//
	// We'll switch between the two equations based on the sign bit of the
	// normalized input value.

	// Model parameters - slope and divisor
	m := [2]int64{1690932006, 16777216}
	d := [2]int64{100, 1}

	for i := 0; i < len(in); i++ {
		normalized := int64(in[i]) - 128
		idx := (normalized & 0x80) >> 7
		out[i] = int32((normalized * m[idx]) / d[idx])
	}
}

func ConvertInt16ToInt32(out []int32, in []int16) {

	// y = (      65536 * x ) / 1      {x  < 0}
	// y = ( 6553800004 * x ) / 100000 {x >= 0}

	// Model parameters - slope and divisor
	m := [2]int64{6553800004, 65536}
	d := [2]int64{100000, 1}

	for i := 0; i < len(in); i++ {
		idx := (int64(in[i]) & 0x8000) >> 15
		out[i] = int32((int64(in[i]) * m[idx]) / d[idx])
	}
}

func ConvertInt24ToInt32(out []int32, in []Int24) {
	// y = (        256 * x ) / 1        {x  < 0}
	// y = ( 2560000304 * x ) / 10000000 {x >= 0}

	// Model parameters - slope and divisor
	m := [2]int64{2560000304, 256}
	d := [2]int64{10000000, 1}

	for i := 0; i < len(in); i++ {
		idx := (in[i][2] & 0x80) >> 7
		out[i] = int32((in[i].AsInt64() * m[idx]) / d[idx])
	}
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
