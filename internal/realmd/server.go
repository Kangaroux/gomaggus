package realmd

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

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

	realmsDb models.RealmService
}

func NewServer(db *sqlx.DB, port int) *Server {
	return &Server{
		port:        port,
		sessionKeys: make(map[string][]byte),
		realmsDb:    models.NewDbRealmService(db),
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

		client := &Client{conn: conn, reconnectData: make([]byte, 16)}
		go s.handleClient(client)
	}
}

func (s *Server) handleClient(c *Client) {
	buf := make([]byte, 4096)

	for {
		n, err := c.conn.Read(buf)

		if err == io.EOF {
			log.Println("client disconnected (closed by client)")
			return
		} else if err != nil {
			log.Printf("error reading from client: %v\n", err)
			c.conn.Close()
			return
		}

		log.Printf("read %d bytes\n", n)

		if err := s.handlePacket(c, buf[:n]); err != nil {
			log.Println(err)
			c.conn.Close()
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

	// https://gtker.com/wow_messages/docs/cmd_auth_logon_challenge_server.html#protocol-version-8
	resp := &bytes.Buffer{}
	resp.WriteByte(OP_LOGIN_CHALLENGE)
	resp.WriteByte(0) // protocol version

	if c.username == MOCK_USERNAME {
		resp.WriteByte(WOW_SUCCESS)
		resp.Write(MOCK_PUBLIC_KEY)
		resp.WriteByte(1)  // generator size (1 byte)
		resp.WriteByte(7)  // generator
		resp.WriteByte(32) // large prime size (32 bytes)
		resp.Write(internal.Reverse(srp.LargeSafePrime))
		resp.Write(MOCK_SALT)
		resp.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) // crc hash
		resp.WriteByte(0)
	} else {
		resp.WriteByte(WOW_FAIL_UNKNOWN_ACCOUNT)
	}

	if _, err := c.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to login challenge")
	c.state = StateAuthProof

	return nil
}

func (s *Server) handleLoginProof(c *Client, data []byte) error {
	log.Println("Starting login proof")
	p := LoginProofPacket{}
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &p); err != nil {
		return err
	}

	clientPublicKey := p.ClientPublicKey[:]
	clientProof := p.ClientProof[:]

	c.sessionKey = srp.CalculateServerSessionKey(
		clientPublicKey, MOCK_PUBLIC_KEY, MOCK_PRIVATE_KEY, MOCK_VERIFIER)
	calculatedClientProof := srp.CalculateClientProof(
		MOCK_USERNAME, MOCK_SALT, clientPublicKey, MOCK_PUBLIC_KEY, c.sessionKey,
	)
	proofMatch := bytes.Equal(calculatedClientProof, clientProof)

	// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_server.html#protocol-version-8
	resp := &bytes.Buffer{}
	resp.WriteByte(OP_LOGIN_PROOF)

	if !proofMatch {
		resp.WriteByte(WOW_FAIL_UNKNOWN_ACCOUNT)
		resp.Write([]byte{0, 0}) // padding
	} else {
		resp.WriteByte(WOW_SUCCESS)
		resp.Write(srp.CalculateServerProof(clientPublicKey, clientProof, c.sessionKey))
		resp.Write([]byte{0, 0, 0, 0}) // Account flag
		resp.Write([]byte{0, 0, 0, 0}) // Hardware survey ID
		resp.Write([]byte{0, 0})       // Unknown

		// Save the session key in case the client needs to reconnect later
		s.sessionKeys[c.username] = c.sessionKey
	}

	if _, err := c.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to login proof")

	if proofMatch {
		c.state = StateAuthenticated
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
