package srp

import (
	"encoding/binary"
	"strings"
	"testing"

	"github.com/kangaroux/gomaggus/internal"
	"github.com/stretchr/testify/assert"
)

func TestClientProof(t *testing.T) {
	rows := internal.LoadTestData("../test_data/srp/calculate_client_proof.csv")

	t.Run("generated test data", func(t *testing.T) {
		for _, row := range rows {
			username := row[0]
			salt := internal.DecodeHex(row[1])
			clientPublic := internal.DecodeHex(row[2])
			serverPublic := internal.DecodeHex(row[3])
			sessionKey := internal.DecodeHex(row[4])
			expected := internal.DecodeHex(row[5])

			assert.Equal(t, expected, CalculateClientProof(username, salt, clientPublic, serverPublic, sessionKey))
		}
	})

	t.Run("username is case insensitive", func(t *testing.T) {
		row := rows[0]
		username := row[0]
		salt := internal.DecodeHex(row[1])
		clientPublic := internal.DecodeHex(row[2])
		serverPublic := internal.DecodeHex(row[3])
		sessionKey := internal.DecodeHex(row[4])

		first := CalculateClientProof(strings.ToLower(username), salt, clientPublic, serverPublic, sessionKey)
		second := CalculateClientProof(strings.ToUpper(username), salt, clientPublic, serverPublic, sessionKey)

		assert.Equal(t, first, second)
	})
}

func TestServerProof(t *testing.T) {
	rows := internal.LoadTestData("../test_data/srp/calculate_server_proof.csv")

	for _, row := range rows {
		clientPublic := internal.DecodeHex(row[0])
		clientProof := internal.DecodeHex(row[1])
		sessionKey := internal.DecodeHex(row[2])
		expected := internal.DecodeHex(row[3])

		assert.Equal(t, expected, CalculateServerProof(clientPublic, clientProof, sessionKey))
	}
}

func TestReconnectProof(t *testing.T) {
	rows := internal.LoadTestData("../test_data/srp/calculate_reconnect_proof.csv")

	t.Run("generated test data", func(t *testing.T) {
		for _, row := range rows {
			username := row[0]
			clientData := internal.DecodeHex(row[1])
			serverData := internal.DecodeHex(row[2])
			sessionKey := internal.DecodeHex(row[3])
			expected := internal.DecodeHex(row[4])

			assert.Equal(t, expected, CalculateReconnectProof(username, clientData, serverData, sessionKey))
		}
	})

	t.Run("username is case insensitive", func(t *testing.T) {
		row := rows[0]
		username := row[0]
		clientData := internal.DecodeHex(row[1])
		serverData := internal.DecodeHex(row[2])
		sessionKey := internal.DecodeHex(row[3])

		first := CalculateReconnectProof(strings.ToLower(username), clientData, serverData, sessionKey)
		second := CalculateReconnectProof(strings.ToUpper(username), clientData, serverData, sessionKey)

		assert.Equal(t, first, second)
	})
}

func TestCalculateWorldProof(t *testing.T) {
	t.Skip("FIXME")

	expected := internal.DecodeHex("6095EB678CD195253F66F32BADA785CA6D9376B2")
	username := "TNDQWSHEBWHPABV2"
	clientSeed := make([]byte, 4)
	serverSeed := make([]byte, 4)
	binary.BigEndian.PutUint32(clientSeed, 1454143186)
	binary.BigEndian.PutUint32(serverSeed, 309086257)
	sessionKey := internal.DecodeHex("914D6219A99109D6BD946F6E6AF12BB611C59A22531C6F1A3F3CF58624D528DC163BE43813112C3D")

	assert.Equal(t, expected, CalculateWorldProof(username, clientSeed, serverSeed, sessionKey))
}
