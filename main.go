package main

import (
	"bytes"
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

func handleConnection(c Client) {
	defer func() {
		c.conn.Close()
		log.Printf("Client %d disconnected", c.id)
	}()

	log.Printf("Client %d connected from %v", c.id, c.conn.RemoteAddr())
	buf := bytes.Buffer{}

	for {
		n, err := buf.ReadFrom(c.conn)
		if err != nil {
			log.Printf("Client %d read failed: %v", c.id, err)
			return
		} else if n == 0 {
			log.Printf("Client %d closed connection", c.id)
			return
		}

		log.Printf("Client %d read %d bytes", c.id, n)
		log.Printf("%v", buf.Bytes())
		opcode, err := buf.ReadByte()

		if err != nil {
			log.Printf("Client %d failed to get opcode: %v", c.id, err)
			return
		}

		log.Printf("Client %d opcode: 0x%x", c.id, opcode)

		buf.Reset()
	}
}
