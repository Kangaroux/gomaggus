package main

import (
	"crypto/rc4"
	"encoding/csv"
	"encoding/hex"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: make test util funcs
func decodeHex(s string) []byte {
	val, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return val
}

func loadTestData(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	return rows
}

func Reverse(data []byte) []byte {
	n := len(data)
	newData := make([]byte, n)
	for i := 0; i < n; i++ {
		newData[i] = data[n-i-1]
	}
	return newData
}

func TestGenerateKey(t *testing.T) {
	rows := loadTestData("test_data/header/wrath_generate_key.csv")

	for _, row := range rows {
		expected := decodeHex(row[2])
		sessionKey := decodeHex(row[0])
		fixedKey := decodeHex(row[1])
		wrath := &WrathHeaderCrypto{sessionKey: sessionKey}
		assert.Equal(t, expected, wrath.GenerateKey(fixedKey))
	}
}
func TestDrop1024(t *testing.T) {
	key := decodeHex("470446575F1EEDC1732473F4BDCACF4C13EE3837A3E2B9B720472E851525855DFC8045BF80D2F4D5")
	data := []byte("hello world")
	result := make([]byte, len(data))

	// Test the output when no bytes are dropped
	c, _ := rc4.NewCipher(key)
	copy(result, data)
	c.XORKeyStream(result, data)
	assert.Equal(t, decodeHex("B53577EC1E7FEAE6CDBB5C"), result)

	// Test the output when the first 1024 bytes are dropped (RC4-drop1024)
	c, _ = rc4.NewCipher(key)
	drop1024(c)
	copy(result, data)
	c.XORKeyStream(result, data)
	assert.Equal(t, decodeHex("9F273843E006E43B0D33B6"), result)
}

func TestEncryptDecrypt(t *testing.T) {
	sessionKey := Reverse(decodeHex("403FCE7B2B1FCDAE43F118B6C7E517D5A1498088180936A3E45B9888978B7675ECBAA7DB4CA4E8DE"))
	data := Reverse(decodeHex("DEBA4C8DCBD613F06E725123E887CF730F7A2B5DCB6812877C4D138AF489BCEE441872FE54DCE6F8675F719F1922E32526DB"))
	expectedDecrypt := Reverse(decodeHex("4657BB6ECEDE761D6780D2A83F3DAA0C780F50E938BDF37874366F4F9DBC5D8D315407127949A3C7CE2DF11AE1E369CEC828"))
	expectedEncrypt := Reverse(decodeHex("462D681786E2CF7266F7AFF9786F30B17646F740FDEBDD9C592CFEDA3614D5533A7ABFCC03E33DC2AF8DB1252C8B3AFE1A41"))

	h := NewWrathHeaderCrypto(sessionKey)
	h.Init()
	assert.Equal(t, expectedDecrypt, h.Decrypt(data))
	assert.Equal(t, expectedEncrypt, h.Encrypt(data))
}
