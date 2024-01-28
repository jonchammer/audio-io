package core

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// ------------------------------------------------------------------------- //
// SampleType
// ------------------------------------------------------------------------- //

func TestSampleType_IsValid(t *testing.T) {
	require.True(t, SampleTypeUint8.IsValid())
	require.True(t, SampleTypeFloat64.IsValid())
	require.False(t, SampleType(99).IsValid())
}

func TestSampleType_Size(t *testing.T) {
	require.Equal(t, 1, SampleTypeUint8.Size())
	require.Equal(t, 2, SampleTypeInt16.Size())
	require.Equal(t, 3, SampleTypeInt24.Size())
	require.Equal(t, 4, SampleTypeInt32.Size())
	require.Equal(t, 4, SampleTypeFloat32.Size())
	require.Equal(t, 8, SampleTypeFloat64.Size())
}

func TestSampleType_String(t *testing.T) {
	require.Equal(t, "Uint8", SampleTypeUint8.String())
	require.Equal(t, "Int16", SampleTypeInt16.String())
	require.Equal(t, "Int24", SampleTypeInt24.String())
	require.Equal(t, "Int32", SampleTypeInt32.String())
	require.Equal(t, "Float32", SampleTypeFloat32.String())
	require.Equal(t, "Float64", SampleTypeFloat64.String())
	require.Equal(t, "SampleType(99)", SampleType(99).String())
}
