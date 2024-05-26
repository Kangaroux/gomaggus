package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const (
	AuthLoginChallengeOpCode uint8 = 0x0
	AuthLoginProofOpCode     uint8 = 0x1
)

const (
	WOW_SUCCESS                              uint8 = 0x00
	WOW_FAIL_BANNED                          uint8 = 0x03
	WOW_FAIL_UNKNOWN_ACCOUNT                 uint8 = 0x04
	WOW_FAIL_INCORRECT_PASSWORD              uint8 = 0x05
	WOW_FAIL_ALREADY_ONLINE                  uint8 = 0x06
	WOW_FAIL_NO_TIME                         uint8 = 0x07
	WOW_FAIL_DB_BUSY                         uint8 = 0x08
	WOW_FAIL_VERSION_INVALID                 uint8 = 0x09
	WOW_FAIL_VERSION_UPDATE                  uint8 = 0x0A
	WOW_FAIL_INVALID_SERVER                  uint8 = 0x0B
	WOW_FAIL_SUSPENDED                       uint8 = 0x0C
	WOW_FAIL_FAIL_NOACCESS                   uint8 = 0x0D
	WOW_SUCCESS_SURVEY                       uint8 = 0x0E
	WOW_FAIL_PARENTCONTROL                   uint8 = 0x0F
	WOW_FAIL_LOCKED_ENFORCED                 uint8 = 0x10
	WOW_FAIL_TRIAL_ENDED                     uint8 = 0x11
	WOW_FAIL_USE_BATTLENET                   uint8 = 0x12
	WOW_FAIL_ANTI_INDULGENCE                 uint8 = 0x13
	WOW_FAIL_EXPIRED                         uint8 = 0x14
	WOW_FAIL_NO_GAME_ACCOUNT                 uint8 = 0x15
	WOW_FAIL_CHARGEBACK                      uint8 = 0x16
	WOW_FAIL_INTERNET_GAME_ROOM_WITHOUT_BNET uint8 = 0x17
	WOW_FAIL_GAME_ACCOUNT_LOCKED             uint8 = 0x18
	WOW_FAIL_UNLOCKABLE_LOCK                 uint8 = 0x19
	WOW_FAIL_CONVERSION_REQUIRED             uint8 = 0x20
	WOW_FAIL_DISCONNECTED                    uint8 = 0xFF
)

func main() {
	listener, err := net.Listen("tcp", ":3724")

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	log.Print("Listening on port 3724")

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(NewClient(conn))
	}
}

func handleLoginChallenge(c *Client, data []byte) error {
	c.log.Print("Received client login challenge")

	p := LoginChallengePacket{}
	reader := bytes.NewReader(data)
	err := binary.Read(reader, binary.LittleEndian, &p)

	if err != nil {
		return err
	}

	accountNameBytes := make([]byte, p.AccountNameLength)

	if _, err := reader.Read(accountNameBytes); err != nil {
		return err
	}

	accountName := string(accountNameBytes)

	// WoW client sends these strings reversed (uint32 converted to little endian?)
	ReverseBytes(p.OSArch[:])
	ReverseBytes(p.OS[:])
	ReverseBytes(p.Locale[:])

	// c.log.Printf("GameName: %s", string(p.GameName[:4]))
	// c.log.Printf("Version: %d.%d.%d.%d", p.Version[0], p.Version[1], p.Version[2], p.Build)
	// c.log.Printf("OSArch: %s", string(p.OSArch[:4]))
	// c.log.Printf("OS: %s", string(p.OS[:4]))
	// c.log.Printf("Locale: %s", string(p.Locale[:4]))
	// c.log.Printf("IP4: %v", netip.AddrFrom4(p.IP))
	// c.log.Printf("AccountNameLength: %v", p.AccountNameLength)
	// c.log.Printf("AccountName: %v", accountName)

	loginChallengeResponse(c, accountName)

	return nil
}

