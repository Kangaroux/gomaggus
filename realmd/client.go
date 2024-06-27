package realmd

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/kangaroux/gomaggus/internal"
	"github.com/kangaroux/gomaggus/model"
	"github.com/mixcode/binarystruct"
	"github.com/phuslu/log"
)

const (
	ClientHeaderSize = 6
)

var (
	nextID atomic.Int64
)

type ClientHeader struct {
	Size   uint16
	Opcode ClientOpcode
}

type Client struct {
	ID            int64
	Conn          net.Conn
	ServerSeed    []byte
	Authenticated bool
	Log           *log.Logger

	// HeaderCrypto decrypts incoming packet headers and encrypts outgoing packet headers. HeaderCrypto
	// is nil if the client has not yet authenticated.
	HeaderCrypto *HeaderCrypto

	// Cancels a pending logout, if there is one. This func is safe to call when there is no pending logout.
	CancelPendingLogout context.CancelFunc
	LogoutPending       bool

	Account   *model.Account
	Character *model.Character
	Realm     *model.Realm
	Session   *model.Session
}

func NewClient(conn net.Conn) (*Client, error) {
	seed := make([]byte, 4)
	if _, err := rand.Read(seed); err != nil {
		return nil, err
	}

	c := &Client{
		ID:         nextID.Add(1),
		Conn:       conn,
		ServerSeed: seed,
		Log:        &log.Logger{},

		// Use a placeholder func so the caller doesn't have to check if it's nil
		CancelPendingLogout: internal.DoNothing,
	}

	return c, nil
}

// BuildHeader returns the server header as a byte array. The returned array contains 4 or 5 bytes
// depending on the size and is encrypted if the client is authenticated. If size is larger than
// SizeFieldMaxValue, BuildHeader returns an error.
func (c *Client) BuildHeader(opcode ServerOpcode, size uint32) ([]byte, error) {
	// Include the opcode in the size
	size += 2

	if size > SizeFieldMaxValue {
		return nil, fmt.Errorf("BuildHeader: size is too large (%d bytes)", size)
	}

	var header []byte

	// The size field in the header can be 2 or 3 bytes. If the size field is 3 bytes, the MSB of the
	// size will be set.
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

// ParseHeader parses and returns the header from data. If data is smaller than 6 bytes, ParseHeader
// returns an error.
func (c *Client) ParseHeader(data []byte) (*ClientHeader, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("ParseHeader: payload should be at least 6 bytes but it's only %d", len(data))
	}

	header := data[:ClientHeaderSize]

	if c.Authenticated {
		if err := c.HeaderCrypto.Decrypt(header); err != nil {
			return nil, err
		}
	}

	h := &ClientHeader{
		// The value of size in the packet includes the opcode, however for packet parsing the opcode
		// is considered part of the header. Subtracting 4 from the size gives the correct size for
		// everything after the header.
		Size:   binary.BigEndian.Uint16(header[:2]) - 4,
		Opcode: ClientOpcode(binary.LittleEndian.Uint32(header[2:6])),
	}

	return h, nil
}

// SendPacket encodes data and sends it to the client. SendPacket expects data to not contain header information.
func (c *Client) SendPacket(opcode ServerOpcode, data interface{}) error {
	var dataBytes []byte

	if data != nil {
		buf := bytes.Buffer{}
		if _, err := binarystruct.Write(&buf, binarystruct.LittleEndian, data); err != nil {
			return err
		}

		dataBytes = buf.Bytes()
	}

	return c.SendPacketBytes(opcode, dataBytes)
}

// SendPacketBytes generates a header and sends a packet containing the header + data. In most cases,
// SendPacket should be used instead.
func (c *Client) SendPacketBytes(opcode ServerOpcode, data []byte) error {
	header, err := c.BuildHeader(opcode, uint32(len(data)))
	if err != nil {
		return err
	}

	c.Log.Debug().Str("op", opcode.String()).Int("size", len(data)).Msg("packet send")
	c.Log.Trace().Str("data", hex.EncodeToString(data)).Msg("send data")

	_, err = c.Conn.Write(append(header, data...))
	return err
}
