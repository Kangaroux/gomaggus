package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/kangaroux/gomaggus/internal/worldd"
)

// Opcodes sent by the server
const (
	OP_SWORLD_AUTH_CHALLENGE uint16 = 0x1EC
)

// Opcodes sent by the client
const (
	OP_CWORLD_AUTH_SESSION uint32 = 0x1ED
)

func main() {
	listener, err := net.Listen("tcp", ":8085")

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	log.Print("listening on port 8085")

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go handleClient(conn)
	}
}

type Client struct {
	conn          net.Conn
	username      string
	authenticated bool
	crypto        *worldd.WrathHeaderCrypto
}

func handleClient(c net.Conn) {
	log.Printf("client connected from %v\n", c.RemoteAddr().String())

	buf := make([]byte, 4096)
	client := &Client{conn: c}

	// The server is the one who initiates the auth challenge here, unlike the login server where
	// the client is the one who initiates it
	if err := sendAuthChallenge(client); err != nil {
		log.Printf("error sending auth challenge: %v\n", err)
		c.Close()
		return
	}

	for {
		log.Println("waiting to read...")
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

		if err := handlePacket(client, buf[:n]); err != nil {
			log.Printf("error handling packet: %v\n", err)
			c.Close()
			return
		}
	}
}

// https://gtker.com/wow_messages/docs/smsg_auth_challenge.html#client-version-335
func sendAuthChallenge(c *Client) error {
	body := &bytes.Buffer{}
	body.Write([]byte{1, 0, 0, 0})                              // unknown
	binary.Write(body, binary.LittleEndian, uint32(0xDEADBEEF)) // server seed

	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		return err
	}
	body.Write(seed) // seed, unused

	resp := &bytes.Buffer{}
	binary.Write(resp, binary.BigEndian, uint16(body.Len())+2)
	binary.Write(resp, binary.LittleEndian, OP_SWORLD_AUTH_CHALLENGE)
	body.WriteTo(resp)

	if _, err := c.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("sent auth challenge")
	return nil
}

type Header struct {
	Size   uint16
	Opcode uint32
}

func parseHeader(c *Client, data []byte) (*Header, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("parseHeader: payload should be at least 6 bytes but it's only %d", len(data))
	}

	headerData := data[:6]

	if c.authenticated {
		if c.crypto == nil {
			return nil, errors.New("parseHeader: client is authenticated but client.crypto is nil")
		}

		headerData = c.crypto.Decrypt(headerData)
	}

	h := &Header{
		Size:   binary.BigEndian.Uint16(headerData[:2]),
		Opcode: binary.LittleEndian.Uint32(headerData[2:6]),
	}

	return h, nil
}

type AuthSessionPacket struct {
	ClientBuild     uint32
	LoginServerId   uint32
	Username        string
	LoginServerType uint32
	ClientSeed      uint32
	RegionId        uint32
	BattlegroundId  uint32
	RealmId         uint32
	DOSResponse     uint64
	ClientProof     [20]byte
	AddonInfo       []byte
}

func handlePacket(c *Client, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("handlePacket: packet is empty")
	}

	header, err := parseHeader(c, data)
	if err != nil {
		return err
	}

	switch header.Opcode {
	case OP_CWORLD_AUTH_SESSION:
		// skip the header
		i := 6
		p := AuthSessionPacket{}
		p.ClientBuild = binary.BigEndian.Uint32(data[i : i+4])

		return nil
	}

	return nil
}
