package packets

// https://gtker.com/wow_messages/docs/cmd_realm_list_server.html#protocol-version-8
type ServerRealm struct {
	Type          byte
	Locked        bool
	Flags         byte
	Name          string `binary:"zstring"`
	Host          string `binary:"zstring"`
	Population    float32
	NumCharacters byte
	Region        byte
	Id            byte
}

type ServerRealmListBody struct {
	_         [4]byte // header padding
	NumRealms uint16
	Realms    []ServerRealm `binary:"[NumRealms]Any"`
	_         [2]byte       // footer padding
}

type ServerRealmListHeader struct {
	Opcode byte
	Size   uint16
}
