package authd

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"log"
	mrand "math/rand"

	"github.com/kangaroux/gomaggus/internal"

	"github.com/kangaroux/gomaggus/internal/models"
	"github.com/kangaroux/gomaggus/internal/srp"
)

func handleLoginChallenge(services *Services, c *Client, data []byte) error {
	log.Println("Starting login challenge")

	var err error

	p := ClientLoginChallenge{}
	if err = p.Read(data); err != nil {
		return err
	}
	c.username = p.Username

	log.Printf("client trying to login as '%s'\n", c.username)

	c.account, err = services.accounts.Get(&models.AccountGetParams{Username: c.username})
	if err != nil {
		return err
	}

	var publicKey []byte
	var salt []byte

	if c.account == nil {
		publicKey = make([]byte, srp.KeySize)
		if _, err := rand.Read(publicKey); err != nil {
			return err
		}

		// A real account will always return the same salt, so our fake account needs to do that, too.
		// Using the username as a seed for the fake salt guarantees we always generate the same data.
		// Ironically, using crypto/rand here is actually less secure.
		//
		// If we didn't do this, a bad actor could send two challenges with the same username and compare
		// the salts. The salts would be the same for real accounts and different for fake accounts.
		// This would allow someone to mine usernames and start building an attack vector.
		seededRand := mrand.New(mrand.NewSource(internal.FastHash(c.username)))
		salt = make([]byte, srp.SaltSize)
		if _, err := seededRand.Read(salt); err != nil {
			return err
		}
	} else {
		if err = c.account.DecodeSrp(); err != nil {
			return err
		}
		publicKey = srp.CalculateServerPublicKey(c.account.Verifier(), c.privateKey)
		c.serverPublicKey = publicKey
		salt = c.account.Salt()
	}

	resp := ServerLoginChallenge{
		Opcode: OpLoginChallenge,

		// Protocol version is always zero for server responses
		ProtocolVersion: 0,

		// Always return success to prevent a bad actor from mining usernames. See above for how
		// fake data is generated when the username doesn't exist
		ErrorCode:      CodeSuccess,
		GeneratorSize:  1,
		Generator:      srp.Generator,
		LargePrimeSize: srp.LargePrimeSize,
		SecurityFlags:  0,
	}
	copy(resp.PublicKey[:], publicKey)
	copy(resp.LargePrime[:], srp.LargePrime())
	copy(resp.Salt[:], salt)
	copy(resp.CrcHash[:], []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	respBuf := bytes.Buffer{}
	// The byte arrays are already little endian so the buffer can be used as-is
	binary.Write(&respBuf, binary.BigEndian, &resp)

	if _, err := c.conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to login challenge")
	c.state = StateAuthProof

	return nil
}
