package main

import (
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

func handlePacket(c Client, data []byte, n int) {
	if n == 0 {
		return
	}

	log.Printf("Client %d read %d bytes", c.id, n)
	log.Printf("%v", data)
	opcode := data[0]
	log.Printf("Client %d opcode: 0x%x", c.id, opcode)
}

func handleConnection(c Client) {
	defer func() {
		c.conn.Close()
		log.Printf("Client %d disconnected", c.id)
	}()

	log.Printf("Client %d connected from %v", c.id, c.conn.RemoteAddr())
	buf := make([]byte, 4096)

	for {
		n, err := c.conn.Read(buf)

		if err != nil && err != io.EOF {
			log.Printf("Client %d read failed: %v", c.id, err)
			return
		}

		handlePacket(c, buf[:n], n)

		if err == io.EOF {
			log.Printf("Client %d closed connection (EOF)", c.id)
			return
		}
	}
}
