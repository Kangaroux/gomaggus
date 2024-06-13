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
// FIELD ORDER MATTERS, DO NOT REORDER
type ClientReconnectProof struct {
	Opcode         Opcode // OpReconnectProof
	ProofData      [srp.ProofDataSize]byte
	ClientProof    [srp.ProofSize]byte
	ClientChecksum [20]byte
	KeyCount       uint8
}

func (p *ClientReconnectProof) Read(data []byte) error {
	reader := bytes.NewReader(data)
	return binary.Read(reader, binary.LittleEndian, p)
}

// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_proof_server.html#protocol-version-8
type ServerReconnectProof struct {
	Opcode    Opcode
	ErrorCode RespCode
	_         [2]byte // padding
}

func ReconnectProof(svc *authd.Service, c *authd.Client, data []byte) error {
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

			p := ClientReconnectProof{}
			if err := p.Read(data); err != nil {
				return err
			}

			serverProof := srp.CalculateReconnectProof(c.Username, p.ProofData[:], c.ReconnectData, c.SessionKey)
			authenticated = bytes.Equal(serverProof, p.ClientProof[:])
		}
	}

	resp := ServerReconnectProof{Opcode: OpReconnectProof}

	if !authenticated {
		resp.ErrorCode = CodeFailUnknownAccount
	} else {
		resp.ErrorCode = CodeSuccess
	}

	respBuf := bytes.Buffer{}
	binary.Write(&respBuf, binary.BigEndian, &resp)

	if _, err := c.Conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to reconnect proof")

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
