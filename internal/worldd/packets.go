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
