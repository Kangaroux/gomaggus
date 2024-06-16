package handler

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"log"
	"time"

	"github.com/kangaroux/gomaggus/authd"
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/srp"
	"github.com/mixcode/binarystruct"
)

// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_h.Client.html#protocol-version-8
type loginProofRequest struct {
	Opcode           authd.Opcode // OpLoginProof
	ClientPublicKey  [srp.KeySize]byte
	ClientProof      [srp.ProofSize]byte
	CRCHash          [20]byte
	NumTelemetryKeys uint8
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
}

func (h *LoginProof) Handle(data []byte) error {
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
		req := loginProofRequest{}
		if _, err := binarystruct.Unmarshal(data, binarystruct.LittleEndian, &req); err != nil {
			return err
		}

		c := h.Client
		acct := h.Client.Account

		c.ClientPublicKey = req.ClientPublicKey[:]
		c.SessionKey = srp.CalculateServerSessionKey(c.ClientPublicKey, c.ServerPublicKey, c.PrivateKey, acct.Verifier())
		calculatedClientProof := srp.CalculateClientProof(acct.Username, acct.Salt(), c.ClientPublicKey, c.ServerPublicKey, c.SessionKey)
		authenticated = bytes.Equal(calculatedClientProof, req.ClientProof[:])

		if authenticated {
			serverProof = srp.CalculateServerProof(c.ClientPublicKey, req.ClientProof[:], c.SessionKey)
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
		err := h.Sessions.UpdateOrCreate(&model.Session{
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
