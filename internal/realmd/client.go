package realmd

import (
	"net"

	"github.com/kangaroux/gomaggus/internal/models"
)

type Client struct {
	conn          net.Conn
	serverSeed    [4]byte
	authenticated bool
	crypto        *WrathHeaderCrypto

	account *models.Account
	realm   *models.Realm
	session *models.Session
}
