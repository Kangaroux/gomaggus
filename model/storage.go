package model

import (
	"time"
)

type AccountStorageType uint8

const (
	AccountData     AccountStorageType = 0
	AccountKeybinds AccountStorageType = 2
	AccountMacros   AccountStorageType = 4
)

type CharacterStorageType uint8

const (
	CharacterConfig   CharacterStorageType = 1
	CharacterKeybinds CharacterStorageType = 3
	CharacterMacros   CharacterStorageType = 5
	CharacterLayout   CharacterStorageType = 6
	CharacterChat     CharacterStorageType = 7
)

const (
	AccountStorageCount   = 3
	CharacterStorageCount = 5
	AllAccountStorage     = (1 << AccountData) | (1 << AccountKeybinds) | (1 << AccountMacros)
	AllCharacterStorage   = (1 << CharacterConfig) | (1 << CharacterKeybinds) | (1 << CharacterMacros) | (1 << CharacterLayout) | (1 << CharacterChat)
)

// AccountStorage stores compressed data from the client that is linked to their account.
type AccountStorage struct {
	AccountId uint32    `db:"account_id"`
	UpdatedAt time.Time `db:"updated_at"`
	Type      AccountStorageType
	Data      []byte
}

// CharacterStorage stores compressed data from the client that is linked to a specific character.
type CharacterStorage struct {
	CharacterId uint32    `db:"character_id"`
	UpdatedAt   time.Time `db:"updated_at"`
	Type        CharacterStorageType
	Data        []byte
}
