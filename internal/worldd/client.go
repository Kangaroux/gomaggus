package worldd

import "net"

type Client struct {
	conn          net.Conn
	username      string
	serverSeed    uint32
	authenticated bool
	crypto        *WrathHeaderCrypto
}
