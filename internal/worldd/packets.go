package worldd

type Header struct {
	Size   uint16
	Opcode uint32
}

type AuthSessionPacket struct {
	ClientBuild     uint32
	LoginServerId   uint32
	Username        string
	LoginServerType uint32
	ClientSeed      uint32
	RegionId        uint32
	BattlegroundId  uint32
	RealmId         uint32
	DOSResponse     uint64
	ClientProof     [20]byte
	AddonInfo       []byte
}
