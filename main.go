package main

import (
	"log"
	"net"
	"time"
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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	log.Printf("Client connected: %v", conn.RemoteAddr())
	log.Print("Disconnecting client in 5 seconds...")

	time.Sleep(5 * time.Second)
	conn.Close()

	log.Print("Client disconnected")
}
