package authd

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"log"
	"time"

	"github.com/kangaroux/gomaggus/internal/authd/packets"
	"github.com/kangaroux/gomaggus/internal/models"
	"github.com/kangaroux/gomaggus/internal/srp"
)

func handleLoginProof(services *Services, c *Client, data []byte) error {
	log.Println("Starting login proof")

	var serverProof []byte
	authenticated := false

	if c.account != nil {
		p := packets.ClientLoginProof{}
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
		resp := packets.ServerLoginProofFail{
			Opcode:    OP_LOGIN_PROOF,
			ErrorCode: WOW_FAIL_UNKNOWN_ACCOUNT,
		}
		binary.Write(&respBuf, binary.BigEndian, &resp)
	} else {
		resp := packets.ServerLoginProofSuccess{
			Opcode:           OP_LOGIN_PROOF,
			ErrorCode:        WOW_SUCCESS,
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
