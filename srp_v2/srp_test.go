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
	assert.Equal(t, expected, CalculateX(username, password, salt))

	expected = decodeHex("E2F9A0F1E824006C98DA753448E743F7DAA1EAA1")
	username = "00XD0QOSA9L8KMXC"
	password = "43R4Z35TKBKFW8JI"
	assert.Equal(t, expected, CalculateX(username, password, salt))
}

func TestVerifier(t *testing.T) {
	expected := decodeHex("21B4153B0A938D0A69D28F2690CC3F79A99A13C40CACB525B3B79D4201EB33FF")
	salt := reverse(decodeHex("AFE5D28E925DBB3DAFED5D91ACA0928940E8FBFEF2D2A3CC154ADA0FE6ABEF6F"))
	username := "LF2BGFQIFQ3HZ1ZF"
	password := "MVRVMUJFWRA0IBVK"
	assert.Equal(t, expected, CalculateVerifier(username, password, salt))
}

func TestServerPublicKey(t *testing.T) {
	expected := decodeHex("85A204C987B68764FA69C523E32B940D1E1822B9E0F134FDC5086B1408A2BB43")
	verifier := decodeHex("870A98A3DA8CCAFE6B2F4B0C43A022A0C6CEF4374BA4A50CEBF3FACA60237DC4")
	privateKey := decodeHex("ACDCB7CB1DE67DB1D5E0A37DAE80068BCCE062AE0EDA0CBEADF560BCDAE6D6B9")
	assert.Equal(t, expected, CalculateServerPublicKey(verifier, privateKey))
}

func TestCalculateU(t *testing.T) {
	expected := decodeHex("1309BD7851A1A505B95D6F60A8D884133458D24E")
	clientPublic := decodeHex("6FCEEEE7D40AAF0C7A08DFE1EFD3FCE80A152AA436CECB77FC06DAF9E9E5BDF3")
	serverPublic := decodeHex("F8CD769BDE603FC8F48B9BE7C5BEAAA7BD597ABDBDAC1AEFCACF0EE13443A3B9")
	assert.Equal(t, expected, CalculateU(clientPublic, serverPublic))
}

func TestServerSKey(t *testing.T) {
	expected := reverse(decodeHex("3503B289A60D6DD59EBD6FD88DF24836833433E39048ECAFF7E887313554F85C"))
	clientPublic := decodeHex("51CCDDFACF7F960EDF5030F09F0B033C0D08DB1E43FCBA3A92ABB4BE3535D1DB")
	verifier := decodeHex("6FC7D4ACFCFFFDCF780EE9BBD17AE507FFCDF586F83B2C9AEE2198F195DB3AB5")
	u := decodeHex("F9CEDDD82E776BEDB1A94852A9A7FFA4FCADD5DE")
	serverPrivate := decodeHex("A5DBBFCB4C7A1B7C3041CAC9DDBD36CD646F9FBABDAD66A019BCBB8FEDF2FAAE")
	assert.Equal(t, expected, CalculateServerSKey(clientPublic, verifier, u, serverPrivate))
}

func TestInterleave(t *testing.T) {
	expected := decodeHex("EE144E1AE08DAC891AB63ABC42BF89738003343422E6B58131BEE4C3087A7027E55A7216D18D556C")
	S := decodeHex("8F4CEBD60DFC34E5C007E51BD4F3A4FF2BC1D930E2D3EA770D8D3EEDFF2DCCFC")
	assert.Equal(t, expected, CalculateInterleave(S))
}

func TestClientProof(t *testing.T) {
	expected := reverse(decodeHex("7D07022B4064CCE633D679F61C6B212B6F8BC5C3"))
	username := "7WG6SHZL33JMGPO4"
	salt := reverse(decodeHex("00A4A09E0B5ACA438B8CD837D0816CA26043DBD1EAEF138EEF72DCF3F696D03D"))
	clientPublic := decodeHex("0095FE039AFE5E1BADE9AC0CAEC3CB73D2D08BBF4CA8ADDBCDF0CE709ED5103F")
	serverPublic := decodeHex("00B0C41F58CCE894CFB816FA72CA344C9FE2ED7CE799452ADBA7ABDCD26EAE75")
	sessionKey := decodeHex("77A4D39CF9C0BF373EF870BD2941C339C575FDD1CBAA31C919EA7BD5023267D303E20FEC9A9C402F")
	assert.Equal(t, expected, CalculateClientProof(username, salt, clientPublic, serverPublic, sessionKey))
}
