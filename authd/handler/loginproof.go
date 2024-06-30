package handler

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"time"

	srp "github.com/kangaroux/go-wow-srp6"
	"github.com/kangaroux/gomaggus/authd"
	"github.com/kangaroux/gomaggus/model"
	"github.com/mixcode/binarystruct"
)

type telemetryKey struct {
	Unknown    uint16
	Unknown2   uint32
	Unknown3   [4]byte
	CDKeyProof [20]byte
}

// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_h.Client.html#protocol-version-8
type loginProofRequest struct {
	Opcode           authd.Opcode // OpLoginProof
	ClientPublicKey  [srp.KeySize]byte
	ClientProof      [srp.ProofSize]byte
	CRCHash          [20]byte
	NumTelemetryKeys uint8
	_                []telemetryKey `binary:"[NumTelemetryKeys]Any"`
	SecurityFlag     uint8
}

// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_server.html#protocol-version-8
type loginProofFailed struct {
	Opcode    authd.Opcode // OpLoginProof
	ErrorCode authd.RespCode
	_         [2]byte // padding
}

type loginProofSuccess struct {
	Opcode           authd.Opcode // OpLoginProof
	ErrorCode        authd.RespCode
	Proof            [srp.ProofSize]byte
	AccountFlags     uint32
	HardwareSurveyId uint32
	_                [2]byte // padding
}

type LoginProof struct {
	Client   *authd.Client
	Sessions model.SessionService
	request  loginProofRequest
}

func (h *LoginProof) Handle() error {
	if h.Client.State != authd.StateAuthProof {
		return &ErrWrongState{
			Handler:  "LoginProof",
			Expected: authd.StateAuthProof,
			Actual:   h.Client.State,
		}
	}

	log.Println("Starting login proof")

	var serverProof []byte
	authenticated := false

	if h.Client.Account != nil {
		c := h.Client
		acct := h.Client.Account

		c.ClientPublicKey = h.request.ClientPublicKey[:]
		c.SessionKey = srp.SessionKey(c.ClientPublicKey, c.ServerPublicKey, c.PrivateKey, acct.Verifier())
		calculatedClientProof := srp.ClientChallengeProof(acct.Username, acct.Salt(), c.ClientPublicKey, c.ServerPublicKey, c.SessionKey)
		authenticated = bytes.Equal(calculatedClientProof, h.request.ClientProof[:])

		if authenticated {
			serverProof = srp.ServerChallengeProof(c.ClientPublicKey, h.request.ClientProof[:], c.SessionKey)
		}
	}

	respBuf := bytes.Buffer{}

	if !authenticated {
		resp := loginProofFailed{
			Opcode:    authd.OpcodeLoginProof,
			ErrorCode: authd.UnknownAccount,
		}
		binary.Write(&respBuf, binary.BigEndian, &resp)
	} else {
		resp := loginProofSuccess{
			Opcode:           authd.OpcodeLoginProof,
			ErrorCode:        authd.Success,
			AccountFlags:     0,
			HardwareSurveyId: 0,
		}
		copy(resp.Proof[:], serverProof)
		binary.Write(&respBuf, binary.BigEndian, &resp)
	}

	if _, err := h.Client.Conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to login proof")

	if authenticated {
		_, err := h.Sessions.UpdateOrCreate(&model.Session{
			AccountId:     h.Client.Account.Id,
			SessionKeyHex: hex.EncodeToString(h.Client.SessionKey),
			Connected:     1,
			ConnectedAt:   sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return err
		}

		h.Client.State = authd.StateAuthenticated
	} else {
		h.Client.State = authd.StateInvalid
	}

	return nil
}

// Read reads the packet data and parses it as a login proof request. If data is too small then
// Read returns ErrPacketReadEOF.
func (h *LoginProof) Read(data []byte) (int, error) {
	n, err := binarystruct.Unmarshal(data, binary.LittleEndian, &h.request)

	if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
		return 0, ErrPacketReadEOF
	} else if err != nil {
		return 0, err
	}

	return n, nil
}
