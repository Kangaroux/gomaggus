package realmd

import (
	"crypto/rand"
	"fmt"
	"net"

	"github.com/kangaroux/gomaggus/model"
)

type Client struct {
	Conn          net.Conn
	ServerSeed    []byte
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
		ServerSeed: seed,
	}

	return c, nil
}

// BuildHeader returns a 4-5 byte server header. It encrypts the header if the client is authenticated.
func (c *Client) BuildHeader(opcode ServerOpcode, size uint32) ([]byte, error) {
	// Include the opcode in the size
	size += 2

	if size > SizeFieldMaxValue {
		return nil, fmt.Errorf("BuildHeader: size is too large (%d bytes)", size)
	}

	var header []byte

	// The size field in the header can be 2 or 3 bytes. The most significant bit in the size field
	// is reserved as a flag to indicate this. In total, server headers are 4 or 5 bytes.
	//
	// The header format is: <size><opcode>
	// <size> is 2-3 bytes big endian
	// <opcode> is 2 bytes little endian
	if size > LargeHeaderThreshold {
		header = []byte{
			byte(size>>16) | LargeHeaderFlag,
			byte(size >> 8),
			byte(size),
			byte(opcode),
			byte(opcode >> 8),
		}
	} else {
		header = []byte{
			byte(size >> 8),
			byte(size),
			byte(opcode),
			byte(opcode >> 8),
		}
	}

	if c.Authenticated {
		header = c.Crypto.Encrypt(header)
	}

	return header, nil
}
