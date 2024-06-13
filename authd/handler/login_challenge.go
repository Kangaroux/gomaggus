package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"log"
	mrand "math/rand"

	"github.com/kangaroux/gomaggus/authd"
	"github.com/kangaroux/gomaggus/internal"
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/srp"
	"github.com/mixcode/binarystruct"
)

// https://gtker.com/wow_messages/docs/cmd_auth_logon_challenge_client.html
type loginChallengeRequest struct {
	Opcode          Opcode // OpLoginChallenge
	ProtocolVersion uint8
	Size            uint16
	GameName        [4]byte
	Version         [3]byte
	Build           uint16
	OSArch          [4]byte
	OS              [4]byte
	Locale          [4]byte
	TimezoneBias    uint32
	IP              [4]byte
	UsernameLength  uint8
	Username        string `binary:"string(UsernameLength)"`
}

func (p *loginChallengeRequest) Read(data []byte) error {
	reader := bytes.NewReader(data)

	if _, err := binarystruct.Read(reader, binarystruct.LittleEndian, p); err != nil {
		return err
	}

	if reader.Len() != 0 {
		return &ErrPacketUnreadBytes{Handler: "LoginChallengePacket", UnreadCount: reader.Len()}
	}

	return nil
}

// https://gtker.com/wow_messages/docs/cmd_auth_logon_challenge_server.html#protocol-version-8
type loginChallengeResponse struct {
	Opcode          Opcode
	ProtocolVersion uint8
	ErrorCode       RespCode
	PublicKey       [srp.KeySize]byte
	GeneratorSize   uint8
	Generator       uint8
	LargePrimeSize  uint8
	LargePrime      [srp.LargePrimeSize]byte
	Salt            [srp.SaltSize]byte
	CrcHash         [16]byte

	// Using any flags would require additional fields but this is set to zero for now
	SecurityFlags byte
}

func LoginChallenge(svc *authd.Service, c *authd.Client, data []byte) error {
	if c.State != authd.StateAuthChallenge {
		return &ErrWrongState{
			Handler:  "LoginChallenge",
			Expected: authd.StateAuthChallenge,
			Actual:   c.State,
		}
	}

	log.Println("Starting login challenge")

	var err error

	p := loginChallengeRequest{}
	if err = p.Read(data); err != nil {
		return err
	}
	c.Username = p.Username

	log.Printf("client trying to login as '%s'\n", c.Username)

	c.Account, err = svc.Accounts.Get(&model.AccountGetParams{Username: c.Username})
	if err != nil {
		return err
	}

	var publicKey []byte
	var salt []byte

	if c.Account == nil {
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
		seededRand := mrand.New(mrand.NewSource(internal.FastHash(c.Username)))
		salt = make([]byte, srp.SaltSize)
		if _, err := seededRand.Read(salt); err != nil {
			return err
		}
	} else {
		if err = c.Account.DecodeSrp(); err != nil {
			return err
		}
		publicKey = srp.CalculateServerPublicKey(c.Account.Verifier(), c.PrivateKey)
		c.ServerPublicKey = publicKey
		salt = c.Account.Salt()
	}

	resp := loginChallengeResponse{
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

	if _, err := c.Conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to login challenge")
	c.State = authd.StateAuthProof

	return nil
}
