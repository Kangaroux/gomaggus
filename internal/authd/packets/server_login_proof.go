package packets

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
