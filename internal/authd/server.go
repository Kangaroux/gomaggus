package authd

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/jmoiron/sqlx"
	"github.com/kangaroux/gomaggus/internal/models"
	"github.com/kangaroux/gomaggus/internal/srp"
)

const (
	DefaultListenAddr = ":3724"
)

type Services struct {
	accounts models.AccountService
	realms   models.RealmService
	sessions models.SessionService
}

type Server struct {
	listenAddr string
	services   *Services
}

func NewServer(db *sqlx.DB, listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		services: &Services{
			accounts: models.NewDbAccountService(db),
			realms:   models.NewDbRealmService(db),
			sessions: models.NewDbSessionService(db),
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

		log.Printf("client connected from %s\n", conn.RemoteAddr().String())

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	client := &Client{
		conn:          conn,
		reconnectData: make([]byte, 16),
		privateKey:    make([]byte, srp.KeySize),
	}

	if _, err := rand.Read(client.privateKey); err != nil {
		return
	}

	buf := make([]byte, 4096)

	for {
		n, err := client.conn.Read(buf)

		if err == io.EOF {
			log.Println("client disconnected (closed by client)")
			return
		} else if err != nil {
			log.Printf("error reading from client: %v\n", err)
			return
		}

		log.Printf("read %d bytes\n", n)

		if err := s.handlePacket(client, buf[:n]); err != nil {
			log.Println(err)
			return
		}
	}
}

func (s *Server) handlePacket(c *Client, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("handlePacket: packet is empty")
	}

	opcode := data[0]

	switch c.state {
	case StateAuthChallenge:
		if opcode == OP_LOGIN_CHALLENGE {
			return handleLoginChallenge(s.services, c, data)
		} else if opcode == OP_RECONNECT_CHALLENGE {
			return handleReconnectChallenge(s.services, c, data)
		}
	case StateAuthProof:
		if opcode == OP_LOGIN_PROOF {
			return handleLoginProof(s.services, c, data)
		}
	case StateReconnectProof:
		if opcode == OP_RECONNECT_PROOF {
			return handleReconnectProof(s.services, c, data)
		}
	case StateAuthenticated:
		if opcode == OP_REALM_LIST {
			return handleRealmList(s.services, c)
		}
	}

	return fmt.Errorf(
		"handlePacket: opcode %d is not valid for current state (%d) or does not exist",
		opcode,
		c.state,
	)
}