package model

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/kangaroux/gomaggus/srp"
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

func (a *Account) String() string {
	return fmt.Sprintf("Account(\"%s\" id=%d)", a.Username, a.Id)
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
	salt, err := hex.DecodeString(acc.SrpSaltHex)
	if err != nil {
		return err
	}

	verifier, err := hex.DecodeString(acc.SrpVerifierHex)
	if err != nil {
		return err
	}

	acc.srpSalt = salt
	acc.srpVerifier = verifier

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
