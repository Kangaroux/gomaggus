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
)

// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_client.html#protocol-version-8
type loginProofRequest struct {
	Opcode           authd.Opcode // OpLoginProof
	ClientPublicKey  [srp.KeySize]byte
	ClientProof      [srp.ProofSize]byte
	CRCHash          [20]byte
	NumTelemetryKeys uint8
}

func (p *loginProofRequest) Read(data []byte) error {
	reader := bytes.NewReader(data)
	return binary.Read(reader, binary.LittleEndian, p)
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

func LoginProof(svc *authd.Service, c *authd.Client, data []byte) error {
	if c.State != authd.StateAuthProof {
		return &ErrWrongState{
			Handler:  "LoginProof",
			Expected: authd.StateAuthProof,
			Actual:   c.State,
		}
	}

	log.Println("Starting login proof")

	var serverProof []byte
	authenticated := false

	if c.Account != nil {
		p := loginProofRequest{}
		if err := p.Read(data); err != nil {
			return err
		}

		c.ClientPublicKey = p.ClientPublicKey[:]
		c.SessionKey = srp.CalculateServerSessionKey(
			c.ClientPublicKey,
			c.ServerPublicKey,
			c.PrivateKey,
			c.Account.Verifier(),
		)
		calculatedClientProof := srp.CalculateClientProof(
			c.Account.Username,
			c.Account.Salt(),
			c.ClientPublicKey,
			c.ServerPublicKey,
			c.SessionKey,
		)
		authenticated = bytes.Equal(calculatedClientProof, p.ClientProof[:])

		if authenticated {
			serverProof = srp.CalculateServerProof(c.ClientPublicKey, p.ClientProof[:], c.SessionKey)
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

	if _, err := c.Conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to login proof")

	if authenticated {
		err := svc.Sessions.UpdateOrCreate(&model.Session{
			AccountId:     c.Account.Id,
			SessionKeyHex: hex.EncodeToString(c.SessionKey),
			Connected:     1,
			ConnectedAt:   sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return err
		}

		c.State = authd.StateAuthenticated
	} else {
		c.State = authd.StateInvalid
	}

	return nil
}
