package server

import (
	"bytes"
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

	chunk := make([]byte, 256)
	buf := bytes.Buffer{}

	for {
		readN, readErr := client.Conn.Read(chunk)
		if readErr != nil && readErr != io.EOF {
			log.Printf("error reading from client: %v\n", readErr)
			return
		}

		log.Printf("read %d bytes\n", readN)

		handleN, err := srv.handlePacket(client, buf.Bytes())
		if err != nil {
			log.Println(err)
			return
		}

		// Discard the bytes that were used to handle the packet. This behaves similar to buf.Read,
		// but only commits to reading the data if the packet was handled successfully.
		io.CopyN(io.Discard, &buf, int64(handleN))

		if readErr == io.EOF {
			log.Println("client disconnected (closed by client)")
			return
		}
	}
}

func (srv *Server) handlePacket(c *authd.Client, data []byte) (int, error) {
	opcode := authd.Opcode(data[0])

	switch opcode {
	case authd.OpcodeLoginChallenge:
		h := handler.LoginChallenge{
			Client:   c,
			Accounts: srv.Accounts,
		}

		n, err := h.Read(data)
		if err != nil {
			return n, err
		}

		return n, h.Handle()

	case authd.OpcodeLoginProof:
		h := handler.LoginProof{
			Client:   c,
			Sessions: srv.Sessions,
		}
		return 0, h.Handle(data)

	case authd.OpcodeReconnectChallenge:
		h := handler.ReconnectChallenge{
			Client:   c,
			Accounts: srv.Accounts,
		}
		return 0, h.Handle(data)

	case authd.OpcodeReconnectProof:
		h := handler.ReconnectProof{
			Client:   c,
			Sessions: srv.Sessions,
		}
		return 0, h.Handle(data)

	case authd.OpcodeRealmList:
		h := handler.RealmList{
			Client: c,
			Realms: srv.Realms,
		}
		return 0, h.Handle()

	default:
		return 0, fmt.Errorf("handlePacket: unknown opcode %x", opcode)
	}
}
