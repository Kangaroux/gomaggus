package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/kangaroux/gomaggus/internal/srp"
)

type Account struct {
	Id        uint32
	CreatedAt time.Time    `db:"created_at"`
	LastLogin sql.NullTime `db:"last_login"`

	Username       string
	Email          string
	SrpSaltHex     string `db:"srp_salt"`
	SrpVerifierHex string `db:"srp_verifier"`

	srpSalt     []byte
	srpVerifier []byte
}

func (acc *Account) SetUsernamePassword(username, password string) error {
	salt := make([]byte, srp.SaltSize)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	acc.Username = username
	acc.SrpSaltHex = hex.EncodeToString(salt)
	acc.SrpVerifierHex = hex.EncodeToString(srp.CalculateVerifier(username, password, salt))

	return nil
}

func (acc *Account) DecodeSrp() error {
	var err error
	if acc.srpSalt, err = hex.DecodeString(acc.SrpSaltHex); err != nil {
		return err
	}
	if acc.srpVerifier, err = hex.DecodeString(acc.SrpVerifierHex); err != nil {
		return err
	}
	return nil
}

func (acc *Account) Salt() []byte {
	if acc.srpSalt == nil {
		panic("DecodeSrp must be called before accessing Salt")
	}
	return acc.srpSalt
}

func (acc *Account) Verifier() []byte {
	if acc.srpVerifier == nil {
		panic("DecodeSrp must be called before accessing Verifier")
	}
	return acc.srpVerifier
}
