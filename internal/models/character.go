package models

import (
	"database/sql"
	"time"
)

type Character struct {
	Id        uint32
	CreatedAt time.Time    `db:"created_at"`
	LastLogin sql.NullTime `db:"last_login"`

	Name       string
	AccountId  uint32 `db:"account_id"`
	RealmId    uint32 `db:"realm_id"`
	Race       byte
	Class      byte
	Gender     byte
	SkinColor  byte `db:"skin_color"`
	Face       byte
	HairStyle  byte `db:"hair_style"`
	HairColor  byte `db:"hair_color"`
	FacialHair byte `db:"facial_hair"`
	OutfitId   byte `db:"outfit_id"`
}
