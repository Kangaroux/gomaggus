package auth

import (
	"crypto/rand"
	"log"

	"github.com/kangaroux/gomaggus/realmd"
)

// https://gtker.com/wow_messages/docs/smsg_auth_challenge.html#client-version-335
// The server initiates the challenge, there is no initial request from the client
type challengeResponse struct {
	Unknown    uint32
	ServerSeed []byte `binary:"[4]byte"`
	UnusedSeed [32]byte
}

func SendChallenge(client *realmd.Client) error {
	resp := challengeResponse{
		Unknown:    0x1,
		ServerSeed: client.ServerSeed,
	}

	// Generate the unused seed
	if _, err := rand.Read(resp.UnusedSeed[:]); err != nil {
		return err
	}

	if err := client.SendPacket(realmd.OpServerAuthChallenge, &resp); err != nil {
		return err
	}

	log.Println("sent auth challenge")
	return nil
}
