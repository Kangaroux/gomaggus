package realmd

import (
	"crypto/rand"
	"encoding/binary"
	"net"

	"github.com/kangaroux/gomaggus/model"
)

type Client struct {
	Conn          net.Conn
	ServerSeed    uint32
	Authenticated bool
	Crypto        *HeaderCrypto

	Account *model.Account
	Realm   *model.Realm
	Session *model.Session
}

func NewClient(conn net.Conn) (*Client, error) {
	seed := make([]byte, 4)
	if _, err := rand.Read(seed); err != nil {
		return nil, err
	}

	c := &Client{
		Conn:       conn,
		ServerSeed: binary.BigEndian.Uint32(seed),
	}

	return c, nil
}
