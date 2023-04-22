package core

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInterleaveSlices(t *testing.T) {
	c1 := []float64{0, 3, 6, 9}
	c2 := []float64{1, 4, 7, 10}
	c3 := []float64{2, 5, 8, 11}

	expected := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	got, err := InterleaveSlices(c1, c2, c3)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestInterleaveSlicesInvalidInput(t *testing.T) {
	c1 := []float64{0, 3, 6, 8}
	c2 := []float64{1, 4, 7}
	c3 := []float64{2, 5}
	_, err := InterleaveSlices(c1, c2, c3)
	require.ErrorIs(t, err, ErrInterleaveInvalidElementCount)
}

func TestDeinterleaveSlices_Normal(t *testing.T) {
	input := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	expected := [][]float64{
		{0, 3, 6, 9},
		{1, 4, 7, 10},
		{2, 5, 8, 11},
	}

	output, err := DeinterleaveSlices(input, 3)
	require.NoError(t, err)
	require.Equal(t, expected, output)
}

func TestDeinterleaveSlices_InvalidInput(t *testing.T) {
	input := []float64{0, 1, 2, 3}
	_, err := DeinterleaveSlices(input, 3)
	require.ErrorIs(t, err, ErrDeinterleaveInvalidElementCount)
}
