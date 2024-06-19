package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"runtime/debug"

	"github.com/jmoiron/sqlx"
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
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
			Accounts: model.NewDbAccountService(db),
			Chars:    model.NewDbCharacterService(db),
			Realms:   model.NewDbRealmService(db),
			Sessions: model.NewDbSessionService(db),
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
	readBuf := bytes.Buffer{}
	headerBuf := make([]byte, realmd.ClientHeaderSize)
	packetBuf := bytes.Buffer{}

	for {
		n, err := io.Copy(&readBuf, conn)
		if err != nil {
			log.Printf("error reading from client: %v\n", err)
			conn.Close()
			return
		} else if n == 0 {
			log.Println("client disconnected (closed by client)")
			return
		}

		// Process the read buffer while it has at least enough data for a packet header
		for readBuf.Len() > realmd.ClientHeaderSize {

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
			if readBuf.Len() < int(header.Size) {
				break
			}

			io.CopyN(&packetBuf, &readBuf, int64(header.Size))

			data := packetBuf.Bytes()

			log.Printf("%d: %x\n", len(data), data)

			if err := s.handlePacket(client, header, data); err != nil {
				log.Printf("error handling packet: %v\n", err)
				conn.Close()
				return
			}

			header = nil
			packetBuf.Reset()
		}
	}
}

func (s *Server) handlePacket(c *realmd.Client, header *realmd.ClientHeader, data []byte) error {
	log.Printf("RECV  op=0x%-4x dsize=%d hsize=%d  \n", header.Opcode, len(data), header.Size)

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
		return session.DataTimesHandler(c)

	case realmd.OpClientPlayerLogin:
		return session.LoginHandler(s.services, c, data)

	case realmd.OpClientLogout:
		return session.LogoutHandler(c)

	case realmd.OpClientLogoutCancel:
		return session.LogoutCancelHandler(c)

	default:
		log.Printf("unknown opcode: 0x%x\n", header.Opcode)
		return nil
	}
}
