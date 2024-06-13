package srp

import (
	"strings"
	"testing"

	"github.com/kangaroux/gomaggus/internal"
	"github.com/stretchr/testify/assert"
)

func TestX(t *testing.T) {
	rows := internal.LoadTestData("../testdata/srp/calculate_x.csv")

	t.Run("generated test data", func(t *testing.T) {
		for _, row := range rows {
			username := row[0]
			password := row[1]
			salt := internal.DecodeHex(row[2])
			expected := internal.DecodeHex(row[3])

			assert.Equal(t, expected, CalculateX(username, password, salt))
		}
	})

	t.Run("username/pass are case insensitive", func(t *testing.T) {
		row := rows[0]
		username := row[0]
		password := row[1]
		salt := internal.DecodeHex(row[2])

		first := CalculateX(strings.ToLower(username), strings.ToLower(password), salt)
		second := CalculateX(strings.ToUpper(username), strings.ToUpper(password), salt)

		assert.Equal(t, first, second)
	})
}

func TestVerifier(t *testing.T) {
	rows := internal.LoadTestData("../testdata/srp/calculate_verifier.csv")

	t.Run("generated test data", func(t *testing.T) {
		for _, row := range rows {
			username := row[0]
			password := row[1]
			salt := internal.DecodeHex(row[2])
			expected := internal.DecodeHex(row[3])

			assert.Equal(t, expected, CalculateVerifier(username, password, salt))
		}
	})

	t.Run("username/pass are case insensitive", func(t *testing.T) {
		row := rows[0]
		username := row[0]
		password := row[1]
		salt := internal.DecodeHex(row[2])

		first := CalculateVerifier(strings.ToLower(username), strings.ToLower(password), salt)
		second := CalculateVerifier(strings.ToUpper(username), strings.ToUpper(password), salt)

		assert.Equal(t, first, second)
	})
}

func TestServerPublicKey(t *testing.T) {
	rows := internal.LoadTestData("../testdata/srp/calculate_server_public_key.csv")

	for _, row := range rows {
		verifier := internal.DecodeHex(row[0])
		privateKey := internal.DecodeHex(row[1])
		expected := internal.DecodeHex(row[2])

		assert.Equal(t, expected, CalculateServerPublicKey(verifier, privateKey))
	}
}

func TestCalculateU(t *testing.T) {
	rows := internal.LoadTestData("../testdata/srp/calculate_u.csv")

	for _, row := range rows {
		clientPublic := internal.DecodeHex(row[0])
		serverPublic := internal.DecodeHex(row[1])
		expected := internal.DecodeHex(row[2])

		assert.Equal(t, expected, CalculateU(clientPublic, serverPublic))
	}
}

func TestServerSKey(t *testing.T) {
	rows := internal.LoadTestData("../testdata/srp/calculate_server_s.csv")

	for _, row := range rows {
		clientPublic := internal.DecodeHex(row[0])
		verifier := internal.DecodeHex(row[1])
		u := internal.DecodeHex(row[2])
		serverPrivate := internal.DecodeHex(row[3])
		expected := internal.DecodeHex(row[4])

		assert.Equal(t, expected, CalculateServerSKey(clientPublic, verifier, u, serverPrivate))
	}
}

func TestInterleave(t *testing.T) {
	rows := internal.LoadTestData("../testdata/srp/calculate_interleaved.csv")

	for _, row := range rows {
		s := internal.DecodeHex(row[0])
		expected := internal.DecodeHex(row[1])

		assert.Equal(t, expected, CalculateInterleave(s))
	}
}

func TestServerSessionKey(t *testing.T) {
	rows := internal.LoadTestData("../testdata/srp/calculate_server_session_key.csv")

	for _, row := range rows {
		clientPublic := internal.DecodeHex(row[0])
		serverPrivate := internal.DecodeHex(row[1])
		verifier := internal.DecodeHex(row[2])
		expected := internal.DecodeHex(row[3])
		serverPublic := CalculateServerPublicKey(verifier, serverPrivate)

		assert.Equal(t, expected, CalculateServerSessionKey(clientPublic, serverPublic, serverPrivate, verifier))
	}
}
