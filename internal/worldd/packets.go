package worldd

type Header struct {
	Size   uint16
	Opcode uint32
}

// https://gtker.com/wow_messages/docs/cmsg_auth_session.html#client-version-335
type AuthSessionPacket struct {
	ClientBuild     uint32
	LoginServerId   uint32
	Username        string
	LoginServerType uint32
	ClientSeed      [4]byte
	RegionId        uint32
	BattlegroundId  uint32
	RealmId         uint32
	DOSResponse     uint64
	ClientProof     [20]byte
	AddonInfo       []byte
}

// https://gtker.com/wow_messages/docs/cmsg_ping.html#client-version-19-client-version-110-client-version-111-client-version-112-client-version-2-client-version-3
type PingPacket struct {
	SequenceId    uint32
	RoundTripTime uint32 // zero if server hasn't responded?
}

// https://gtker.com/wow_messages/docs/cmsg_realm_split.html
type RealmSplitPacket struct {
	RealmId uint32
}

// https://gtker.com/wow_messages/docs/cmsg_char_create.html#client-version-32-client-version-33
type CharCreatePacket struct {
	// Name string
	Race       Race
	Class      Class
	Gender     Gender
	SkinColor  byte
	Face       byte
	HairStyle  byte
	HairColor  byte
	FacialHair byte
	OutfitId   byte
}

// https://gtker.com/wow_messages/docs/cmsg_char_delete.html
type CharDeletePacket struct {
	CharacterId uint64
}
