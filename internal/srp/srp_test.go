package srp

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/kangaroux/gomaggus/internal"
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
	rows := internal.LoadTestData("../test_data/srp/calculate_x.csv")

	t.Run("generated test data", func(t *testing.T) {
		for _, row := range rows {
			username := row[0]
			password := row[1]
			salt := decodeHex(row[2])
			expected := decodeHex(row[3])

			assert.Equal(t, expected, CalculateX(username, password, salt))
		}
	})

	t.Run("username/pass are case insensitive", func(t *testing.T) {
		row := rows[0]
		username := row[0]
		password := row[1]
		salt := decodeHex(row[2])

		first := CalculateX(strings.ToLower(username), strings.ToLower(password), salt)
		second := CalculateX(strings.ToUpper(username), strings.ToUpper(password), salt)

		assert.Equal(t, first, second)
	})
}

func TestVerifier(t *testing.T) {
	rows := internal.LoadTestData("../test_data/srp/calculate_verifier.csv")

	t.Run("generated test data", func(t *testing.T) {
		for _, row := range rows {
			username := row[0]
			password := row[1]
			salt := decodeHex(row[2])
			expected := decodeHex(row[3])

			assert.Equal(t, expected, CalculateVerifier(username, password, salt))
		}
	})

	t.Run("username/pass are case insensitive", func(t *testing.T) {
		row := rows[0]
		username := row[0]
		password := row[1]
		salt := decodeHex(row[2])

		first := CalculateVerifier(strings.ToLower(username), strings.ToLower(password), salt)
		second := CalculateVerifier(strings.ToUpper(username), strings.ToUpper(password), salt)

		assert.Equal(t, first, second)
	})
}

func TestServerPublicKey(t *testing.T) {
	rows := internal.LoadTestData("../test_data/srp/calculate_server_public_key.csv")

	for _, row := range rows {
		verifier := decodeHex(row[0])
		privateKey := decodeHex(row[1])
		expected := decodeHex(row[2])

		assert.Equal(t, expected, CalculateServerPublicKey(verifier, privateKey))
	}
}

func TestCalculateU(t *testing.T) {
	rows := internal.LoadTestData("../test_data/srp/calculate_u.csv")

	for _, row := range rows {
		clientPublic := decodeHex(row[0])
		serverPublic := decodeHex(row[1])
		expected := decodeHex(row[2])

		assert.Equal(t, expected, CalculateU(clientPublic, serverPublic))
	}
}

func TestServerSKey(t *testing.T) {
	rows := internal.LoadTestData("../test_data/srp/calculate_server_s.csv")

	for _, row := range rows {
		clientPublic := decodeHex(row[0])
		verifier := decodeHex(row[1])
		u := decodeHex(row[2])
		serverPrivate := decodeHex(row[3])
		expected := decodeHex(row[4])

		assert.Equal(t, expected, CalculateServerSKey(clientPublic, verifier, u, serverPrivate))
	}
}

func TestInterleave(t *testing.T) {
	rows := internal.LoadTestData("../test_data/srp/calculate_interleaved.csv")

	for _, row := range rows {
		s := decodeHex(row[0])
		expected := decodeHex(row[1])

		assert.Equal(t, expected, CalculateInterleave(s))
	}
}

func TestServerSessionKey(t *testing.T) {
	rows := internal.LoadTestData("../test_data/srp/calculate_server_session_key.csv")

	for _, row := range rows {
		clientPublic := decodeHex(row[0])
		serverPrivate := decodeHex(row[1])
		verifier := decodeHex(row[2])
		expected := decodeHex(row[3])
		serverPublic := CalculateServerPublicKey(verifier, serverPrivate)

		assert.Equal(t, expected, CalculateServerSessionKey(clientPublic, serverPublic, serverPrivate, verifier))
	}
}

func TestClientProof(t *testing.T) {
	rows := internal.LoadTestData("../test_data/srp/calculate_client_proof.csv")

	t.Run("generated test data", func(t *testing.T) {
		for _, row := range rows {
			username := row[0]
			salt := decodeHex(row[1])
			clientPublic := decodeHex(row[2])
			serverPublic := decodeHex(row[3])
			sessionKey := decodeHex(row[4])
			expected := decodeHex(row[5])

			assert.Equal(t, expected, CalculateClientProof(username, salt, clientPublic, serverPublic, sessionKey))
		}
	})

	t.Run("username is case insensitive", func(t *testing.T) {
		row := rows[0]
		username := row[0]
		salt := decodeHex(row[1])
		clientPublic := decodeHex(row[2])
		serverPublic := decodeHex(row[3])
		sessionKey := decodeHex(row[4])

		first := CalculateClientProof(strings.ToLower(username), salt, clientPublic, serverPublic, sessionKey)
		second := CalculateClientProof(strings.ToUpper(username), salt, clientPublic, serverPublic, sessionKey)

		assert.Equal(t, first, second)
	})
}

func TestServerProof(t *testing.T) {
	rows := internal.LoadTestData("../test_data/srp/calculate_server_proof.csv")

	for _, row := range rows {
		clientPublic := decodeHex(row[0])
		clientProof := decodeHex(row[1])
		sessionKey := decodeHex(row[2])
		expected := decodeHex(row[3])

		assert.Equal(t, expected, CalculateServerProof(clientPublic, clientProof, sessionKey))
	}
}

func TestReconnectProof(t *testing.T) {
	rows := internal.LoadTestData("../test_data/srp/calculate_reconnect_proof.csv")

	t.Run("generated test data", func(t *testing.T) {
		for _, row := range rows {
			username := row[0]
			clientData := decodeHex(row[1])
			serverData := decodeHex(row[2])
			sessionKey := decodeHex(row[3])
			expected := decodeHex(row[4])

			assert.Equal(t, expected, CalculateReconnectProof(username, clientData, serverData, sessionKey))
		}
	})

	t.Run("username is case insensitive", func(t *testing.T) {
		row := rows[0]
		username := row[0]
		clientData := decodeHex(row[1])
		serverData := decodeHex(row[2])
		sessionKey := decodeHex(row[3])

		first := CalculateReconnectProof(strings.ToLower(username), clientData, serverData, sessionKey)
		second := CalculateReconnectProof(strings.ToUpper(username), clientData, serverData, sessionKey)

		assert.Equal(t, first, second)
	})
}
