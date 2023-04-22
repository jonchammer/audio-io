package core

// QuantizeToUint8 linearly maps input values in the range [-1, 1] to the range
// [0, 255], with input 0.0 mapping to output 128.
func QuantizeToUint8(input []float64) []uint8 {

	// [-1, 1] -> [0, 2] -> [0, 1] -> [0, 255]
	//   -> (x + 1) * 127.5 + 0.5
	//   -> (x * 127.5) + 128.0
	res := make([]uint8, len(input))
	for i := 0; i < len(input); i++ {
		res[i] = uint8((input[i] * 127.5) + 128.0)
	}
	return res
}

// QuantizeToInt16 linearly maps input values in the range [-1, 1] to the range
// [-32768, 32767], with input 0.0 mapping to output 0.
func QuantizeToInt16(input []float64) []int16 {

	// [-1, 1] -> [0, 2] -> [0, 65535] -> [-32768, 32767]
	//   -> (x + 1) * 32767.5 - 32768.0
	//   -> (x * 32767.5) + 32767.5 - 32768.0
	//   -> (x * 32767.5) - 0.5
	//
	// Note that 0.0 -> -0.5, but truncation towards 0 ensures that the result
	// is actually 0, as intended.
	res := make([]int16, len(input))
	for i := 0; i < len(input); i++ {
		res[i] = int16((input[i] * 32767.5) - 0.5)
	}
	return res
}

// QuantizeToInt24 linearly maps input values in the range [-1, 1] to the range
// [-8388608, 8388607], with input 0.0 mapping to output 0.
//
// Note that because int24 isn't a native type in Go, we'll use the larger type
// int32 as a container.
func QuantizeToInt24(input []float64) []int32 {

	// [-1, 1] -> [0, 2] -> [0, 16777215] -> [-8388608, 8388607]
	//   -> (x + 1) * 8388607.5 - 8388608.0
	//   -> (x * 8388607.5) + 8388607.5 - 8388608.0
	//   -> (x * 8388607.5) - 0.5
	//
	// Note that 0.0 -> -0.5, but truncation towards 0 ensures that the result
	// is actually 0, as intended.
	res := make([]int32, len(input))
	for i := 0; i < len(input); i++ {
		res[i] = int32((input[i] * 8388607.5) - 0.5)
	}
	return res
}

// QuantizeToInt32 linearly maps input values in the range [-1, 1] to the range
// [-2147483648, 2147483647], with input 0.0 mapping to output 0.
func QuantizeToInt32(input []float64) []int32 {

	// [-1, 1] -> [0, 2] -> [0, 4294967295] -> [-2147483648, 2147483647]
	//   -> (x + 1) * 2147483647.5 - 2147483648.0
	//   -> (x * 2147483647.5) + 2147483647.5 - 2147483648.0
	//   -> (x * 2147483647.5) - 0.5
	//
	// Note that 0.0 -> -0.5, but truncation towards 0 ensures that the result
	// is actually 0, as intended.
	res := make([]int32, len(input))
	for i := 0; i < len(input); i++ {
		res[i] = int32((input[i] * 2147483647.5) - 0.5)
	}
	return res
}

// QuantizeToFloat32 reduces the precision of the inputs from float64 to
// float32 by performing a direct cast on each element.
func QuantizeToFloat32(input []float64) []float32 {
	res := make([]float32, len(input))
	for i := 0; i < len(input); i++ {
		res[i] = float32(input[i])
	}
	return res
}
