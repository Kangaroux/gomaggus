package realmd

type RealmVersion struct {
	Major uint8
	Minor uint8
	Patch uint8
	Build uint16
}

// https://gtker.com/wow_messages/docs/realm.html#protocol-version-8
type Realm struct {
	Type   RealmType
	Locked bool
	Flags  RealmFlag
	Name   string // C-style NUL terminated, e.g. "Test Realm\x00"
	Host   string // C-style NUL terminated, e.g. "localhost:8085\x00"

	// A percentage of how full the server is with active sessions. Mangos has the upper limit of this
	// value as 2.0 for some reason. The game client only seems to interpret this value on an absolute
	// scale if there is only one realm. It seems like when there are multiple realms, it compares the
	// pop relatively, i.e. whatever realm has the highest pop is now the upper limit. Suffice to say,
	// it's not important whether this value is accurate.
	Population      float32
	NumCharsOnRealm uint8 // Number of characters for the logged in account
	Region          RealmRegion
	Id              uint8
	Version         RealmVersion // included only if REALMFLAG_SPECIFY_BUILD flag is set
}
