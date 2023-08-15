package main

import (
	"net"
	"time"
)

var nextId int64 = 1

type Client struct {
	conn        net.Conn
	id          int64
	connectedAt time.Time
}

func NewClient(conn net.Conn) Client {
	id := nextId
	nextId++

	return Client{
		conn:        conn,
		id:          id,
		connectedAt: time.Now(),
	}
}
