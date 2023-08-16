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

func handleLoginChallenge(c *Client, data []byte, dataLen int) error {
	c.log.Print("start login challenge")

	p := LoginChallengePacket{}
	reader := bytes.NewReader(data)
	err := binary.Read(reader, binary.LittleEndian, &p)

	if err != nil {
		return err
	}

	accountName := make([]byte, p.AccountNameLength)

	if _, err := reader.Read(accountName); err != nil {
		return err
	}

	// WoW client sends these strings reversed (uint32 converted to little endian?)
	reverseBytes(p.OSArch[:], 4)
	reverseBytes(p.OS[:], 4)
	reverseBytes(p.Locale[:], 4)

	c.log.Printf("GameName: %s", string(p.GameName[:4]))
	c.log.Printf("Version: v%d.%d.%d.%d", p.Version[0], p.Version[1], p.Version[2], p.Build)
	c.log.Printf("OSArch: %s", string(p.OSArch[:4]))
	c.log.Printf("OS: %s", string(p.OS[:4]))
	c.log.Printf("Locale: %s", string(p.Locale[:4]))
	// c.log.Printf("AccountNameFirstLetter: %v", string(p.AccountNameFirstLetter))
	c.log.Printf("AccountNameLength: %v", p.AccountNameLength)
	c.log.Printf("IP4: %v", netip.AddrFrom4(p.IP))

	return nil
}

func handleLoginProof(c *Client, data []byte, dataLen int) error {
	c.log.Print("start login proof")

	return nil
}

func handlePacket(c *Client, data []byte, dataLen int) error {
	if dataLen == 0 {
		return nil
	}

	c.log.Printf("read %d bytes", dataLen)
	c.log.Printf("%v", data)

	opcode := data[0]

	c.log.Printf("opcode: 0x%x", opcode)

	switch opcode {
	case 0:
		return handleLoginChallenge(c, data, dataLen)
	case 1:
		return handleLoginProof(c, data, dataLen)
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
