package realmd

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"net"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
)

const (
	DefaultListenAddr = ":8085"
)

type Server struct {
	listenAddr string

	services *realmd.Service
}

func NewServer(db *sqlx.DB, listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		services: &realmd.Service{
			Accounts: model.NewDbAccountService(db),
			Chars:    model.NewDbCharacterervice(db),
			Realms:   model.NewDbRealmService(db),
			Sessions: model.NewDbSessionService(db),
		},
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp4", s.listenAddr)

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	log.Printf("listening on %s\n", listener.Addr().String())

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recovered from panic: %v", err)

			if err := conn.Close(); err != nil {
				log.Printf("error closing after recover: %v", err)
			}
		}
	}()

	log.Printf("client connected from %v\n", conn.RemoteAddr().String())

	client := &realmd.Client{Conn: conn}
	binary.BigEndian.PutUint32(client.ServerSeed[:], mrand.Uint32())

	// The server is the one who initiates the auth challenge here, unlike the login server where
	// the client is the one who initiates it
	if err := s.sendAuthChallenge(client); err != nil {
		log.Printf("error sending auth challenge: %v\n", err)
		conn.Close()
		return
	}

	buf := make([]byte, 4096)

	for {
		log.Println("waiting to read...")
		n, err := conn.Read(buf)

		if err == io.EOF {
			log.Println("client disconnected (closed by client)")
			return
		} else if err != nil {
			log.Printf("error reading from client: %v\n", err)
			conn.Close()
			return
		}

		log.Printf("read %d bytes\n", n)

		if err := s.handlePacket(client, buf[:n]); err != nil {
			log.Printf("error handling packet: %v\n", err)
			conn.Close()
			return
		}
	}
}

// https://gtker.com/wow_messages/docs/smsg_auth_challenge.html#client-version-335
func (s *Server) sendAuthChallenge(c *realmd.Client) error {
	body := &bytes.Buffer{}
	body.Write([]byte{1, 0, 0, 0}) // unknown
	binary.Write(body, binary.BigEndian, c.ServerSeed)

	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		return err
	}
	body.Write(seed) // seed, unused. This differs from the 4 byte server seed

	resp := &bytes.Buffer{}
	respHeader, err := realmd.BuildHeader(OpServerAuthChallenge, uint32(body.Len()))
	if err != nil {
		return err
	}
	resp.Write(respHeader)
	resp.Write(body.Bytes())

	if _, err := c.Conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("sent auth challenge")
	return nil
}

func parseHeader(c *realmd.Client, data []byte) (*Header, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("parseHeader: payload should be at least 6 bytes but it's only %d", len(data))
	}

	headerData := data[:6]

	if c.Authenticated {
		if c.Crypto == nil {
			return nil, errors.New("parseHeader: client is authenticated but client.crypto is nil")
		}

		headerData = c.Crypto.Decrypt(headerData)
	}

	h := &Header{
		Size:   binary.BigEndian.Uint16(headerData[:2]),
		Opcode: ClientOpcode(binary.LittleEndian.Uint32(headerData[2:6])),
	}

	return h, nil
}

func readCString(r *bytes.Reader) (string, error) {
	s := strings.Builder{}

	for {
		b, err := r.ReadByte()

		if err != nil {
			return "", err
		} else if b == 0x0 {
			break
		}

		s.WriteByte(b)
	}

	return s.String(), nil
}

func (s *Server) handlePacket(c *realmd.Client, data []byte) error {
	var err error

	if len(data) == 0 {
		return fmt.Errorf("handlePacket: packet is empty")
	}

	header, err := parseHeader(c, data)
	if err != nil {
		return err
	}

	switch header.Opcode {
	case OpClientAuthSession:
		return handleAuthSession(s.services, c, data)

	case OpClientPing:
		return handlePing(c, data)

	case OpClientReadyForAccountDataTimes:
		return handleAccountDataTimes(c)

	case OpClientCharList:
		return handleCharList(s.services, c)

	case OpClientRealmSplit:
		return handleRealmSplit(c, data)

	case OpClientCharCreate:
		return handleCharCreate(s.services, c, data)

	case OpClientCharDelete:
		return handleCharDelete(s.services, c, data)

	case OpClientPlayerLogin:
		return handlePlayerLogin(s.services, c, data)

	default:
		log.Printf("unknown opcode: 0x%x\n", header.Opcode)
	}

	return nil
}

func getPowerTypeForClass(c model.Class) PowerType {
	switch c {
	case model.ClassWarrior:
		return PowerTypeRage

	case model.ClassPaladin,
		model.ClassHunter,
		model.ClassPriest,
		model.ClassShaman,
		model.ClassMage,
		model.ClassWarlock,
		model.ClassDruid:
		return PowerTypeMana

	case model.ClassRogue:
		return PowerTypeEnergy

	default:
		log.Println("getPowerTypeForClass: got unexpected class", c)
		return PowerTypeMana
	}
}

// packGuid returns a packed *little-endian* representation of an 8-byte integer. The packing works
// by creating a bit mask to mark which bytes are non-zero. Any bytes which are zero are discarded.
// The result is a byte array with the first byte as the bitmask, followed by the remaining
// undiscarded bytes. The bytes after the bitmask are little-endian.
func packGuid(val uint64) []byte {
	// At its largest, a packed guid takes up 9 bytes (1 byte mask + 8 bytes)
	result := make([]byte, 9)
	n := 0

	for i := 0; i < 8; i++ {
		if byte(val) > 0 {
			// Set the mask bit
			result[0] |= 1 << i
			// Add the byte to the result. The loop traverses the bytes from right-to-left but they
			// are written to the result from left-to-right, which swaps it to little-endian.
			result[1] = byte(val)
			n++
		}
		// Move to the next byte
		val >>= 8
	}

	return result[:n+1]
}

type UpdateMask struct {
	largestBit int
	mask       []uint32
}

func NewUpdateMask() *UpdateMask {
	return &UpdateMask{mask: make([]uint32, 16)}
}

// Mask returns the smallest []uint32 to represent all of the mask bits that were set.
func (m *UpdateMask) Mask() []uint32 {
	largestBitIndex := m.largestBit / 32
	return m.mask[:largestBitIndex+1]
}

// SetFieldMask sets all the bits necessary for the provided field mask.
func (m *UpdateMask) SetFieldMask(fieldMask FieldMask) {
	for i := 0; i < fieldMask.Size; i++ {
		m.SetBit(int(fieldMask.Offset) + i)
	}
}

// SetBit sets the nth bit in the update mask. The bit is zero-indexed with the first bit being zero.
func (m *UpdateMask) SetBit(bit int) {
	index := bit / 32
	bitPos := bit % 32
	m.resize(index)

	if bit > m.largestBit {
		m.largestBit = bit
	}

	m.mask[index] |= 1 << bitPos
}

// Resizes the mask to fit up to n uint32s.
func (m *UpdateMask) resize(n int) {
	if len(m.mask) > n {
		return
	}

	// Grow the array exponentially
	newSize := len(m.mask)
	newSize *= newSize

	// If it's still too small just use the desired size
	if newSize < n {
		newSize = n
	}

	oldMask := m.mask
	m.mask = make([]uint32, newSize)
	copy(m.mask, oldMask)
}
