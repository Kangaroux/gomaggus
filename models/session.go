package models

import (
	"database/sql"
	"encoding/hex"
)

type Session struct {
	AccountId      uint32       `db:"account_id"`
	SessionKeyHex  string       `db:"session_key"`
	Connected      uint8        // TODO: add types
	ConnectedAt    sql.NullTime `db:"connected_at"`
	DisconnectedAt sql.NullTime `db:"disconnected_at"`

	sessionKey []byte
}

func (s *Session) Decode() error {
	var err error
	if s.sessionKey, err = hex.DecodeString(s.SessionKeyHex); err != nil {
		return err
	}
	return nil
}

func (s *Session) SessionKey() []byte {
	if s.sessionKey == nil {
		panic("Decode must be called before accessing SessionKey")
	}
	return s.sessionKey
}
