package models

import "database/sql"

type Session struct {
	Id uint32

	AccountId      uint32       `db:"account_id"`
	SessionKeyHex  string       `db:"session_key"`
	Connected      uint8        // TODO: add types
	ConnectedAt    sql.NullTime `db:"connected_at"`
	DisconnectedAt sql.NullTime `db:"disconnected_at"`
}
