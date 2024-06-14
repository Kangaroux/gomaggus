package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"log"

	"github.com/kangaroux/gomaggus/realmd"
)

// https://gtker.com/wow_messages/docs/smsg_auth_challenge.html#client-version-335
func SendChallenge(c *realmd.Client) error {
	body := &bytes.Buffer{}
	body.Write([]byte{1, 0, 0, 0}) // unknown
	binary.Write(body, binary.BigEndian, c.ServerSeed)

	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		return err
	}
	body.Write(seed) // seed, unused. This differs from the 4 byte server seed

	resp := &bytes.Buffer{}
	respHeader, err := realmd.BuildHeader(realmd.OpServerAuthChallenge, uint32(body.Len()))
	if err != nil {
		return err
	}
	resp.Write(respHeader)
	resp.Write(body.Bytes())

	if _, err := c.Conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("sent auth challenge")
	return nil
}
