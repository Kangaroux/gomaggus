package server

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net"
	"runtime/debug"

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
	Accounts   model.AccountService
	Realms     model.RealmService
	Sessions   model.SessionService
}

func New(db *sqlx.DB, listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		Accounts:   model.NewDbAccountService(db),
		Realms:     model.NewDbRealmService(db),
		Sessions:   model.NewDbSessionService(db),
	}
}

func (srv *Server) Start() {
	listener, err := net.Listen("tcp4", srv.listenAddr)

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

		go srv.handleConnection(conn)
	}
}

func (srv *Server) handleConnection(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recovered from panic: %v\n", err)
			debug.PrintStack()

			if err := conn.Close(); err != nil {
				log.Printf("error closing after recover: %v\n", err)
			}
		}
	}()

	log.Printf("client connected from %v\n", conn.RemoteAddr().String())

	client := &authd.Client{
		Conn:          conn,
		ReconnectData: make([]byte, handler.ReconnectDataLen),
		PrivateKey:    make([]byte, srp.KeySize),
	}

	if _, err := rand.Read(client.PrivateKey); err != nil {
		return
	}

	buf := make([]byte, 4096)

	for {
		n, err := client.Conn.Read(buf)

		if err == io.EOF {
			log.Println("client disconnected (closed by client)")
			return
		} else if err != nil {
			log.Printf("error reading from client: %v\n", err)
			return
		}

		log.Printf("read %d bytes\n", n)

		if err := srv.handlePacket(client, buf[:n]); err != nil {
			log.Println(err)
			return
		}
	}
}

func (srv *Server) handlePacket(c *authd.Client, data []byte) error {
	opcode := authd.Opcode(data[0])

	switch opcode {
	case authd.OpcodeLoginChallenge:
		h := handler.LoginChallenge{
			Client:   c,
			Accounts: srv.Accounts,
		}
		return h.Handle(data)

	case authd.OpcodeLoginProof:
		h := handler.LoginProof{
			Client:   c,
			Sessions: srv.Sessions,
		}
		return h.Handle(data)

	case authd.OpcodeReconnectChallenge:
		h := handler.ReconnectChallenge{
			Client:   c,
			Accounts: srv.Accounts,
		}
		return h.Handle(data)

	case authd.OpcodeReconnectProof:
		h := handler.ReconnectProof{
			Client:   c,
			Sessions: srv.Sessions,
		}
		return h.Handle(data)

	case authd.OpcodeRealmList:
		h := handler.RealmList{
			Client: c,
			Realms: srv.Realms,
		}
		return h.Handle()

	default:
		return fmt.Errorf("handlePacket: unknown opcode %x", opcode)
	}
}
