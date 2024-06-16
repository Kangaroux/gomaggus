package realmd

import (
	"bytes"
	"crypto/rc4"
	"testing"

	"github.com/kangaroux/gomaggus/internal"
	"github.com/stretchr/testify/assert"
)

func TestGenerateKey(t *testing.T) {
	rows := internal.MustLoadTestData("../testdata/header/wrath_generate_key.csv")

	for _, row := range rows {
		expected := internal.MustDecodeHex(row[2])
		sessionKey := internal.MustDecodeHex(row[0])
		fixedKey := internal.MustDecodeHex(row[1])
		wrath := &HeaderCrypto{sessionKey: sessionKey}
		assert.Equal(t, expected, wrath.GenerateKey(fixedKey))
	}
}
func TestDrop1024(t *testing.T) {
	key := internal.MustDecodeHex("470446575F1EEDC1732473F4BDCACF4C13EE3837A3E2B9B720472E851525855DFC8045BF80D2F4D5")
	data := []byte("hello world")
	result := make([]byte, len(data))

	// Test the output when no bytes are dropped
	c, _ := rc4.NewCipher(key)
	copy(result, data)
	c.XORKeyStream(result, data)
	assert.Equal(t, internal.MustDecodeHex("B53577EC1E7FEAE6CDBB5C"), result)

	// Test the output when the first 1024 bytes are dropped (RC4-drop1024)
	c, _ = rc4.NewCipher(key)
	drop1024(c)
	copy(result, data)
	c.XORKeyStream(result, data)
	assert.Equal(t, internal.MustDecodeHex("9F273843E006E43B0D33B6"), result)
}

func TestEncryptDecrypt(t *testing.T) {
	sessionKey := internal.Reverse(internal.MustDecodeHex("403FCE7B2B1FCDAE43F118B6C7E517D5A1498088180936A3E45B9888978B7675ECBAA7DB4CA4E8DE"))
	data := internal.Reverse(internal.MustDecodeHex("DEBA4C8DCBD613F06E725123E887CF730F7A2B5DCB6812877C4D138AF489BCEE441872FE54DCE6F8675F719F1922E32526DB"))
	expectedDecrypt := internal.Reverse(internal.MustDecodeHex("4657BB6ECEDE761D6780D2A83F3DAA0C780F50E938BDF37874366F4F9DBC5D8D315407127949A3C7CE2DF11AE1E369CEC828"))
	expectedEncrypt := internal.Reverse(internal.MustDecodeHex("462D681786E2CF7266F7AFF9786F30B17646F740FDEBDD9C592CFEDA3614D5533A7ABFCC03E33DC2AF8DB1252C8B3AFE1A41"))

	h := NewHeaderCrypto(sessionKey)
	h.Init()

	actualDecrypt := bytes.Clone(data)
	h.Decrypt(data)
	actualEncrypt := bytes.Clone(data)
	h.Encrypt(data)

	assert.Equal(t, expectedDecrypt, actualDecrypt)
	assert.Equal(t, expectedEncrypt, actualEncrypt)
}

func TestHeaderParse(t *testing.T) {
	t.Skip("TODO")

	// sessionKey := []byte{
	// 	0x2E, 0xFE, 0xE7, 0xB0, 0xC1, 0x77, 0xEB, 0xBD, 0xFF, 0x66, 0x76, 0xC5, 0x6E, 0xFC, 0x23,
	// 	0x39, 0xBE, 0x9C, 0xAD, 0x14, 0xBF, 0x8B, 0x54, 0xBB, 0x5A, 0x86, 0xFB, 0xF8, 0x1F, 0x6D,
	// 	0x42, 0x4A, 0xA2, 0x3C, 0xC9, 0xA3, 0x14, 0x9F, 0xB1, 0x75,
	// }
	// // seed := gom.DecodeHex("DEADBEEF")
	// data := gom.DecodeHex("60B1D4C5E50485EB")
	// // expectedSize := 19826
	// // expectedOpcode := 2589630381

	// h := NewWrathHeaderCrypto(sessionKey)
	// assert.NoError(t, h.Init())

	// // decrypting client header (6 bytes).
	// decrypted := h.Decrypt(data)
	// fmt.Printf("%x\n", decrypted)
	// size := binary.BigEndian.Uint16(decrypted[:2])
	// opcode := binary.LittleEndian.Uint32(decrypted[2:6])
	// fmt.Printf("size: %d\n", size)
	// fmt.Printf("opcode: %d\n", opcode)
}