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

// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_client.html#protocol-version-8
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

func LoginProof(svc *authd.Service, client *authd.Client, data []byte) error {
	if client.State != authd.StateAuthProof {
		return &ErrWrongState{
			Handler:  "LoginProof",
			Expected: authd.StateAuthProof,
			Actual:   client.State,
		}
	}

	log.Println("Starting login proof")

	var serverProof []byte
	authenticated := false

	if client.Account != nil {
		req := loginProofRequest{}
		if _, err := binarystruct.Unmarshal(data, binarystruct.LittleEndian, &req); err != nil {
			return err
		}

		client.ClientPublicKey = req.ClientPublicKey[:]
		client.SessionKey = srp.CalculateServerSessionKey(client.ClientPublicKey, client.ServerPublicKey, client.PrivateKey, client.Account.Verifier())

		calculatedClientProof := srp.CalculateClientProof(client.Account.Username, client.Account.Salt(), client.ClientPublicKey, client.ServerPublicKey, client.SessionKey)
		authenticated = bytes.Equal(calculatedClientProof, req.ClientProof[:])

		if authenticated {
			serverProof = srp.CalculateServerProof(client.ClientPublicKey, req.ClientProof[:], client.SessionKey)
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

	if _, err := client.Conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to login proof")

	if authenticated {
		err := svc.Sessions.UpdateOrCreate(&model.Session{
			AccountId:     client.Account.Id,
			SessionKeyHex: hex.EncodeToString(client.SessionKey),
			Connected:     1,
			ConnectedAt:   sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return err
		}

		client.State = authd.StateAuthenticated
	} else {
		client.State = authd.StateInvalid
	}

	return nil
}
