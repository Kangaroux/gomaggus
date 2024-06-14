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

// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_proof_client.html
type reconnectProofRequest struct {
	Opcode         authd.Opcode // OpReconnectProof
	ProofData      [srp.ProofDataSize]byte
	ClientProof    [srp.ProofSize]byte
	ClientChecksum [20]byte
	KeyCount       uint8
}

func (p *reconnectProofRequest) Read(data []byte) error {
	reader := bytes.NewReader(data)
	return binary.Read(reader, binary.LittleEndian, p)
}

// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_proof_server.html#protocol-version-8
type reconnectProofResponse struct {
	Opcode    authd.Opcode
	ErrorCode authd.RespCode
	_         [2]byte // padding
}

func ReconnectProof(svc *authd.Service, c *authd.Client, data []byte) error {
	if c.State != authd.StateAuthProof {
		return &ErrWrongState{
			Handler:  "RealmList",
			Expected: authd.StateAuthProof,
			Actual:   c.State,
		}
	}

	log.Println("Starting reconnect proof")

	authenticated := false

	if c.Account != nil {
		session, err := svc.Sessions.Get(c.Account.Id)
		if err != nil {
			return err
		}

		// We can only try to reconnect the client if we have a previous session key
		if session != nil {
			if err := session.Decode(); err != nil {
				return err
			}
			c.SessionKey = session.SessionKey()

			p := reconnectProofRequest{}
			if err := p.Read(data); err != nil {
				return err
			}

			serverProof := srp.CalculateReconnectProof(c.Username, p.ProofData[:], c.ReconnectData, c.SessionKey)
			authenticated = bytes.Equal(serverProof, p.ClientProof[:])
		}
	}

	resp := reconnectProofResponse{Opcode: authd.OpcodeReconnectProof}

	if !authenticated {
		resp.ErrorCode = authd.UnknownAccount
	} else {
		resp.ErrorCode = authd.Success
	}

	respBuf := bytes.Buffer{}
	binary.Write(&respBuf, binary.BigEndian, &resp)

	if _, err := c.Conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to reconnect proof")

	if authenticated {
		session := model.Session{
			AccountId:     c.Account.Id,
			SessionKeyHex: hex.EncodeToString(c.SessionKey),
			Connected:     1,
			ConnectedAt:   sql.NullTime{Time: time.Now(), Valid: true},
		}
		if err := svc.Sessions.UpdateOrCreate(&session); err != nil {
			return err
		}
		c.State = authd.StateAuthenticated
	} else {
		c.State = authd.StateInvalid
	}

	return nil
}
