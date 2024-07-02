package server

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
	"github.com/kangaroux/gomaggus/realmd/handler/account"
	"github.com/kangaroux/gomaggus/realmd/handler/auth"
	"github.com/kangaroux/gomaggus/realmd/handler/char"
	"github.com/kangaroux/gomaggus/realmd/handler/player"
	"github.com/kangaroux/gomaggus/realmd/handler/realm"
	"github.com/kangaroux/gomaggus/realmd/handler/session"
	"github.com/kangaroux/gomaggus/realmd/handler/world"
	"github.com/phuslu/log"
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
			Accounts:         model.NewDbAccountService(db),
			AccountStorage:   model.NewDbAccountStorageService(db),
			CharacterStorage: model.NewDbCharacterStorageService(db),
			Characters:       model.NewDbCharacterService(db),
			Realms:           model.NewDbRealmService(db),
			Sessions:         model.NewDbSessionService(db),
		},
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("error setting up tcp server")
	}

	defer listener.Close()
	log.Info().Str("listen", listener.Addr().String()).Msg("realmd start")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error().Err(err).Msg("error accepting client")
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Warn().Stack().Any("err", err).Msg("recovered from panic")
		}
		conn.Close()
	}()

	client, err := realmd.NewClient(conn)
	if err != nil {
		log.Error().Err(err).Msg("error setting up client")
		return
	}

	ip := strings.Split(client.Conn.RemoteAddr().String(), ":")[0]

	// Create a logger for the client that includes the client's ID/IP
	*client.Log = log.DefaultLogger
	client.Log.Context = log.NewContext(nil).
		Int64("cid", client.ID).
		Str("ip", ip).
		Value()

	client.Log.Info().Msg("client connected")

	// In realmd the server initiates the auth challenge
	if err := auth.SendChallenge(client); err != nil {
		client.Log.Error().Err(err).Msg("error sending auth challenge")
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
			client.Log.Error().Err(err).Msg("error reading socket")
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
					client.Log.Error().Err(err).Msg("error parsing packet header")
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
				if err, ok := err.(*realmd.ErrKickClient); ok {
					client.Log.Warn().Str("reason", err.Reason).Msg("kicking client")
					return
				}

				client.Log.Error().Err(err).Msg("error handling packet")
				return
			}

			header = nil
			packetBuf.Reset()
		}

		if err == io.EOF {
			client.Log.Info().Msg("client disconnected")
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

	c.Log.Debug().Str("op", opName).Uint16("size", header.Size).Msg("packet recv")
	c.Log.Trace().
		Func(func(e *log.Entry) { // Skip encoding unless it's actually needed
			e.Str("data", hex.EncodeToString(data))
		}).
		Msg("recv data")

	switch header.Opcode {
	// Ignored
	case realmd.OpClientLogoutForce:
		c.Log.Debug().Str("op", opName).Msg("packet ignored")
		return nil

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
		h := &account.StorageTimesHandler{
			Client:  c,
			Service: s.services,
		}
		return h.Handle()

	case realmd.OpClientPlayerLogin:
		return session.LoginHandler(s.services, c, data)

	case realmd.OpClientLogoutRequest:
		return session.LogoutHandler(c)

	case realmd.OpClientLogoutCancel:
		return session.LogoutCancelHandler(c)

	case realmd.OpClientPutStorage:
		h := &account.StoragePutHandler{
			Client:  c,
			Service: s.services,
		}
		return h.Handle(data)

	case realmd.OpClientGetStorage:
		h := &account.StorageGetHandler{
			Client:  c,
			Service: s.services,
		}
		return h.Handle(data)

	case realmd.OpClientGetUnixTime:
		return world.UnixTimeHandler(c)

	case realmd.OpClientStandStateChange:
		return player.StandStateHandler(c, data)

	case realmd.OpClientGetPlayerName:
		return player.NameHandler(*s.services, c, data)

	default:
		return nil
	}
}
