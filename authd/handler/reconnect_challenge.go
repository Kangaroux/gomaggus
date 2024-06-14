package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"log"

	"github.com/kangaroux/gomaggus/authd"
	"github.com/kangaroux/gomaggus/model"
)

const (
	ReconnectDataLen = 16
)

type reconnectChallengeRequest = loginChallengeRequest

// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_challenge_server.html#protocol-version-8
type reconnectChallengeResponse struct {
	Opcode        authd.Opcode // OpReconnectChallenge
	ErrorCode     authd.RespCode
	ReconnectData [ReconnectDataLen]byte
	ChecksumSalt  [16]byte
}

func ReconnectChallenge(svc *authd.Service, c *authd.Client, data []byte) error {
	if c.State != authd.StateAuthChallenge {
		return &ErrWrongState{
			Handler:  "RealmList",
			Expected: authd.StateAuthChallenge,
			Actual:   c.State,
		}
	}

	log.Println("Starting reconnect challenge")

	var err error
	p := reconnectChallengeRequest{}
	if err = p.Read(data); err != nil {
		return err
	}
	c.Username = p.Username

	log.Printf("client trying to reconnect as '%s'\n", c.Username)

	c.Account, err = svc.Accounts.Get(&model.AccountGetParams{Username: c.Username})
	if err != nil {
		return err
	}

	// Generate random data that will be used for the reconnect proof
	if _, err := rand.Read(c.ReconnectData); err != nil {
		return err
	}

	resp := reconnectChallengeResponse{
		Opcode: authd.OpcodeReconnectChallenge,

		// Always return success to prevent a bad actor from mining usernames.
		ErrorCode:    authd.Success,
		ChecksumSalt: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	copy(resp.ReconnectData[:], c.ReconnectData)

	respBuf := bytes.Buffer{}
	binary.Write(&respBuf, binary.BigEndian, &resp)

	if _, err := c.Conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to reconnect challenge")

	c.State = authd.StateReconnectProof

	return nil
}
