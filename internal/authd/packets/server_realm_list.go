package packets

type ServerRealm struct {
	Type          byte
	Locked        bool
	Flags         byte
	Name          string
	Host          string
	Population    float32
	NumCharacters byte
	Region        byte
	Id            byte
}

// https://gtker.com/wow_messages/docs/cmd_realm_list_server.html#protocol-version-8
type ServerRealmList struct {
	Opcode byte
	Size   uint16
	_      [4]byte // header padding
	Realms []ServerRealm
	_      [2]byte // footer padding
}
