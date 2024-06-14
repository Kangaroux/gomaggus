package realmd

import (
	"net"

	"github.com/kangaroux/gomaggus/model"
)

type Client struct {
	Conn          net.Conn
	ServerSeed    [4]byte
	Authenticated bool
	Crypto        *HeaderCrypto

	Account *model.Account
	Realm   *model.Realm
	Session *model.Session
}
