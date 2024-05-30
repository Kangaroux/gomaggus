package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	OP_WORLD_AUTH_CHALLENGE uint16 = 0xEC01
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
	conn       net.Conn
	username   string
	sessionKey []byte
}

func handleClient(c net.Conn) {
	log.Printf("client connected from %v\n", c.RemoteAddr().String())

	buf := make([]byte, 4096)
	client := &Client{conn: c}

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
	binary.Write(resp, binary.BigEndian, OP_WORLD_AUTH_CHALLENGE)
	body.WriteTo(resp)

	if _, err := c.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("sent auth challenge")
	return nil
}

func handlePacket(c *Client, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("error: packet is empty")
	}

	// switch data[0] {
	// case OP_LOGIN_CHALLENGE:
	// }

	return nil
}
