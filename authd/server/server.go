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

	srp "github.com/kangaroux/go-wow-srp6"
	"github.com/kangaroux/gomaggus/authd"
	"github.com/kangaroux/gomaggus/authd/handler"
	"github.com/kangaroux/gomaggus/model"
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
	log.Println("listening on", listener.Addr())

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go srv.handleConnection(conn)
	}
}

func (srv *Server) handleConnection(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("recovered from panic:", err)
			debug.PrintStack()

			if err := conn.Close(); err != nil {
				log.Println("error closing after recover:", err)
			}
		}
	}()

	log.Println("client connected from", conn.RemoteAddr())

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
			log.Println("error reading from client:", readErr)
			return
		} else if readN == 0 {
			log.Println("client disconnected (closed by client)")
			return
		}

		log.Printf("read %d bytes", readN)
		buf.Write(chunk[:readN])

		handleN, err := srv.handlePacket(client, buf.Bytes())
		if err == handler.ErrPacketReadEOF {
			log.Println("handler wants more data, reading...")
			continue
		} else if err != nil {
			log.Println("error handling packet:", err)
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

// handlePacket parses and handles data sent by the client. It returns the number of bytes that were
// parsed. If the packet is incomplete and needs more data, handlePacket returns handler.ErrPacketReadEOF.
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
			return 0, err
		}

		return n, h.Handle()

	case authd.OpcodeLoginProof:
		h := handler.LoginProof{
			Client:   c,
			Sessions: srv.Sessions,
		}

		n, err := h.Read(data)
		if err != nil {
			return 0, err
		}

		return n, h.Handle()

	case authd.OpcodeReconnectChallenge:
		h := handler.ReconnectChallenge{
			Client:   c,
			Accounts: srv.Accounts,
		}

		n, err := h.Read(data)
		if err != nil {
			return 0, err
		}

		return n, h.Handle()

	case authd.OpcodeReconnectProof:
		h := handler.ReconnectProof{
			Client:   c,
			Sessions: srv.Sessions,
		}

		n, err := h.Read(data)
		if err != nil {
			return 0, err
		}

		return n, h.Handle()

	case authd.OpcodeRealmList:
		h := handler.RealmList{
			Client: c,
			Realms: srv.Realms,
		}

		n, err := h.Read(data)
		if err != nil {
			return 0, err
		}

		return n, h.Handle()

	default:
		return 0, fmt.Errorf("handlePacket: unknown opcode %x", opcode)
	}
}
