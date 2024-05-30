package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	srpv2 "github.com/kangaroux/go-realmd/srp_v2"
)

const (
	OP_LOGIN_CHALLENGE     byte = 0
	OP_LOGIN_PROOF         byte = 1
	OP_RECONNECT_CHALLENGE byte = 2
	OP_RECONNECT_PROOF     byte = 3

	WOW_SUCCESS              byte = 0
	WOW_FAIL_UNKNOWN_ACCOUNT byte = 4

	MOCK_USERNAME = "TEST"
	MOCK_PASSWORD = "PASSWORD"
)

var (
	MOCK_SALT        []byte
	MOCK_VERIFIER    []byte
	MOCK_PRIVATE_KEY []byte
	MOCK_PUBLIC_KEY  []byte
)

func init() {
	MOCK_SALT = make([]byte, 32)
	if _, err := rand.Read(MOCK_SALT); err != nil {
		log.Fatalf("error generating salt: %v\n", err)
	}

	MOCK_VERIFIER = srpv2.CalculateVerifier(MOCK_USERNAME, MOCK_PASSWORD, MOCK_SALT)
	MOCK_PRIVATE_KEY = srpv2.NewPrivateKey()
	MOCK_PUBLIC_KEY = srpv2.CalculateServerPublicKey(MOCK_VERIFIER, MOCK_PRIVATE_KEY)
}

func main() {
	listener, err := net.Listen("tcp", ":3724")

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	log.Print("listening on port 3724")

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go handleClient(conn)
	}
}

func handleClient(c net.Conn) {
	buf := make([]byte, 4096)

	log.Printf("client connected from %v\n", c.RemoteAddr().String())

	for {
		n, err := c.Read(buf)

		if err == io.EOF {
			log.Println("client disconnected (closed by client)")
			return
		} else if err != nil {
			log.Printf("error reading from client: %v\n", err)
			c.Close()
			return
		}

		log.Printf("read %d bytes\n", n)

		if err := handlePacket(c, buf[:n]); err != nil {
			log.Printf("error handling packet: %v\n", err)
			c.Close()
			return
		}
	}
}

type LoginChallengePacket struct {
	Opcode         byte // 0x0
	Error          byte // unused?
	Size           uint16
	GameName       [4]byte
	Version        [3]byte
	Build          uint16
	OSArch         [4]byte
	OS             [4]byte
	Locale         [4]byte
	TimezoneBias   uint32
	IP             [4]byte
	UsernameLength uint8

	// The username is a variable size and needs to be read manually
	// Username string
}

type LoginProofPacket struct {
	Opcode           byte // 0x1
	ClientPublicKey  [32]byte
	ClientProof      [20]byte
	CRCHash          [20]byte // unused
	NumTelemetryKeys uint8    // unused
}

func handlePacket(c net.Conn, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("error: packet is empty")
	}

	switch data[0] {
	case OP_LOGIN_CHALLENGE:
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
		username := string(usernameBytes)
		log.Printf("client trying to login as '%s'\n", username)

		// https://gtker.com/wow_messages/docs/cmd_auth_logon_challenge_server.html#protocol-version-8
		resp := bytes.Buffer{}
		resp.WriteByte(OP_LOGIN_CHALLENGE)
		resp.WriteByte(0) // protocol version
		resp.WriteByte(WOW_SUCCESS)
		resp.Write(MOCK_PUBLIC_KEY)
		resp.WriteByte(1)  // generator size (1 byte)
		resp.WriteByte(7)  // generator
		resp.WriteByte(32) // large prime size (32 bytes)
		resp.Write(srpv2.Reverse(srpv2.LargeSafePrime))
		resp.Write(MOCK_SALT)
		resp.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) // crc hash
		resp.WriteByte(0)

		if _, err := c.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("Replied to login challenge")
		return nil
	case OP_LOGIN_PROOF:
		log.Println("Starting login proof")
		p := LoginProofPacket{}
		reader := bytes.NewReader(data)
		if err := binary.Read(reader, binary.LittleEndian, &p); err != nil {
			return err
		}

		clientPublicKey := p.ClientPublicKey[:]
		clientProof := p.ClientProof[:]

		sessionKey := srpv2.CalculateServerSessionKey(
			clientPublicKey, MOCK_PUBLIC_KEY, MOCK_PRIVATE_KEY, MOCK_VERIFIER)
		calculatedClientProof := srpv2.CalculateClientProof(
			MOCK_USERNAME, MOCK_SALT, clientPublicKey, MOCK_PUBLIC_KEY, sessionKey,
		)

		// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_server.html#protocol-version-8
		resp := bytes.Buffer{}
		resp.WriteByte(OP_LOGIN_PROOF)

		if !bytes.Equal(calculatedClientProof, clientProof) {
			resp.WriteByte(WOW_FAIL_UNKNOWN_ACCOUNT)
			resp.Write([]byte{0, 0}) // padding
		} else {
			resp.WriteByte(WOW_SUCCESS)
			resp.Write(srpv2.CalculateServerProof(clientPublicKey, clientProof, sessionKey))
			resp.Write([]byte{0, 0, 0, 0}) // Account flag
			resp.Write([]byte{0, 0, 0, 0}) // Hardware survey ID
			resp.Write([]byte{0, 0})       // Unknown
		}

		if _, err := c.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("Replied to login proof")
		return nil
	case OP_RECONNECT_CHALLENGE:
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
		username := string(usernameBytes)
		log.Printf("client trying to reconnect as '%s'\n", username)

		randBytes := make([]byte, 16)
		if _, err := rand.Read(randBytes); err != nil {
			return err
		}

		// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_challenge_server.html#protocol-version-8
		resp := bytes.Buffer{}
		resp.WriteByte(OP_RECONNECT_CHALLENGE)
		resp.Write(randBytes)
		resp.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) // checksum salt

		if _, err := c.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("Replied to reconnect challenge")
		return nil
	case OP_RECONNECT_PROOF:
		return nil
	default:
		return fmt.Errorf("error: unknown opcode (%v)", data[0])
	}
}
