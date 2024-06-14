package realmd

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/kangaroux/gomaggus/internal"
	"github.com/kangaroux/gomaggus/model"
	"github.com/mixcode/binarystruct"
)

const (
	ClientHeaderSize = 6
)

type ClientHeader struct {
	Size   uint16
	Opcode ClientOpcode
}

type ClientPacket struct {
	Header  *ClientHeader
	Payload []byte
}

type Client struct {
	Conn          net.Conn
	ServerSeed    []byte
	Authenticated bool
	HeaderCrypto  *HeaderCrypto

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

// BuildHeader returns a 4-5 byte server header. The header is encrypted if the client is authenticated.
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
		if err := c.HeaderCrypto.Encrypt(header); err != nil {
			return nil, err
		}
	}

	return header, nil
}

// ParseHeader returns a parse header from a client packet. The header is decrypted if the client is
// authenticated.
func (c *Client) ParseHeader(data []byte) (*ClientHeader, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("ParseHeader: payload should be at least 6 bytes but it's only %d", len(data))
	}

	header := data[:6]

	if c.Authenticated {
		if err := c.HeaderCrypto.Decrypt(header); err != nil {
			return nil, err
		}
	}

	h := &ClientHeader{
		Size:   binary.BigEndian.Uint16(header[:2]),
		Opcode: ClientOpcode(binary.LittleEndian.Uint32(header[2:6])),
	}

	return h, nil
}

// SendPacket builds and sends a packet to the client. `data` should be either a struct pointer or
// a byte array. The data should NOT contain any header information, including size or opcode.
func (c *Client) SendPacket(opcode ServerOpcode, data interface{}) error {
	buf := bytes.Buffer{}
	if _, err := binarystruct.Write(&buf, binarystruct.LittleEndian, data); err != nil {
		return err
	}

	header, err := c.BuildHeader(opcode, uint32(buf.Len()))
	if err != nil {
		return err
	}

	_, err = c.Conn.Write(internal.ConcatBytes(header, buf.Bytes()))
	return err
}
