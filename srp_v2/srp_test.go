package srpv2

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func decodeHex(s string) []byte {
	val, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return val
}

func TestReverse(t *testing.T) {
	assert.Equal(t, []byte{0, 1}, reverse([]byte{1, 0}))
	assert.Equal(t, []byte{0, 1, 2}, reverse([]byte{2, 1, 0}))
}

func TestX(t *testing.T) {
	expected := reverse(decodeHex("D927E98BE3E9AF84FDC99DE9034F8E70ED7E90D6"))
	salt := reverse(decodeHex("CAC94AF32D817BA64B13F18FDEDEF92AD4ED7EF7AB0E19E9F2AE13C828AEAF57"))
	username := "USERNAME123"
	password := "PASSWORD123"
	assert.Equal(t, expected, calculateX(username, password, salt))
}

func TestVerifier(t *testing.T) {
	expected := reverse(decodeHex("21B4153B0A938D0A69D28F2690CC3F79A99A13C40CACB525B3B79D4201EB33FF"))
	salt := reverse(decodeHex("AFE5D28E925DBB3DAFED5D91ACA0928940E8FBFEF2D2A3CC154ADA0FE6ABEF6F"))
	username := "LF2BGFQIFQ3HZ1ZF"
	password := "MVRVMUJFWRA0IBVK"
	assert.Equal(t, expected, calculateVerifier(username, password, salt))
}
