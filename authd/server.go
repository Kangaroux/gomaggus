package authd

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/jmoiron/sqlx"

	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/srp"
)

const (
	DefaultListenAddr = ":3724"
)

type Services struct {
	accounts model.AccountService
	realms   model.RealmService
	sessions model.SessionService
}

type Server struct {
	listenAddr string
	services   *Services
}

func NewServer(db *sqlx.DB, listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		services: &Services{
			accounts: model.NewDbAccountService(db),
			realms:   model.NewDbRealmService(db),
			sessions: model.NewDbSessionService(db),
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
		if err := recover(); err != nil {
			log.Printf("recovered from panic: %v", err)

			if err := conn.Close(); err != nil {
				log.Printf("error closing after recover: %v", err)
			}
		}
	}()

	log.Printf("client connected from %v\n", conn.RemoteAddr().String())

	client := &Client{
		conn:          conn,
		reconnectData: make([]byte, ReconnectDataLen),
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
	opcode := Opcode(data[0])

	switch c.state {
	case StateAuthChallenge:
		if opcode == OpLoginChallenge {
			return handleLoginChallenge(s.services, c, data)
		} else if opcode == OpReconnectChallenge {
			return handleReconnectChallenge(s.services, c, data)
		}
	case StateAuthProof:
		if opcode == OpLoginProof {
			return handleLoginProof(s.services, c, data)
		}
	case StateReconnectProof:
		if opcode == OpReconnectProof {
			return handleReconnectProof(s.services, c, data)
		}
	case StateAuthenticated:
		if opcode == OpRealmList {
			return handleRealmList(s.services, c)
		}
	}

	return fmt.Errorf(
		"handlePacket: opcode %d is not valid for current state (%d) or does not exist",
		opcode,
		c.state,
	)
}
