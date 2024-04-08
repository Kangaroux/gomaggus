package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetBytes(t *testing.T) {
	expected := []byte{1, 2, 3}
	eb := EndianBytes{}

	eb.SetBytes(expected, true)
	assert.Equal(t, expected, eb.Bytes())

	// Data should have been cloned
	expected[0] = 2
	assert.NotEqual(t, expected, eb.Bytes())
}

func Test_EndianSwap(t *testing.T) {
	expectedBig := []byte{1, 2, 3}
	expectedLittle := []byte{3, 2, 1}
	eb := EndianBytes{}
	eb.SetBytes(expectedBig, true)

	eb.ToLittleEndian()
	assert.Equal(t, expectedLittle, eb.Bytes())

	// Little -> Little has no effect
	eb.ToLittleEndian()
	assert.Equal(t, expectedLittle, eb.Bytes())

	eb.ToBigEndian()
	assert.Equal(t, expectedBig, eb.Bytes())

	// Big -> Big has no effect
	eb.ToBigEndian()
	assert.Equal(t, expectedBig, eb.Bytes())
}

func Test_BytesFromHex(t *testing.T) {
	var b *EndianBytes
	var err error

	b, err = BytesFromHex("test", true)
	assert.Error(t, err)
	assert.Nil(t, b)

	b, err = BytesFromHex("010203", true)
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, b.Bytes())

	b.ToLittleEndian()
	assert.Equal(t, []byte{3, 2, 1}, b.Bytes())

	b, err = BytesFromHex("010203", false)
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, b.Bytes())

	b.ToBigEndian()
	assert.Equal(t, []byte{3, 2, 1}, b.Bytes())
}
