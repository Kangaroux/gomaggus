package model

import (
	"database/sql"
	"fmt"
	"time"
)

type Character struct {
	Id        uint32
	CreatedAt time.Time    `db:"created_at"`
	LastLogin sql.NullTime `db:"last_login"`

	Name      string
	AccountId uint32 `db:"account_id"`
	RealmId   uint32 `db:"realm_id"`
	Race      Race
	Class     Class
	Gender    Gender
	SkinColor byte `db:"skin_color"`
	Face      byte
	HairStyle byte `db:"hair_style"`
	HairColor byte `db:"hair_color"`

	// Facial hair, piercings, etc.
	ExtraCosmetic byte `db:"extra_cosmetic"`
	OutfitId      byte `db:"outfit_id"`
}

func (c *Character) String() string {
	return fmt.Sprintf("Character(%d, %s)", c.Id, c.Name)
}

type Race byte

const (
	RaceHuman             Race = 1
	RaceOrc               Race = 2
	RaceDwarf             Race = 3
	RaceNightElf          Race = 4
	RaceUndead            Race = 5
	RaceTauren            Race = 6
	RaceGnome             Race = 7
	RaceTroll             Race = 8
	RaceGoblin            Race = 9
	RaceBloodElf          Race = 10
	RaceDraenei           Race = 11
	RaceFelOrc            Race = 12
	RaceNaga              Race = 13
	RaceBroken            Race = 14
	RaceSkeleton          Race = 15
	RaceVrykul            Race = 16
	RaceTuskarr           Race = 17
	RaceForestTroll       Race = 18
	RaceTaunka            Race = 19
	RaceNorthrendSkeleton Race = 20
	RaceIceTroll          Race = 21
)

type Class byte

const (
	ClassWarrior     Class = 1
	ClassPaladin     Class = 2
	ClassHunter      Class = 3
	ClassRogue       Class = 4
	ClassPriest      Class = 5
	ClassDeathKnight Class = 6
	ClassShaman      Class = 7
	ClassMage        Class = 8
	ClassWarlock     Class = 9
	ClassDruid       Class = 11
)

type Gender byte

const (
	GenderMale   Gender = 0
	GenderFemale Gender = 1
	GenderNone   Gender = 2 // used by pets?
)
