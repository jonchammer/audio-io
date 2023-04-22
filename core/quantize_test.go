package core

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQuantizeToUint8(t *testing.T) {
	input := []float64{-1.0, 0.0, +1.0}
	expected := []uint8{0, 128, 255}
	require.Equal(t, expected, QuantizeToUint8(input))
}

func TestQuantizeToInt16(t *testing.T) {
	input := []float64{-1.0, 0.0, +1.0}
	expected := []int16{-32768, 0, 32767}
	require.Equal(t, expected, QuantizeToInt16(input))
}

func TestQuantizeToInt24(t *testing.T) {
	input := []float64{-1.0, 0.0, +1.0}
	expected := []int32{-8388608, 0, 8388607}
	require.Equal(t, expected, QuantizeToInt24(input))
}

func TestQuantizeToInt32(t *testing.T) {
	input := []float64{-1.0, 0.0, +1.0}
	expected := []int32{-2147483648, 0, 2147483647}
	require.Equal(t, expected, QuantizeToInt32(input))
}

func TestQuantizeToFloat32(t *testing.T) {
	input := []float64{-1.0, 0.0, +1.0}
	expected := []float32{-1.0, 0.0, +1.0}
	require.Equal(t, expected, QuantizeToFloat32(input))
}
