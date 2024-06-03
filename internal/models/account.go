package models

import (
	"database/sql"
	"time"
)

type Account struct {
	Id             uint32
	CreatedAt      time.Time    `db:"created_at"`
	LastLogin      sql.NullTime `db:"last_login"`
	Username       string
	SrpVerifierHex string `db:"srp_verifier"`
	SrpSaltHex     string `db:"srp_salt"`
	Email          string
	RealmId        uint32 `db:"realm_id"`
}
