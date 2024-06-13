package authd

import "github.com/kangaroux/gomaggus/internal/srp"

// https://gtker.com/wow_messages/docs/cmd_auth_logon_challenge_server.html#protocol-version-8
// FIELD ORDER MATTERS, DO NOT REORDER
type ServerLoginChallenge struct {
	Opcode          byte
	ProtocolVersion byte
	ErrorCode       byte
	PublicKey       [srp.KeySize]byte
	GeneratorSize   byte
	Generator       byte
	LargePrimeSize  byte
	LargePrime      [srp.LargePrimeSize]byte
	Salt            [srp.SaltSize]byte
	CrcHash         [16]byte

	// Using any flags would require additional fields but this is set to zero for now
	SecurityFlags byte
}

// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_server.html#protocol-version-8
type ServerLoginProofFail struct {
	Opcode    byte
	ErrorCode byte
	_         [2]byte // padding
}

type ServerLoginProofSuccess struct {
	Opcode           byte
	ErrorCode        byte
	Proof            [20]byte
	AccountFlags     uint32
	HardwareSurveyId uint32
	_                [2]byte // padding
}

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

// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_challenge_server.html#protocol-version-8
type ServerReconnectChallenge struct {
	Opcode        byte
	ErrorCode     byte
	ReconnectData [ReconnectDataLen]byte
	ChecksumSalt  [16]byte
}

// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_proof_server.html#protocol-version-8
type ServerReconnectProof struct {
	Opcode    byte
	ErrorCode byte
	_         [2]byte // padding
}
