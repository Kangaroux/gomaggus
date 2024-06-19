package model

import (
	"database/sql"
	"time"
)

type StorageType uint8

const (
	AccountData     StorageType = 0
	AccountKeybinds StorageType = 2
	AccountMacros   StorageType = 4

	CharacterConfig   StorageType = 1
	CharacterKeybinds StorageType = 3
	CharacterMacros   StorageType = 5
	CharacterLayout   StorageType = 6
	CharacterChat     StorageType = 7
)

// AccountStorage stores compressed client data.
type AccountStorage struct {
	AccountId uint32
	UpdatedAt time.Time
	Type      StorageType
	Data      sql.RawBytes
}