func loginChallengeResponse(c *Client, username string) error {
	buf := bytes.Buffer{}
	buf.WriteByte(AuthLoginChallengeOpCode)
	buf.WriteByte(0) // protocol version

	if username != "TEST" {
		c.log.Printf("Unknown username '%s', disconnecting\n", username)

		buf.WriteByte(WOW_FAIL_UNKNOWN_ACCOUNT)
		c.conn.Write(buf.Bytes())
		c.conn.Close()
		return nil
	}

	buf.WriteByte(WOW_SUCCESS)
	c.verifier = passVerify(MOCK_PASSWORD, MOCK_PASSWORD, MOCK_SALT)
	c.serverPublicKey = calcServerPublicKey(c.verifier, MOCK_PRIVATE_KEY)
	buf.Write(c.serverPublicKey.Bytes())
	buf.WriteByte(1) // generator length
	buf.WriteByte(bigG().Bytes()[0])
	buf.WriteByte(32) // N length
	buf.Write(largeSafePrime.LittleEndian().Bytes())
	buf.Write(MOCK_SALT.Bytes())
	buf.Write([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) // CRC hash (unused)
	buf.WriteByte(0)                                                  // security flag (none)

	c.log.Println("Username OK, requesting client login proof")
	if _, err := c.conn.Write(buf.Bytes()); err != nil {
		c.log.Fatal(err)
	}

	return nil
}

func handleLoginProof(c *Client, data []byte) error {
	c.log.Print("Received client login proof")

	p := LoginProofPacket{}
	reader := bytes.NewReader(data)
	err := binary.Read(reader, binary.LittleEndian, &p)

	if err != nil {
		return err
	}

	// c.log.Printf("PublicKey: %x\n", p.ClientPublicKey)
	// c.log.Printf("Proof: %x\n", p.ClientProof)

	loginProofResponse(c, NewByteArray(p.ClientPublicKey[:], 32, false).BigInt(), NewByteArray(p.ClientProof[:], 20, false))

	return nil
}

func loginProofResponse(c *Client, clientPublicKey BigInteger, clientProof *ByteArray) error {
	buf := bytes.Buffer{}
	buf.WriteByte(AuthLoginProofOpCode)

	sessionKey := calcServerSessionKey(
		clientPublicKey, c.serverPublicKey, c.verifier, MOCK_PRIVATE_KEY,
	)
	expectedProof := calcClientProof(
		MOCK_USERNAME, sessionKey, clientPublicKey, c.serverPublicKey, MOCK_SALT)

	if clientProof != expectedProof {
		c.log.Println("Client proof does not match, disconnecting")

		buf.Write([]byte{WOW_FAIL_INCORRECT_PASSWORD, 0, 0}) // Add 2 bytes of padding so packet is 4 bytes
		c.conn.Write(buf.Bytes())
		c.conn.Close()
		return nil
	}

	serverProof := calcServerProof(clientPublicKey, clientProof, sessionKey)
	buf.WriteByte(WOW_SUCCESS)
	buf.Write(serverProof.Bytes())
	buf.Write([]byte{0, 0, 0, 0}) // Account flag (uint32)
	buf.Write([]byte{0, 0, 0, 0}) // Hardware survey ID (uint32, unused)
	buf.Write([]byte{0, 0})       // Unknown flags (uint16, unused)

	c.log.Println("Client proof OK, sending server proof")
	c.conn.Write(buf.Bytes())

	return nil
}

func handlePacket(c *Client, data []byte, dataLen int) error {
	if dataLen == 0 {
		return nil
	}

	// c.log.Printf("read %d bytes", dataLen)
	// c.log.Printf("%v", data)

	opcode := data[0]

	// c.log.Printf("opcode: 0x%x", opcode)

	switch opcode {
	case AuthLoginChallengeOpCode:
		return handleLoginChallenge(c, data)
	case AuthLoginProofOpCode:
		return handleLoginProof(c, data)
	default:
		return fmt.Errorf("unknown opcode: 0x%x", opcode)
	}
}

func handleConnection(c *Client) {
	defer func() {
		c.conn.Close()
		c.log.Print("disconnected")
	}()

	c.log.Printf("connected from %v", c.conn.RemoteAddr())
	buf := make([]byte, 4096)

	time.AfterFunc(time.Second*5, func() { c.conn.Close() })

	for {
		n, err := c.conn.Read(buf)

		if err != nil && err != io.EOF {
			c.log.Printf("read failed: %v", err)
			return
		}

		if err := handlePacket(c, buf[:n], n); err != nil {
			c.log.Printf("handle packet failed: %v", err)
			return
		}

		if err == io.EOF {
			c.log.Print("closed connection (EOF)")
			return
		}
	}
}
