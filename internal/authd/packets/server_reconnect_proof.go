package packets

// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_proof_server.html#protocol-version-8
type ServerReconnectProof struct {
	Opcode    byte
	ErrorCode byte
	_         [2]byte // padding
}
