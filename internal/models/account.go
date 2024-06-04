package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/kangaroux/gomaggus/srp"
)

type Account struct {
	Id        uint32
	CreatedAt time.Time    `db:"created_at"`
	LastLogin sql.NullTime `db:"last_login"`

	Username       string
	Email          string
	RealmId        uint32 `db:"realm_id"`
	SrpSaltHex     string `db:"srp_salt"`
	SrpVerifierHex string `db:"srp_verifier"`
}

func (acc *Account) SetUsernamePassword(username, password string) error {
	var err error
	var salt []byte

	if acc.SrpSaltHex == "" {
		salt = make([]byte, 32)
		if _, err = rand.Read(salt); err != nil {
			return err
		}
		acc.SrpSaltHex = hex.EncodeToString(salt)
	} else {
		if salt, err = hex.DecodeString(acc.SrpSaltHex); err != nil {
			return err
		}
	}

	acc.Username = username
	acc.SrpVerifierHex = hex.EncodeToString(srp.CalculateVerifier(username, password, salt))

	return nil
}
