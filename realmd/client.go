package realmd

import (
	"net"

	"github.com/kangaroux/gomaggus/model"
)

type Client struct {
	conn          net.Conn
	serverSeed    [4]byte
	authenticated bool
	crypto        *WrathHeaderCrypto

	account *model.Account
	realm   *model.Realm
	session *model.Session
}
