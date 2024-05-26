package main

import (
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustDecodeHex(s string) []byte {
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func hexToByteArray(s string, bigEndian bool) *ByteArray {
	return NewByteArray(mustDecodeHex(s), 0, bigEndian)
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

func Test_constants(t *testing.T) {
	assert.Equal(t, bigK(), big.NewInt(3))
	assert.Equal(t, bigG(), big.NewInt(7))
	assert.Equal(t, bigN().Bytes(), []byte{
		0x89, 0x4B, 0x64, 0x5E, 0x89, 0xE1, 0x53, 0x5B,
		0xBD, 0xAD, 0x5B, 0x8B, 0x29, 0x06, 0x50, 0x53,
		0x08, 0x01, 0xB1, 0x8E, 0xBF, 0xBF, 0x5E, 0x8F,
		0xAB, 0x3C, 0x82, 0x87, 0x2A, 0x3E, 0x9B, 0xB7,
	})

	hBigN := sha1.Sum(largeSafePrime.LittleEndian().Bytes())
	hBigG := sha1.Sum(bigG().Bytes())
	computedXorHash := make([]byte, 20)

	for i := 0; i < 20; i++ {
		computedXorHash[i] = hBigN[i] ^ hBigG[i]
	}

	assert.Equal(t, xorHash.Bytes(), computedXorHash)
}

func Test_calcX(t *testing.T) {
	t.Run("generated test data", func(t *testing.T) {
		rows := loadTestData("test_data/srp/calculate_x.csv")

		for _, row := range rows {
			username := row[0]
			password := row[1]
			salt := hexToByteArray(row[2], false)
			expected := hexToByteArray(row[3], false).BigInt()

			assert.Equal(t, expected, calcX(username, password, salt))
		}
	})

	t.Run("case insensitive username and password", func(t *testing.T) {
		salt := hexToByteArray("F3DCABA1165E23534CDC7D709E87B409C505C28D26D3DC14247796BE29CC4D24", false)

		assert.Equal(t, calcX("username", "password", salt), calcX("USERNAME", "PASSWORD", salt))
	})
}

func Test_passVerify(t *testing.T) {
	t.Run("generated test data", func(t *testing.T) {
		rows := loadTestData("test_data/srp/calculate_verifier.csv")

		for _, row := range rows {
			username := row[0]
			password := row[1]
			salt := hexToByteArray(row[2], false)
			expected := hexToByteArray(row[3], false).BigInt()

			assert.Equal(t, expected, passVerify(username, password, salt))
		}
	})

	t.Run("case insensitive username and password", func(t *testing.T) {
		salt := hexToByteArray("F3DCABA1165E23534CDC7D709E87B409C505C28D26D3DC14247796BE29CC4D24", false)

		assert.Equal(t, passVerify("username", "password", salt), passVerify("USERNAME", "PASSWORD", salt))
	})
}

func Test_calcServerPublicKey(t *testing.T) {
	rows := loadTestData("test_data/srp/calculate_server_public_key.csv")

	for _, row := range rows {
		verifier := hexToByteArray(row[0], false).BigInt()
		serverPrivateKey := hexToByteArray(row[1], false).BigInt()
		expected := hexToByteArray(row[2], false).BigInt()

		assert.Equal(t, expected, calcServerPublicKey(verifier, serverPrivateKey))
	}
}

func Test_calcClientSKey(t *testing.T) {
	rows := loadTestData("test_data/srp/calculate_client_s.csv")

	for _, row := range rows {
		serverPublicKey := hexToByteArray(row[0], false).BigInt()
		clientPrivateKey := hexToByteArray(row[1], false).BigInt()
		x := hexToByteArray(row[2], false).BigInt()
		u := hexToByteArray(row[3], false).BigInt()
		expected := hexToByteArray(row[4], false)

		assert.Equal(t, expected, calcClientSKey(clientPrivateKey, serverPublicKey, x, u))
	}
}

func Test_calcServerSKey(t *testing.T) {
	rows := loadTestData("test_data/srp/calculate_server_s.csv")

	for _, row := range rows {
		clientPublicKey := hexToByteArray(row[0], false).BigInt()
		serverPrivateKey := hexToByteArray(row[1], false).BigInt()
		verifier := hexToByteArray(row[2], false).BigInt()
		u := hexToByteArray(row[3], false).BigInt()
		expected := hexToByteArray(row[4], false)

		assert.Equal(t, expected, calcServerSKey(clientPublicKey, verifier, u, serverPrivateKey))
	}
}

func Test_calcU(t *testing.T) {
	rows := loadTestData("test_data/srp/calculate_u.csv")

	for _, row := range rows {
		clientPublicKey := hexToByteArray(row[0], false).BigInt()
		serverPublicKey := hexToByteArray(row[1], false).BigInt()
		expected := hexToByteArray(row[2], false).BigInt()

		assert.Equal(t, expected, calcU(clientPublicKey, serverPublicKey))
	}
}

func Test_prepareInterleave(t *testing.T) {
	type testCase struct {
		S        []byte
		expected []byte
	}

	testCases := []testCase{
		{[]byte{0, 1}, []byte{}},
		{[]byte{0, 1, 2, 3}, []byte{2, 3}},
		{[]byte{1, 0, 0, 0}, []byte{1, 0, 0, 0}},
		{[]byte{0, 0, 0, 0}, []byte{}},
	}

	for _, tc := range testCases {
		S := NewByteArray(tc.S, 0, false)
		assert.Equal(t, tc.expected, prepareInterleave(S))
	}
}

func Test_calcInterleave(t *testing.T) {
	rows := loadTestData("test_data/srp/calculate_interleaved.csv")

	for _, row := range rows {
		S := hexToByteArray(row[0], false)
		expected := hexToByteArray(row[1], false)

		assert.Equal(t, expected, calcInterleave(S))
	}
}

func Test_calcServerSessionKey(t *testing.T) {
	rows := loadTestData("test_data/srp/calculate_server_session_key.csv")

	for _, row := range rows {
		clientPublicKey := hexToByteArray(row[0], false).BigInt()
		verifier := hexToByteArray(row[1], false).BigInt()
		serverPrivateKey := hexToByteArray(row[2], false).BigInt()
		expected := hexToByteArray(row[3], false)
		serverPublicKey := calcServerPublicKey(verifier, serverPrivateKey)

		assert.Equal(
			t,
			expected,
			calcServerSessionKey(clientPublicKey,
				serverPublicKey,
				verifier,
				serverPrivateKey,
			),
		)
	}
}

func Test_calcClientSessionKey(t *testing.T) {
	rows := loadTestData("test_data/srp/calculate_client_session_key.csv")

	for _, row := range rows {
		username := row[0]
		password := row[1]
		clientPublicKey := hexToByteArray(row[2], false).BigInt()
		clientPrivateKey := hexToByteArray(row[3], false).BigInt()
		serverPublicKey := hexToByteArray(row[4], false).BigInt()
		salt := hexToByteArray(row[5], false)
		expected := hexToByteArray(row[6], false)

		assert.Equal(
			t,
			expected,
			calcClientSessionKey(username,
				password,
				serverPublicKey,
				clientPrivateKey,
				clientPublicKey,
				salt,
			),
		)
	}
}

func Test_calcServerProof(t *testing.T) {
	rows := loadTestData("test_data/srp/calculate_server_proof.csv")

	for _, row := range rows {
		clientPublicKey := hexToByteArray(row[0], false).BigInt()
		clientProof := hexToByteArray(row[1], false)
		sessionKey := hexToByteArray(row[2], false)
		expected := hexToByteArray(row[3], false).LittleEndian()

		assert.Equal(
			t,
			expected,
			calcServerProof(clientPublicKey,
				clientProof,
				sessionKey,
			),
		)
	}
}

func Test_calcClientProof(t *testing.T) {
	rows := loadTestData("test_data/srp/calculate_client_proof.csv")

	for _, row := range rows {
		username := row[0]
		sessionKey := hexToByteArray(row[1], false)
		clientPublicKey := hexToByteArray(row[2], false).BigInt()
		serverPublicKey := hexToByteArray(row[3], false).BigInt()
		salt := hexToByteArray(row[4], false)
		expected := hexToByteArray(row[5], false)

		assert.Equal(
			t,
			expected,
			calcClientProof(username, sessionKey, clientPublicKey, serverPublicKey, salt),
		)
	}
}
