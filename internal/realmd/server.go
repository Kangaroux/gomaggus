package realmd

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kangaroux/gomaggus/internal"
	"github.com/kangaroux/gomaggus/internal/models"
	"github.com/kangaroux/gomaggus/srp"
)

const (
	DefaultPort = 3724
)

type Server struct {
	port int

	// Maps usernames to session keys to allow reconnecting.
	// FIXME?: clients can't reconnect if the realmd server restarts since this isn't persisted
	sessionKeys map[string][]byte

	accountsDb models.AccountService
	realmsDb   models.RealmService
	sessionsDb models.SessionService
}

func NewServer(db *sqlx.DB, port int) *Server {
	return &Server{
		port:        port,
		sessionKeys: make(map[string][]byte),
		accountsDb:  models.NewDbAccountService(db),
		realmsDb:    models.NewDbRealmService(db),
		sessionsDb:  models.NewDbSessionService(db),
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	log.Printf("listening on port %d\n", s.port)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("client connected from %s\n", conn.RemoteAddr().String())

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	client := &Client{
		conn:          conn,
		reconnectData: make([]byte, 16),
		privateKey:    make([]byte, srp.KeySize),
	}

	if _, err := rand.Read(client.privateKey); err != nil {
		return
	}

	buf := make([]byte, 4096)

	for {
		n, err := client.conn.Read(buf)

		if err == io.EOF {
			log.Println("client disconnected (closed by client)")
			return
		} else if err != nil {
			log.Printf("error reading from client: %v\n", err)
			return
		}

		log.Printf("read %d bytes\n", n)

		if err := s.handlePacket(client, buf[:n]); err != nil {
			log.Println(err)
			return
		}
	}
}

func (s *Server) handlePacket(c *Client, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("handlePacket: packet is empty")
	}

	opcode := data[0]

	switch c.state {
	case StateAuthChallenge:
		if opcode == OP_LOGIN_CHALLENGE {
			return s.handleLoginChallenge(c, data)
		} else if opcode == OP_RECONNECT_CHALLENGE {
			return s.handleReconnectChallenge(c, data)
		}
	case StateAuthProof:
		if opcode == OP_LOGIN_PROOF {
			return s.handleLoginProof(c, data)
		}
	case StateReconnectProof:
		if opcode == OP_RECONNECT_PROOF {
			return s.handleReconnectProof(c, data)
		}
	case StateAuthenticated:
		if opcode == OP_REALM_LIST {
			return s.handleRealmList(c)
		}
	}

	return fmt.Errorf(
		"handlePacket: opcode %d is not valid for current state (%d) or does not exist",
		opcode,
		c.state,
	)
}

func (s *Server) handleLoginChallenge(c *Client, data []byte) error {
	log.Println("Starting login challenge")

	var err error
	p := LoginChallengePacket{}

	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &p); err != nil {
		return err
	}
	usernameBytes := make([]byte, p.UsernameLength)
	if _, err := reader.Read(usernameBytes); err != nil {
		return err
	}
	c.username = strings.ToUpper(string(usernameBytes))
	log.Printf("client trying to login as '%s'\n", c.username)

	c.account, err = s.accountsDb.Get(&models.AccountGetParams{Username: c.username})
	if err != nil {
		return err
	}
	if err = c.account.DecodeSrp(); err != nil {
		return err
	}

	// https://gtker.com/wow_messages/docs/cmd_auth_logon_challenge_server.html#protocol-version-8
	resp := &bytes.Buffer{}
	resp.WriteByte(OP_LOGIN_CHALLENGE)
	resp.WriteByte(0) // protocol version

	// Always return success to prevent a bad actor from mining usernames
	resp.WriteByte(WOW_SUCCESS)

	var publicKey []byte
	var salt []byte

	if c.account == nil {
		publicKey = make([]byte, srp.KeySize)
		if _, err := rand.Read(publicKey); err != nil {
			return err
		}
		salt = make([]byte, srp.SaltSize)
		if _, err := rand.Read(salt); err != nil {
			return err
		}
	} else {
		publicKey = srp.CalculateServerPublicKey(c.account.Verifier(), c.privateKey)
		c.serverPublicKey = publicKey
		salt = c.account.Salt()
	}

	resp.WriteByte(WOW_SUCCESS)
	resp.Write(publicKey)
	resp.WriteByte(1)  // generator size (1 byte)
	resp.WriteByte(7)  // generator
	resp.WriteByte(32) // large prime size (32 bytes)
	resp.Write(internal.Reverse(srp.LargeSafePrime()))
	resp.Write(salt)
	resp.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) // crc hash
	resp.WriteByte(0)

	if _, err := c.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to login challenge")
	c.state = StateAuthProof

	return nil
}

