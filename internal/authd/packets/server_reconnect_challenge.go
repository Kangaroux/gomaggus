package packets

const (
	ReconnectDataLen = 16
)

// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_challenge_server.html#protocol-version-8
type ServerReconnectChallenge struct {
	Opcode        byte
	ErrorCode     byte
	ReconnectData [ReconnectDataLen]byte
	ChecksumSalt  [16]byte
}
