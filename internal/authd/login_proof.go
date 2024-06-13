package authd

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"log"
	"time"

	"github.com/kangaroux/gomaggus/internal/models"
	"github.com/kangaroux/gomaggus/internal/srp"
)

// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_client.html#protocol-version-8
// FIELD ORDER MATTERS, DO NOT REORDER
type ClientLoginProof struct {
	Opcode           Opcode // OpLoginProof
	ClientPublicKey  [srp.KeySize]byte
	ClientProof      [srp.ProofSize]byte
	CRCHash          [20]byte
	NumTelemetryKeys uint8
}

func (p *ClientLoginProof) Read(data []byte) error {
	reader := bytes.NewReader(data)
	return binary.Read(reader, binary.LittleEndian, p)
}

// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_server.html#protocol-version-8
type ServerLoginProofFail struct {
	Opcode    Opcode // OpLoginProof
	ErrorCode ErrorCode
	_         [2]byte // padding
}

type ServerLoginProofSuccess struct {
	Opcode           Opcode // OpLoginProof
	ErrorCode        ErrorCode
	Proof            [srp.ProofSize]byte
	AccountFlags     uint32
	HardwareSurveyId uint32
	_                [2]byte // padding
}

func handleLoginProof(services *Services, c *Client, data []byte) error {
	log.Println("Starting login proof")

	var serverProof []byte
	authenticated := false

	if c.account != nil {
		p := ClientLoginProof{}
		if err := p.Read(data); err != nil {
			return err
		}

		c.clientPublicKey = p.ClientPublicKey[:]
		c.sessionKey = srp.CalculateServerSessionKey(
			c.clientPublicKey,
			c.serverPublicKey,
			c.privateKey,
			c.account.Verifier(),
		)
		calculatedClientProof := srp.CalculateClientProof(
			c.account.Username,
			c.account.Salt(),
			c.clientPublicKey,
			c.serverPublicKey,
			c.sessionKey,
		)
		authenticated = bytes.Equal(calculatedClientProof, p.ClientProof[:])

		if authenticated {
			serverProof = srp.CalculateServerProof(c.clientPublicKey, p.ClientProof[:], c.sessionKey)
		}
	}

	respBuf := bytes.Buffer{}

	if !authenticated {
		resp := ServerLoginProofFail{
			Opcode:    OpLoginProof,
			ErrorCode: CodeFailUnknownAccount,
		}
		binary.Write(&respBuf, binary.BigEndian, &resp)
	} else {
		resp := ServerLoginProofSuccess{
			Opcode:           OpLoginProof,
			ErrorCode:        CodeSuccess,
			AccountFlags:     0,
			HardwareSurveyId: 0,
		}
		copy(resp.Proof[:], serverProof)
		binary.Write(&respBuf, binary.BigEndian, &resp)
	}

	if _, err := c.conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to login proof")

	if authenticated {
		err := services.sessions.UpdateOrCreate(&models.Session{
			AccountId:     c.account.Id,
			SessionKeyHex: hex.EncodeToString(c.sessionKey),
			Connected:     1,
			ConnectedAt:   sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return err
		}

		c.state = StateAuthenticated
	} else {
		c.state = StateInvalid
	}

	return nil
}