func (s *Server) handleLoginProof(c *Client, data []byte) error {
	log.Println("Starting login proof")

	var serverProof []byte
	authenticated := false

	if c.account != nil {
		p := LoginProofPacket{}

		reader := bytes.NewReader(data)
		if err := binary.Read(reader, binary.LittleEndian, &p); err != nil {
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

	// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_server.html#protocol-version-8
	resp := &bytes.Buffer{}
	resp.WriteByte(OP_LOGIN_PROOF)

	if !authenticated {
		resp.WriteByte(WOW_FAIL_UNKNOWN_ACCOUNT)
		resp.Write([]byte{0, 0}) // padding
	} else {
		resp.WriteByte(WOW_SUCCESS)
		resp.Write(serverProof)
		resp.Write([]byte{0, 0, 0, 0}) // Account flag
		resp.Write([]byte{0, 0, 0, 0}) // Hardware survey ID
		resp.Write([]byte{0, 0})       // Unknown
	}

	if _, err := c.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to login proof")

	if authenticated {
		c.state = StateAuthenticated
		s.sessionsDb.UpdateOrCreate(&models.Session{
			AccountId:     c.account.Id,
			SessionKeyHex: hex.EncodeToString(c.sessionKey),
			Connected:     1,
			ConnectedAt:   sql.NullTime{Time: time.Now(), Valid: true},
		})
	} else {
		c.state = StateInvalid
	}

	return nil
}

func (s *Server) handleReconnectChallenge(c *Client, data []byte) error {
	log.Println("Starting reconnect challenge")

	p := LoginChallengePacket{}

	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &p); err != nil {
		return err
	}
	usernameBytes := make([]byte, p.UsernameLength)
	if _, err := reader.Read(usernameBytes); err != nil {
		return err
	}
	c.username = strings.ToUpper(string(usernameBytes))
	log.Printf("client trying to login as '%s'\n", c.username)

	// Generate random data that will be used for the reconnect proof
	if _, err := rand.Read(c.reconnectData); err != nil {
		return err
	}

	sessionKey, hasSessionKey := s.sessionKeys[c.username]
	canReconnect := c.username == MOCK_USERNAME && hasSessionKey

	// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_challenge_server.html#protocol-version-8
	resp := &bytes.Buffer{}
	resp.WriteByte(OP_RECONNECT_CHALLENGE)

	if canReconnect {
		resp.WriteByte(WOW_SUCCESS)
		resp.Write(c.reconnectData)
		resp.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) // checksum salt

		c.sessionKey = sessionKey
	} else {
		resp.WriteByte(WOW_FAIL_UNKNOWN_ACCOUNT)
	}

	if _, err := c.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to reconnect challenge")

	if canReconnect {
		c.state = StateReconnectProof
	} else {
		c.state = StateInvalid
	}

	return nil
}

func (s *Server) handleReconnectProof(c *Client, data []byte) error {
	log.Println("Starting reconnect proof")

	p := ReconnectProofPacket{}

	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &p); err != nil {
		return err
	}

	serverProof := srp.CalculateReconnectProof(c.username, p.ProofData[:], c.reconnectData, c.sessionKey)
	proofMatch := bytes.Equal(serverProof, p.ClientProof[:])

	// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_server.html#protocol-version-8
	resp := &bytes.Buffer{}
	resp.WriteByte(OP_RECONNECT_PROOF)
	if !proofMatch {
		resp.WriteByte(WOW_FAIL_UNKNOWN_ACCOUNT)
		resp.Write([]byte{0, 0}) // padding
	} else {
		resp.WriteByte(WOW_SUCCESS)
		resp.Write(serverProof)
		resp.Write([]byte{0, 0, 0, 0}) // Account flag
		resp.Write([]byte{0, 0, 0, 0}) // Hardware survey ID
		resp.Write([]byte{0, 0})       // Unknown
	}

	if _, err := c.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to reconnect proof")

	if proofMatch {
		c.state = StateAuthenticated
	} else {
		c.state = StateInvalid
	}

	return nil
}

func (s *Server) handleRealmList(c *Client) error {
	realmList, err := s.realmsDb.List()
	if err != nil {
		return err
	}

	// https://gtker.com/wow_messages/docs/cmd_realm_list_server.html#protocol-version-8
	resp := &bytes.Buffer{}
	resp.WriteByte(OP_REALM_LIST)

	inner := &bytes.Buffer{}
	inner.Write([]byte{0, 0, 0, 0}) // header padding
	binary.Write(inner, binary.LittleEndian, uint16(len(realmList)))
	for _, r := range realmList {
		inner.WriteByte(byte(r.Type))
		inner.WriteByte(0)                    // locked
		inner.WriteByte(byte(REALMFLAG_NONE)) // TODO?
		inner.WriteString(r.Name)
		inner.WriteByte(0) // name is NUL-terminated
		inner.WriteString(r.Host)
		inner.WriteByte(0)                                   // host is NUL-terminated
		binary.Write(inner, binary.LittleEndian, float32(0)) // TODO: population
		inner.WriteByte(byte(0))                             // TODO: number of chars on realm
		inner.WriteByte(byte(r.Region))
		inner.WriteByte(byte(r.Id))
	}
	inner.Write([]byte{0, 0}) // footer padding

	// Write size of realm list payload
	binary.Write(resp, binary.LittleEndian, uint16(inner.Len()))
	// Concat to main payload
	inner.WriteTo(resp)

	if _, err := c.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to realm list")

	return nil
}
