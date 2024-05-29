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

func TestX(t *testing.T) {
	expected := decodeHex("D927E98BE3E9AF84FDC99DE9034F8E70ED7E90D6")
	salt := reverse(decodeHex("CAC94AF32D817BA64B13F18FDEDEF92AD4ED7EF7AB0E19E9F2AE13C828AEAF57"))
	username := "USERNAME123"
	password := "PASSWORD123"
	assert.Equal(t, expected, calculateX(username, password, salt))

	expected = decodeHex("E2F9A0F1E824006C98DA753448E743F7DAA1EAA1")
	username = "00XD0QOSA9L8KMXC"
	password = "43R4Z35TKBKFW8JI"
	assert.Equal(t, expected, calculateX(username, password, salt))
}

func TestVerifier(t *testing.T) {
	expected := decodeHex("21B4153B0A938D0A69D28F2690CC3F79A99A13C40CACB525B3B79D4201EB33FF")
	salt := reverse(decodeHex("AFE5D28E925DBB3DAFED5D91ACA0928940E8FBFEF2D2A3CC154ADA0FE6ABEF6F"))
	username := "LF2BGFQIFQ3HZ1ZF"
	password := "MVRVMUJFWRA0IBVK"
	assert.Equal(t, expected, calculateVerifier(username, password, salt))
}

func TestServerPublicKey(t *testing.T) {
	verifier := decodeHex("870A98A3DA8CCAFE6B2F4B0C43A022A0C6CEF4374BA4A50CEBF3FACA60237DC4")
	privateKey := decodeHex("ACDCB7CB1DE67DB1D5E0A37DAE80068BCCE062AE0EDA0CBEADF560BCDAE6D6B9")
	expected := decodeHex("85A204C987B68764FA69C523E32B940D1E1822B9E0F134FDC5086B1408A2BB43")
	assert.Equal(t, expected, calculateServerPublicKey(verifier, privateKey))
}
