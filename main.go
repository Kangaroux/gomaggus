package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/netip"
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
	c.log.Print("start login challenge")

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

	c.log.Printf("GameName: %s", string(p.GameName[:4]))
	c.log.Printf("Version: %d.%d.%d.%d", p.Version[0], p.Version[1], p.Version[2], p.Build)
	c.log.Printf("OSArch: %s", string(p.OSArch[:4]))
	c.log.Printf("OS: %s", string(p.OS[:4]))
	c.log.Printf("Locale: %s", string(p.Locale[:4]))
	c.log.Printf("IP4: %v", netip.AddrFrom4(p.IP))
	c.log.Printf("AccountNameLength: %v", p.AccountNameLength)
	c.log.Printf("AccountName: %v", accountName)

	return nil
}

func handleLoginProof(c *Client, data []byte) error {
	c.log.Print("start login proof")

	return nil
}

func handlePacket(c *Client, data []byte, dataLen int) error {
	if dataLen == 0 {
		return nil
	}

	c.log.Printf("read %d bytes", dataLen)
	// c.log.Printf("%v", data)

	opcode := data[0]

	c.log.Printf("opcode: 0x%x", opcode)

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
