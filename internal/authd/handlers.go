package authd

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"log"
	mrand "math/rand"
	"time"

	"github.com/kangaroux/gomaggus/internal"
	"github.com/kangaroux/gomaggus/internal/models"
	"github.com/kangaroux/gomaggus/internal/srp"
	"github.com/mixcode/binarystruct"
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

func handleLoginProof(services *Services, c *Client, data []byte) error {
	log.Println("Starting login proof")

	var serverProof []byte
	authenticated := false

	if c.account != nil {
		p := ClientLoginProof{}
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
		resp := ServerLoginProofFail{
			Opcode:    OpLoginProof,
			ErrorCode: CodeFailUnknownAccount,
		}
		binary.Write(&respBuf, binary.BigEndian, &resp)
	} else {
		resp := ServerLoginProofSuccess{
			Opcode:           OpLoginProof,
			ErrorCode:        CodeSuccess,
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

func handleRealmList(services *Services, c *Client) error {
	realmList, err := services.realms.List()
	if err != nil {
		return err
	}

	respBody := ServerRealmListBody{
		NumRealms: uint16(len(realmList)),
		Realms:    make([]ServerRealm, len(realmList)),
	}

	for i, r := range realmList {
		respBody.Realms[i] = ServerRealm{
			Type:          r.Type,
			Locked:        false,
			Flags:         RealmFlagNone,
			Name:          r.Name,
			Host:          r.Host,
			Population:    0, // TODO
			NumCharacters: 0, // TODO
			Region:        r.Region,
			Id:            byte(r.Id),
		}
	}

	bodyBytes, err := binarystruct.Marshal(&respBody, binarystruct.LittleEndian)
	if err != nil {
		return err
	}

	respHeader := ServerRealmListHeader{
		Opcode: OpRealmList,
		Size:   uint16(len(bodyBytes)),
	}

	headerBytes, err := binarystruct.Marshal(&respHeader, binarystruct.LittleEndian)
	if err != nil {
		return err
	}

	respBuf := bytes.Buffer{}
	respBuf.Write(headerBytes)
	respBuf.Write(bodyBytes)

	if _, err := c.conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to realm list")

	return nil
}

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

			p := ClientReconnectProof{}
			if err := p.Read(data); err != nil {
				return err
			}

			serverProof := srp.CalculateReconnectProof(c.username, p.ProofData[:], c.reconnectData, c.sessionKey)
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
