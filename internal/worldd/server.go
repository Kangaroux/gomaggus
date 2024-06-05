package worldd

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
	"github.com/kangaroux/gomaggus/internal/models"
)

const (
	DefaultPort = 8085
)

type Server struct {
	port int

	accountsDb models.AccountService
	realmsDb   models.RealmService
	sessionsDb models.SessionService
}

func NewServer(db *sqlx.DB, port int) *Server {
	return &Server{
		port:       port,
		accountsDb: models.NewDbAccountService(db),
		realmsDb:   models.NewDbRealmService(db),
		sessionsDb: models.NewDbSessionService(db),
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", ":8085")

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	log.Print("listening on port 8085")

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

	buf := make([]byte, 4096)
	client := &Client{
		conn:       conn,
		serverSeed: mrand.Uint32(),
		crypto:     NewWrathHeaderCrypto(nil /* TODO session key */),
	}

	// The server is the one who initiates the auth challenge here, unlike the login server where
	// the client is the one who initiates it
	if err := s.sendAuthChallenge(client); err != nil {
		log.Printf("error sending auth challenge: %v\n", err)
		conn.Close()
		return
	}

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
func (s *Server) sendAuthChallenge(c *Client) error {
	body := &bytes.Buffer{}
	body.Write([]byte{1, 0, 0, 0}) // unknown
	binary.Write(body, binary.LittleEndian, c.serverSeed)

	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		return err
	}
	body.Write(seed) // seed, unused. This differs from the 4 byte server seed

	resp := &bytes.Buffer{}
	respHeader, err := makeServerHeader(OP_AUTH_CHALLENGE, uint32(body.Len()))
	if err != nil {
		return err
	}
	resp.Write(respHeader)
	resp.Write(body.Bytes())

	if _, err := c.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("sent auth challenge")
	return nil
}

func parseHeader(c *Client, data []byte) (*Header, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("parseHeader: payload should be at least 6 bytes but it's only %d", len(data))
	}

	headerData := data[:6]

	if c.authenticated {
		if c.crypto == nil {
			return nil, errors.New("parseHeader: client is authenticated but client.crypto is nil")
		}

		headerData = c.crypto.Decrypt(headerData)
	}

	h := &Header{
		Size:   binary.BigEndian.Uint16(headerData[:2]),
		Opcode: binary.LittleEndian.Uint32(headerData[2:6]),
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

func (s *Server) handlePacket(c *Client, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("handlePacket: packet is empty")
	}

	header, err := parseHeader(c, data)
	if err != nil {
		return err
	}

	switch header.Opcode {
	case OP_AUTH_SESSION:
		log.Println("starting auth session")

		r := bytes.NewReader(data)

		// Skip the header
		r.Seek(6, io.SeekStart)

		p := AuthSessionPacket{}
		if err = binary.Read(r, binary.BigEndian, &p.ClientBuild); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &p.LoginServerId); err != nil {
			return err
		}
		if p.Username, err = readCString(r); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &p.LoginServerType); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &p.ClientSeed); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &p.RegionId); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &p.BattlegroundId); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &p.RealmId); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &p.DOSResponse); err != nil {
			return err
		}
		if _, err = r.Read(p.ClientProof[:]); err != nil {
			return err
		}
		addonInfoBuf := bytes.Buffer{}
		if _, err = r.WriteTo(&addonInfoBuf); err != nil {
			return err
		}
		p.AddonInfo = addonInfoBuf.Bytes()

		// TODO: Check client proof

		// https://gtker.com/wow_messages/docs/smsg_auth_response.html#client-version-335
		inner := bytes.Buffer{}
		inner.WriteByte(RespCodeAuthOk)
		inner.Write([]byte{0, 0, 0, 0})   // billing time
		inner.WriteByte(0x0)              // billing flags
		inner.Write([]byte{0, 0, 0, 0})   // billing rested
		inner.WriteByte(ExpansionVanilla) // exp
		resp := bytes.Buffer{}
		respHeader, err := makeServerHeader(OP_AUTH_RESPONSE, uint32(inner.Len()))
		if err != nil {
			return err
		}
		resp.Write(c.crypto.Encrypt(respHeader))
		resp.Write(inner.Bytes())

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("sent auth response")

		return nil
	}

	return nil
}

func makeServerHeader(opcode uint16, size uint32) ([]byte, error) {
	// Include the opcode in the size
	size += 2

	if size > SizeFieldMaxValue {
		return nil, fmt.Errorf("makeServerHeader: size is too large (%d bytes)", size)
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

	return header, nil
}
