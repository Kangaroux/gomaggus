package server

import (
	"fmt"
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

	buf := make([]byte, 4096)

	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			log.Println("client disconnected (closed by client)")
			return
		} else if err != nil {
			log.Printf("error reading from client: %v\n", err)
			conn.Close()
			return
		}

		if err := s.handlePacket(client, buf[:n]); err != nil {
			log.Printf("error handling packet: %v\n", err)
			conn.Close()
			return
		}
	}
}

func (s *Server) handlePacket(c *realmd.Client, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("handlePacket: packet is empty")
	}

	header, err := c.ParseHeader(data)
	if err != nil {
		return err
	}

	log.Printf("RECV  op=0x%-4x dsize=%d hsize=%d  \n", header.Opcode, len(data), header.Size)

	packet := &realmd.ClientPacket{
		Header:  header,
		Payload: data[realmd.ClientHeaderSize:],
	}

	switch header.Opcode {
	case realmd.OpClientPing:
		return session.PingHandler(c, packet)

	case realmd.OpClientAuthSession:
		return auth.ProofHandler(s.services, c, packet)

	case realmd.OpClientCharList:
		return char.ListHandler(s.services, c)

	case realmd.OpClientRealmSplit:
		return realm.SplitInfoHandler(c, packet)

	case realmd.OpClientCharCreate:
		return char.CreateHandler(s.services, c, packet)

	case realmd.OpClientCharDelete:
		return char.DeleteHandler(s.services, c, packet)

	case realmd.OpClientReadyForAccountDataTimes:
		return session.DataTimesHandler(c)

	// FIXME: Sometimes the login packet header is garbled, and it consistently happens when receiving
	// the login packet. Once this happens, all future incoming headers are garbled. I'm not sure which
	// side is causing it, but the RC4 stream is getting out of sync and makes the connection useless.
	//
	// I suspect the reason it only happens sometimes is because the cipher is seeded using the session key.
	// This could be verified by replaying the packets with the same session key and seeing if it
	// consistently breaks. If it's reproducible based on the session key, there may be a pattern that
	// can be identified by generating random keys and seeing which break. For example, a leading or
	// trailing zero in the key.
	//
	// 2024/06/18 18:59:05 RECV  op=0x38c   dsize=92     hsize=8
	// 2024/06/18 18:59:05 SENT  op=0x38b   size=17
	// 2024/06/18 18:59:05 sent realm split
	// 2024/06/18 18:59:09 RECV  op=0x6e5b92f8  dsize=14     hsize=1396
	// 2024/06/18 18:59:09 unknown opcode: 0x6e5b92f8
	case realmd.OpClientPlayerLogin:
		return session.LoginHandler(s.services, c, packet)

	case realmd.OpClientLogout:
		return session.LogoutHandler(c)

	case realmd.OpClientLogoutCancel:
		return session.LogoutCancelHandler(c)

	default:
		log.Printf("unknown opcode: 0x%x\n", header.Opcode)
		return nil
	}
}
