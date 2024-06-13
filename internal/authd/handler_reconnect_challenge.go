package authd

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"log"

	"github.com/kangaroux/gomaggus/internal/models"
)

func handleReconnectChallenge(services *Services, c *Client, data []byte) error {
	log.Println("Starting reconnect challenge")

	var err error
	p := ClientLoginChallenge{}
	if err = p.Read(data); err != nil {
		return err
	}
	c.username = p.Username

	log.Printf("client trying to reconnect as '%s'\n", c.username)

	c.account, err = services.accounts.Get(&models.AccountGetParams{Username: c.username})
	if err != nil {
		return err
	}

	// Generate random data that will be used for the reconnect proof
	if _, err := rand.Read(c.reconnectData); err != nil {
		return err
	}

	resp := ServerReconnectChallenge{
		Opcode: OpReconnectChallenge,

		// Always return success to prevent a bad actor from mining usernames.
		ErrorCode:    CodeSuccess,
		ChecksumSalt: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	copy(resp.ReconnectData[:], c.reconnectData)

	respBuf := bytes.Buffer{}
	binary.Write(&respBuf, binary.BigEndian, &resp)

	if _, err := c.conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to reconnect challenge")

	c.state = StateReconnectProof

	return nil
}
