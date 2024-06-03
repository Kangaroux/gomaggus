package realmd

const (
	OP_LOGIN_CHALLENGE     byte = 0
	OP_LOGIN_PROOF         byte = 1
	OP_RECONNECT_CHALLENGE byte = 2
	OP_RECONNECT_PROOF     byte = 3
	OP_REALM_LIST          byte = 16
)

const (
	WOW_SUCCESS              byte = 0
	WOW_FAIL_UNKNOWN_ACCOUNT byte = 4
)

type RealmFlag uint8

const (
	REALMFLAG_NONE          RealmFlag = 0
	REALMFLAG_INVALID       RealmFlag = 1 // Realm is greyed out and can't be selected
	REALMFLAG_OFFLINE       RealmFlag = 2 // Population: "Offline" and can't be selected
	REALMFLAG_SPECIFY_BUILD RealmFlag = 4 // Includes version in realm name
	REALMFLAG_UNKNOWN1      RealmFlag = 8
	REALMFLAG_UNKNOWN2      RealmFlag = 16
	REALMFLAG_NEW_PLAYERS   RealmFlag = 32  // Population: "New Players" in blue text
	REALMFLAG_NEW_SERVER    RealmFlag = 64  // Population: "New" in green text
	REALMFLAG_FULL          RealmFlag = 128 // Population: "Full" in red text
)
