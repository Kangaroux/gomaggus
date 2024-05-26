package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

var nextId int64 = 1

type Client struct {
	conn            net.Conn
	id              int64
	connectedAt     time.Time
	log             log.Logger
	serverPublicKey BigInteger
	verifier        BigInteger
}

func NewClient(conn net.Conn) *Client {
	id := nextId
	nextId++

	return &Client{
		conn:        conn,
		id:          id,
		connectedAt: time.Now(),
		log:         *log.New(log.Writer(), fmt.Sprintf("Client %d: ", id), log.Flags()|log.Lmsgprefix),
	}
}
