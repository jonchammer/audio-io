package core

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInt24_Normal(t *testing.T) {
	x := Int24FromInt32(8388607)
	require.Equal(t, int32(8388607), x.AsInt32())

	x = Int24FromInt32(0)
	require.Equal(t, int32(0), x.AsInt32())

	x = Int24FromInt32(-8388608)
	require.Equal(t, int32(-8388608), x.AsInt32())
}

func TestInt24_Slice(t *testing.T) {
	data := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
	}

	aliasedData := AliasAs[Int24](data)
	require.Equal(t, []Int24{
		[3]byte{0x01, 0x02, 0x03},
		[3]byte{0x04, 0x05, 0x06},
		[3]byte{0x07, 0x08, 0x09},
	}, aliasedData)

	aliasedData[2] = Int24FromInt32(8388607)
	require.Equal(t, []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0xFF, 0xFF, 0x7F,
	}, data)
}
