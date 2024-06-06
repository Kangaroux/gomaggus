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
	"time"

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
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	log.Printf("listening on port %d\n", s.port)

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
		conn:   conn,
		crypto: NewWrathHeaderCrypto(nil /* TODO session key */),
	}
	binary.BigEndian.PutUint32(client.serverSeed[:], mrand.Uint32())

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
	binary.Write(body, binary.BigEndian, c.serverSeed)

	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		return err
	}
	body.Write(seed) // seed, unused. This differs from the 4 byte server seed

	resp := &bytes.Buffer{}
	respHeader, err := makeServerHeader(OP_SRV_AUTH_CHALLENGE, uint32(body.Len()))
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
	var err error

	if len(data) == 0 {
		return fmt.Errorf("handlePacket: packet is empty")
	}

	header, err := parseHeader(c, data)
	if err != nil {
		return err
	}

	switch header.Opcode {
	case OP_CL_AUTH_SESSION:
		log.Println("starting auth session")

		r := bytes.NewReader(data[6:])

		// https://gtker.com/wow_messages/docs/cmsg_auth_session.html#client-version-335
		p := AuthSessionPacket{}
		if err = binary.Read(r, binary.LittleEndian, &p.ClientBuild); err != nil {
			return err
		}
		if err = binary.Read(r, binary.LittleEndian, &p.LoginServerId); err != nil {
			return err
		}
		if p.Username, err = readCString(r); err != nil {
			return err
		}
		if err = binary.Read(r, binary.LittleEndian, &p.LoginServerType); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &p.ClientSeed); err != nil {
			return err
		}
		if err = binary.Read(r, binary.LittleEndian, &p.RegionId); err != nil {
			return err
		}
		if err = binary.Read(r, binary.LittleEndian, &p.BattlegroundId); err != nil {
			return err
		}
		if err = binary.Read(r, binary.LittleEndian, &p.RealmId); err != nil {
			return err
		}
		if err = binary.Read(r, binary.LittleEndian, &p.DOSResponse); err != nil {
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

		c.authenticated, err = s.authenticateClient(c, &p)
		if err != nil {
			return err
		}

		if !c.authenticated {
			// We can't return an error to the client due to the header encryption, just drop the connection
			return errors.New("client could not be authenticated")
		}

		inner := bytes.Buffer{}
		inner.WriteByte(RespCodeAuthOk)
		inner.Write([]byte{0, 0, 0, 0})   // billing time
		inner.WriteByte(0x0)              // billing flags
		inner.Write([]byte{0, 0, 0, 0})   // billing rested
		inner.WriteByte(ExpansionVanilla) // exp

		// https://gtker.com/wow_messages/docs/smsg_auth_response.html#client-version-335
		resp := bytes.Buffer{}
		respHeader, err := makeServerHeader(OP_SRV_AUTH_RESPONSE, uint32(inner.Len()))
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
	case OP_CL_PING:
		log.Println("starting ping")

		r := bytes.NewReader(data[6:])
		p := PingPacket{}
		if err = binary.Read(r, binary.LittleEndian, &p.SequenceId); err != nil {
			return err
		}
		if err = binary.Read(r, binary.LittleEndian, &p.RoundTripTime); err != nil {
			return err
		}

		resp := bytes.Buffer{}
		respHeader, err := makeServerHeader(OP_SRV_PONG, 4)
		if err != nil {
			return err
		}
		resp.Write(c.crypto.Encrypt(respHeader))
		binary.Write(&resp, binary.LittleEndian, p.SequenceId)

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("sent pong")

		return nil

	case OP_CL_READY_FOR_ACCOUNT_DATA_TIMES:
		log.Println("starting account data times")

		inner := bytes.Buffer{}
		binary.Write(&inner, binary.LittleEndian, uint32(time.Now().Unix()))
		inner.WriteByte(1)              // unknown, mangos uses 1
		inner.Write([]byte{0, 0, 0, 0}) // mask(?)

		resp := bytes.Buffer{}
		respHeader, err := makeServerHeader(OP_SRV_ACCOUNT_DATA_TIMES, uint32(inner.Len()))
		if err != nil {
			return err
		}
		resp.Write(c.crypto.Encrypt(respHeader))
		resp.Write(inner.Bytes())

		if _, err := c.conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("sent account data times")

		return nil

	case OP_CL_CHAR_ENUM:
		log.Println("starting character list")

		resp := bytes.Buffer{}
		respHeader, err := makeServerHeader(OP_SRV_CHAR_ENUM, 1)
		if err != nil {
			return err
		}
		resp.Write(c.crypto.Encrypt(respHeader))
		resp.WriteByte(0) // number of characters

		log.Println("sent character list")

		return nil

	default:
		log.Printf("unknown opcode: 0x%x\n", header.Opcode)
	}

	return nil
}

func (s *Server) authenticateClient(c *Client, p *AuthSessionPacket) (bool, error) {
	var err error

	if c.account, err = s.accountsDb.Get(&models.AccountGetParams{Username: p.Username}); err != nil {
		return false, err
	} else if c.account == nil {
		log.Printf("no account with username %s exists", p.Username)
		return false, nil
	}

	if c.realm, err = s.realmsDb.Get(p.RealmId); err != nil {
		return false, err
	} else if c.realm == nil {
		log.Printf("no realm with id %d exists", p.RealmId)
		return false, nil
	}

	if c.session, err = s.sessionsDb.Get(c.account.Id); err != nil {
		return false, err
	} else if c.session == nil {
		log.Printf("no session for username %s exists", c.account.Username)
		return false, nil
	}

	if err := c.session.Decode(); err != nil {
		return false, err
	}

	c.crypto = NewWrathHeaderCrypto(c.session.SessionKey())
	if err := c.crypto.Init(); err != nil {
		return false, err
	}

	proof := CalculateWorldProof(p.Username, p.ClientSeed[:], c.serverSeed[:], c.session.SessionKey())

	if !bytes.Equal(proof, p.ClientProof[:]) {
		log.Println("proofs don't match")
		log.Printf("got:    %x\n", p.ClientProof)
		log.Printf("wanted: %x\n", proof)
		return false, nil
	}

	log.Println("client authenticated successfully")

	return true, nil
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
