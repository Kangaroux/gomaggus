package server

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"net"

	"github.com/jmoiron/sqlx"
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
	"github.com/kangaroux/gomaggus/realmd/handler/auth"
	"github.com/kangaroux/gomaggus/realmd/handler/char"
	"github.com/kangaroux/gomaggus/realmd/handler/player"
	"github.com/kangaroux/gomaggus/realmd/handler/realm"
	"github.com/kangaroux/gomaggus/realmd/handler/session"
)

const (
	DefaultListenAddr = ":8085"
)

type Server struct {
	listenAddr string

	services *realmd.Service
}

func New(db *sqlx.DB, listenAddr string) *Server {
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
	if err := auth.SendChallenge(client); err != nil {
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

func (s *Server) handlePacket(c *realmd.Client, data []byte) error {
	var err error

	if len(data) == 0 {
		return fmt.Errorf("handlePacket: packet is empty")
	}

	header, err := realmd.ParseClientHeader(c, data)
	if err != nil {
		return err
	}

	switch header.Opcode {
	case realmd.OpClientAuthSession:
		return auth.ProofHandler(s.services, c, data)

	case realmd.OpClientPing:
		return session.PingHandler(c, data)

	case realmd.OpClientReadyForAccountDataTimes:
		return session.DataTimesHandler(c)

	case realmd.OpClientCharList:
		return char.ListHandler(s.services, c)

	case realmd.OpClientRealmSplit:
		return realm.SplitInfoHandler(c, data)

	case realmd.OpClientCharCreate:
		return char.CreateHandler(s.services, c, data)

	case realmd.OpClientCharDelete:
		return char.DeleteHandler(s.services, c, data)

	case realmd.OpClientPlayerLogin:
		return player.LoginHandler(s.services, c, data)

	default:
		log.Printf("unknown opcode: 0x%x\n", header.Opcode)
	}

	return nil
}
