package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"runtime/debug"

	"github.com/jmoiron/sqlx"
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
	"github.com/kangaroux/gomaggus/realmd/handler/account"
	"github.com/kangaroux/gomaggus/realmd/handler/auth"
	"github.com/kangaroux/gomaggus/realmd/handler/char"
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
	log.SetFlags(log.Lmicroseconds)
	return &Server{
		listenAddr: listenAddr,
		services: &realmd.Service{
			Accounts:       model.NewDbAccountService(db),
			AccountStorage: model.NewDbAccountStorageService(db),
			Chars:          model.NewDbCharacterService(db),
			Realms:         model.NewDbRealmService(db),
			Sessions:       model.NewDbSessionService(db),
		},
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", s.listenAddr)
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
			log.Printf("recovered from panic: %v\n", err)
			debug.PrintStack()

			if err := conn.Close(); err != nil {
				log.Printf("error closing after recover: %v\n", err)
			}
		}
	}()

	log.Printf("client connected from %v\n", conn.RemoteAddr().String())

	client, err := realmd.NewClient(conn)
	if err != nil {
		log.Printf("error setting up client: %v\n", err)
		conn.Close()
		return
	}

	// In realmd the server initiates the auth challenge
	if err := auth.SendChallenge(client); err != nil {
		log.Printf("error sending auth challenge: %v\n", err)
		conn.Close()
		return
	}

	var header *realmd.ClientHeader
	chunk := make([]byte, 4096)
	readBuf := bytes.Buffer{}
	packetBuf := bytes.Buffer{}
	headerBuf := make([]byte, realmd.ClientHeaderSize)

	for {
		n, err := conn.Read(chunk)
		if err != nil && err != io.EOF {
			log.Printf("error reading from client: %v\n", err)
			return
		}

		readBuf.Write(chunk[:n])

		// Process the read buffer while it has at least enough data for a packet header
		for readBuf.Len() >= realmd.ClientHeaderSize {

			// Ready to process a new packet, read the header
			if header == nil {
				readBuf.Read(headerBuf)

				h, err := client.ParseHeader(headerBuf)
				if err != nil {
					log.Printf("failed to parse header: %v\n", err)
					conn.Close()
					return
				}
				header = h
			}

			// The buffer contains a partial packet, go back to reading
			if bufLen := readBuf.Len(); bufLen < int(header.Size) {
				break
			}

			io.CopyN(&packetBuf, &readBuf, int64(header.Size))

			if err := s.handlePacket(client, header, packetBuf.Bytes()); err != nil {
				log.Printf("error handling packet: %v\n", err)
				conn.Close()
				return
			}

			header = nil
			packetBuf.Reset()
		}

		// TODO: use cancel context?
		if err == io.EOF {
			log.Println("client disconnected (closed by client)")
			return
		}
	}
}

func (s *Server) handlePacket(c *realmd.Client, header *realmd.ClientHeader, data []byte) error {
	var opName string

	if header.Opcode.IsAClientOpcode() {
		opName = header.Opcode.String()
	} else {
		opName = fmt.Sprintf("UNKNOWN(0x%x)", uint32(header.Opcode))
	}

	log.Printf("[IN]    %s size=%d", opName, header.Size)

	switch header.Opcode {
	case realmd.OpClientPing:
		return session.PingHandler(c, data)

	case realmd.OpClientAuthSession:
		return auth.ProofHandler(s.services, c, data)

	case realmd.OpClientCharList:
		return char.ListHandler(s.services, c)

	case realmd.OpClientRealmSplit:
		return realm.SplitInfoHandler(c, data)

	case realmd.OpClientCharCreate:
		return char.CreateHandler(s.services, c, data)

	case realmd.OpClientCharDelete:
		return char.DeleteHandler(s.services, c, data)

	case realmd.OpClientReadyForAccountDataTimes:
		return account.StorageTimesHandler(s.services, c)

	case realmd.OpClientPlayerLogin:
		return session.LoginHandler(s.services, c, data)

	case realmd.OpClientLogout:
		return session.LogoutHandler(c)

	case realmd.OpClientLogoutCancel:
		return session.LogoutCancelHandler(c)

	default:
		return nil
	}
}
