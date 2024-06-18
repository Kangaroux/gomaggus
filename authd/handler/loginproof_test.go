package handler

import (
	"bytes"
	"io"
	"testing"

	"github.com/kangaroux/gomaggus/authd"
	"github.com/kangaroux/gomaggus/authd/mock"
	"github.com/kangaroux/gomaggus/internal"
	"github.com/kangaroux/gomaggus/model"
	"github.com/mixcode/binarystruct"
	"github.com/stretchr/testify/assert"
)

func TestLoginProof(t *testing.T) {
	var client *authd.Client
	var conn *mock.Conn
	var sessions *mock.SessionService

	newHandler := func() *LoginProof {
		conn = &mock.Conn{}
		client = &authd.Client{
			Conn:  conn,
			State: authd.StateAuthProof,
		}
		sessions = &mock.SessionService{}
		return &LoginProof{
			Client:   client,
			Sessions: sessions,
		}
	}

	t.Run("invalid state", func(t *testing.T) {
		h := newHandler()
		client.State = authd.StateAuthChallenge
		err, ok := h.Handle([]byte{}).(*ErrWrongState)

		assert.Error(t, err)
		assert.True(t, ok)
	})

	t.Run("malformed packet", func(t *testing.T) {
		assert.Equal(t, io.EOF, newHandler().Handle([]byte{}))
	})

	t.Run("from fake challenge response", func(t *testing.T) {
		h := newHandler()
		packet := &loginProofFailed{
			Opcode:    authd.OpcodeLoginProof,
			ErrorCode: authd.UnknownAccount,
		}
		expectedResp := internal.MustMarshal(packet, binarystruct.LittleEndian)
		request := internal.MustMarshal(loginProofRequest{}, binarystruct.LittleEndian)

		conn.OnWrite = func(actual []byte) (int, error) {
			// Server sent the expected bytes
			assert.True(t, bytes.Equal(expectedResp, actual))
			return 0, nil
		}
		// A nil account means the challenge response was faked
		client.Account = nil

		assert.NoError(t, h.Handle(request))
		assert.Equal(t, authd.StateInvalid, client.State) // invalid state
	})

	t.Run("success", func(t *testing.T) {
		h := newHandler()

		requestPacket := &loginProofRequest{
			// Leaving everything blank/zero will result in this proof
			ClientProof: [20]byte(internal.MustDecodeHex("9E224007DEE3D15873D71FCF7D8CD8D94C53DCAA")),
		}
		request := internal.MustMarshal(requestPacket, binarystruct.LittleEndian)

		respPacket := &loginProofSuccess{
			Opcode:           authd.OpcodeLoginProof,
			ErrorCode:        authd.Success,
			Proof:            [20]byte(internal.MustDecodeHex("979F4506AF22E2A3C3BA8C122350BB2B9D144CE2")),
			AccountFlags:     0,
			HardwareSurveyId: 0,
		}
		expectedResp := internal.MustMarshal(respPacket, binarystruct.LittleEndian)

		conn.OnWrite = func(actual []byte) (int, error) {
			// Server sent the expected bytes
			assert.True(t, bytes.Equal(expectedResp, actual))
			return 0, nil
		}
		client.Account = &model.Account{}
		client.Account.DecodeSrp()

		assert.NoError(t, h.Handle(request))
		assert.Equal(t, authd.StateAuthenticated, client.State) // authenticated
	})
}
