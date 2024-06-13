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

func handleReconnectProof(services *Services, c *Client, data []byte) error {
	log.Println("Starting reconnect proof")

	authenticated := false

	if c.account != nil {
		session, err := services.sessions.Get(c.account.Id)
		if err != nil {
			return err
		}

		// We can only try to reconnect the client if we have a previous session key
		if session != nil {
			if err := session.Decode(); err != nil {
				return err
			}
			c.sessionKey = session.SessionKey()

			p := packets.ClientReconnectProof{}
			if err := p.Read(data); err != nil {
				return err
			}

			serverProof := srp.CalculateReconnectProof(c.username, p.ProofData[:], c.reconnectData, c.sessionKey)
			authenticated = bytes.Equal(serverProof, p.ClientProof[:])
		}
	}

	resp := packets.ServerReconnectProof{Opcode: OpReconnectProof}

	if !authenticated {
		resp.ErrorCode = CodeFailUnknownAccount
	} else {
		resp.ErrorCode = CodeSuccess
	}

	respBuf := bytes.Buffer{}
	binary.Write(&respBuf, binary.BigEndian, &resp)

	if _, err := c.conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to reconnect proof")

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
