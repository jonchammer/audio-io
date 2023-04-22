package core

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDequantizeUint8(t *testing.T) {
	input := []uint8{0, 128, 255}
	expected := []float64{-1.0, 0.0, +1.0}
	require.Equal(t, expected, DequantizeUint8(input))
}

func TestDequantizeInt16(t *testing.T) {
	input := []int16{-32768, 0, 32767}
	expected := []float64{-1.0, 0.0, 1.0}
	require.Equal(t, expected, DequantizeInt16(input))
}

func TestDequantizeInt24(t *testing.T) {
	input := []int32{-8388608, 0, 8388607}
	expected := []float64{-1.0, 0.0, 1.0}
	require.Equal(t, expected, DequantizeInt24(input))
}

func TestDequantizeInt32(t *testing.T) {
	input := []int32{-2147483648, 0, 2147483647}
	expected := []float64{-1.0, 0.0, 1.0}
	require.Equal(t, expected, DequantizeInt32(input))
}

func TestDequantizeFloat32(t *testing.T) {
	input := []float32{-1.0, 0.0, 1.0}
	expected := []float64{-1.0, 0.0, 1.0}
	require.Equal(t, expected, DequantizeFloat32(input))
}
