package server

import (
	"bytes"
	"fmt"
	"io"
	golog "log"
	"net"
	"runtime/debug"
	"strings"

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
	golog.SetFlags(golog.Lmicroseconds)
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
		golog.Fatal(err)
	}

	defer listener.Close()
	golog.Println("listening on", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			golog.Fatal(err)
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			golog.Println("recovered from panic:", err)
			debug.PrintStack()

			if err := conn.Close(); err != nil {
				golog.Println("error closing after recover:", err)
			}
		}
	}()

	golog.Println("client connected from", conn.RemoteAddr())

	client, err := realmd.NewClient(conn)
	if err != nil {
		golog.Println("error setting up client:", err)
		conn.Close()
		return
	}

	// In realmd the server initiates the auth challenge
	if err := auth.SendChallenge(client); err != nil {
		golog.Println("error sending auth challenge:", err)
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
			golog.Println("error reading from client:", err)
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
					golog.Println("failed to parse header:", err)
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
				golog.Println("error handling packet:", err)
				conn.Close()
				return
			}

			header = nil
			packetBuf.Reset()
		}

		// TODO: use cancel context?
		if err == io.EOF {
			var msg strings.Builder
			msg.WriteString("Disconnected ")

			if client.Account != nil {
				msg.WriteString(client.Account.String())
			} else {
				msg.WriteString("<unknown>")
			}

			if client.Realm != nil {
				msg.WriteString(" from ")
				msg.WriteString(client.Realm.String())
			}

			golog.Println(msg)
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

	golog.Printf("[IN]    %s size=%d", opName, header.Size)

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
