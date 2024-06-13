package model

import (
	"time"
)

type Realm struct {
	Id        uint32
	CreatedAt time.Time `db:"created_at"`

	Name   string
	Type   RealmType
	Host   string
	Region RealmRegion
}

type RealmType = uint8

const (
	REALMTYPE_PVE   RealmType = 0
	REALMTYPE_PVP   RealmType = 1
	REALMTYPE_RP    RealmType = 6
	REALMTYPE_RPPVP RealmType = 8
)

type RealmRegion = uint8

const (
	REALMREGION_DEV           RealmRegion = 1
	REALMREGION_US            RealmRegion = 2
	REALMREGION_OCEANIC       RealmRegion = 3
	REALMREGION_LATIN_AMERICA RealmRegion = 4
	REALMREGION_TOURNAMENT    RealmRegion = 5
	REALMREGION_KOREA         RealmRegion = 6
	REALMREGION_TOURNAMENT2   RealmRegion = 7
	REALMREGION_ENGLISH       RealmRegion = 8
	REALMREGION_GERMAN        RealmRegion = 9
	REALMREGION_FRENCH        RealmRegion = 10
	REALMREGION_SPANISH       RealmRegion = 11
	REALMREGION_RUSSIAN       RealmRegion = 12
	REALMREGION_TOURNAMENT3   RealmRegion = 13
	REALMREGION_TAIWAN        RealmRegion = 14
	REALMREGION_TOURNAMENT4   RealmRegion = 15
	REALMREGION_CHINA         RealmRegion = 16
	REALMREGION_CN1           RealmRegion = 17
	REALMREGION_CN2           RealmRegion = 18
	REALMREGION_CN3           RealmRegion = 19
	REALMREGION_CN4           RealmRegion = 20
	REALMREGION_CN5           RealmRegion = 21
	REALMREGION_CN6           RealmRegion = 22
	REALMREGION_CN7           RealmRegion = 23
	REALMREGION_CN8           RealmRegion = 24
	REALMREGION_TOURNAMENT5   RealmRegion = 25
	REALMREGION_TEST          RealmRegion = 26
	REALMREGION_TOURNAMENT6   RealmRegion = 27
	REALMREGION_QA            RealmRegion = 28
	REALMREGION_CN9           RealmRegion = 29
	REALMREGION_TEST2         RealmRegion = 30
	REALMREGION_CN10          RealmRegion = 31
	REALMREGION_CTC           RealmRegion = 32
	REALMREGION_CNC           RealmRegion = 33
	REALMREGION_CN1_4         RealmRegion = 34
	REALMREGION_CN2_6_9       RealmRegion = 35
	REALMREGION_CN3_7         RealmRegion = 36
	REALMREGION_CN5_8         RealmRegion = 37
)
