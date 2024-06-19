package model

import (
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

type StorageMask byte

const (
	StorageMaskAccount   StorageMask = (1 << AccountData) | (1 << AccountKeybinds) | (1 << AccountMacros)
	StorageMaskCharacter StorageMask = (1 << CharacterConfig) | (1 << CharacterKeybinds) | (1 << CharacterMacros) | (1 << CharacterLayout) | (1 << CharacterChat)
	StorageMaskAll       StorageMask = StorageMaskAccount | StorageMaskCharacter
)

// AccountStorage stores compressed client data.
type AccountStorage struct {
	AccountId uint32    `db:"account_id"`
	UpdatedAt time.Time `db:"updated_at"`
	Type      StorageType
	Data      []byte
}
