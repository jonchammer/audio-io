package core

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAliasAsUint8(t *testing.T) {

	// Verify that the alias shows the correct data
	src := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	aliased := AliasAs[uint8](src)
	require.Equal(
		t,
		[]uint8{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07},
		aliased,
	)

	// Changes to the aliased slice should be reflected in 'src'
	aliased[3] = 0xFF
	require.Equal(
		t,
		[]byte{0x00, 0x01, 0x02, 0xFF, 0x04, 0x05, 0x06, 0x07},
		src,
	)
}

func TestAliasAsInt16(t *testing.T) {

	// Verify that the alias shows the correct data
	src := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	aliased := AliasAs[int16](src)
	require.Equal(
		t,
		[]int16{0x0100, 0x0302, 0x0504, 0x0706},
		aliased,
	)

	// Changes to the aliased slice should be reflected in 'src'
	aliased[3] = 0x07FF
	require.Equal(
		t,
		[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0xFF, 0x07},
		src,
	)
}

func TestAliasAsInt32(t *testing.T) {

	// Verify that the alias shows the correct data
	src := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	aliased := AliasAs[int32](src)
	require.Equal(
		t,
		[]int32{0x03020100, 0x07060504},
		aliased,
	)

	// Changes to the aliased slice should be reflected in 'src'
	aliased[1] = 0x07FF0504
	require.Equal(
		t,
		[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0xFF, 0x07},
		src,
	)
}
