package auth

import (
	"bytes"
	"crypto/rand"
	"log"

	"github.com/kangaroux/gomaggus/internal"
	"github.com/kangaroux/gomaggus/realmd"
	"github.com/mixcode/binarystruct"
)

// https://gtker.com/wow_messages/docs/smsg_auth_challenge.html#client-version-335
// The server initiates the challenge, there is no initial request from the client
type challengeResponse struct {
	Unknown    uint32
	ServerSeed uint32
	UnusedSeed [32]byte
}

func SendChallenge(c *realmd.Client) error {
	resp := challengeResponse{
		Unknown:    0x1,
		ServerSeed: c.ServerSeed,
	}

	// Generate the unused seed
	if _, err := rand.Read(resp.UnusedSeed[:]); err != nil {
		return err
	}

	buf := bytes.Buffer{}
	if _, err := binarystruct.Write(&buf, binarystruct.LittleEndian, &resp); err != nil {
		return err
	}

	header, err := realmd.BuildHeader(realmd.OpServerAuthChallenge, uint32(buf.Len()))
	if err != nil {
		return err
	}

	if _, err := c.Conn.Write(internal.ConcatBytes(header, buf.Bytes())); err != nil {
		return err
	}

	log.Println("sent auth challenge")
	return nil
}
