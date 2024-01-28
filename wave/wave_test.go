package wave

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/jonchammer/audio-io/core"
)

// ------------------------------------------------------------------------- //
// FormatCode
// ------------------------------------------------------------------------- //

func TestFormatCode_IsValid(t *testing.T) {
	require.True(t, FormatCodePCM.IsValid())
	require.False(t, FormatCode(99).IsValid())
}

func TestFormatCode_String(t *testing.T) {
	require.Equal(t, "PCM", FormatCodePCM.String())
	require.Equal(t, "IEEE Float", FormatCodeIEEEFloat.String())
	require.Equal(t, "Extensible", FormatCodeExtensible.String())
	require.Equal(t, "FormatCode(99)", FormatCode(99).String())
}

func TestSampleType_EffectiveFormatCode(t *testing.T) {
	require.Equal(t, FormatCodePCM, effectiveFormatCode(core.SampleTypeUint8))
	require.Equal(t, FormatCodePCM, effectiveFormatCode(core.SampleTypeInt16))
	require.Equal(t, FormatCodePCM, effectiveFormatCode(core.SampleTypeInt24))
	require.Equal(t, FormatCodePCM, effectiveFormatCode(core.SampleTypeInt32))
	require.Equal(t, FormatCodeIEEEFloat, effectiveFormatCode(core.SampleTypeFloat32))
	require.Equal(t, FormatCodeIEEEFloat, effectiveFormatCode(core.SampleTypeFloat64))
}
