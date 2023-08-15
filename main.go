package main

import (
	"fmt"
	"io"
	"log"
	"net"
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

func handleLoginChallenge(c *Client, data []byte, n int) {
	c.log.Print("start login challenge")
}

func handleLoginProof(c *Client, data []byte, n int) {
	c.log.Print("start login proof")
}

func handlePacket(c *Client, data []byte, n int) error {
	if n == 0 {
		return nil
	}

	c.log.Printf("read %d bytes", n)
	c.log.Printf("%v", data)

	opcode := data[0]

	c.log.Printf("opcode: 0x%x", opcode)

	switch opcode {
	case 0:
		handleLoginChallenge(c, data, n)
	case 1:
		handleLoginProof(c, data, n)
	default:
		return fmt.Errorf("unknown opcode: 0x%x", opcode)
	}

	return nil
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
