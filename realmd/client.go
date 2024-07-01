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

	"github.com/kangaroux/go-wow-srp6/header"
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

	// Header manages packet header encryption/decryption as well as encoding server headers.
	Header *header.WrathHeader

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
		Header:     &header.WrathHeader{},

		// Use a placeholder func so the caller doesn't have to check if it's nil
		CancelPendingLogout: internal.DoNothing,
	}

	return c, nil
}

// ParseHeader parses and returns the header from data. If data is smaller than 6 bytes, ParseHeader
// returns an error.
func (c *Client) ParseHeader(data []byte) (*ClientHeader, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("ParseHeader: payload should be at least 6 bytes but it's only %d", len(data))
	}

	header := data[:ClientHeaderSize]

	if c.Authenticated {
		if err := c.Header.Decrypt(header); err != nil {
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
	header, err := c.Header.Encode(uint16(opcode), uint32(len(data)))
	if err != nil {
		return err
	}

	c.Log.Debug().
		Str("op", opcode.String()).
		Int("size", len(data)).
		Msg("packet send")

	c.Log.Trace().
		Func(func(e *log.Entry) { // Skip encoding unless it's actually needed
			e.Str("data", hex.EncodeToString(data))
		}).
		Msg("send data")

	_, err = c.Conn.Write(append(header, data...))
	return err
}
