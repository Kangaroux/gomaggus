package models

import "database/sql"

type Session struct {
	Id             uint32
	Connected      uint8        // TODO: add types
	SessionKeyHex  string       `db:"session_key"`
	ConnectedAt    sql.NullTime `db:"connected_at"`
	DisconnectedAt sql.NullTime `db:"disconnected_at"`
}
