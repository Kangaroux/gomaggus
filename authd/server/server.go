package server

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/jmoiron/sqlx"

	"github.com/kangaroux/gomaggus/authd"
	"github.com/kangaroux/gomaggus/authd/handler"
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/srp"
)

const (
	DefaultListenAddr = ":3724"
)

type Server struct {
	listenAddr string
	services   *authd.Service
}

func New(db *sqlx.DB, listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		services: &authd.Service{
			Accounts: model.NewDbAccountService(db),
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

	c := &authd.Client{
		Conn:          conn,
		ReconnectData: make([]byte, handler.ReconnectDataLen),
		PrivateKey:    make([]byte, srp.KeySize),
	}

	if _, err := rand.Read(c.PrivateKey); err != nil {
		return
	}

	buf := make([]byte, 4096)

	for {
		n, err := c.Conn.Read(buf)

		if err == io.EOF {
			log.Println("client disconnected (closed by client)")
			return
		} else if err != nil {
			log.Printf("error reading from client: %v\n", err)
			return
		}

		log.Printf("read %d bytes\n", n)

		if err := s.handlePacket(c, buf[:n]); err != nil {
			log.Println(err)
			return
		}
	}
}

func (s *Server) handlePacket(c *authd.Client, data []byte) error {
	opcode := handler.Opcode(data[0])

	switch c.State {
	case authd.StateAuthChallenge:
		if opcode == handler.OpLoginChallenge {
			return handler.LoginChallenge(s.services, c, data)
		} else if opcode == handler.OpReconnectChallenge {
			return handler.ReconnectChallenge(s.services, c, data)
		}
	case authd.StateAuthProof:
		if opcode == handler.OpLoginProof {
			return handler.LoginProof(s.services, c, data)
		}
	case authd.StateReconnectProof:
		if opcode == handler.OpReconnectProof {
			return handler.ReconnectProof(s.services, c, data)
		}
	case authd.StateAuthenticated:
		if opcode == handler.OpRealmList {
			return handler.RealmList(s.services, c)
		}
	}

	return fmt.Errorf(
		"handlePacket: opcode %d is not valid for current state (%d) or does not exist",
		opcode,
		c.State,
	)
}
