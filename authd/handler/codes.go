package handler

type Opcode byte

const (
	OpLoginChallenge     Opcode = 0x0
	OpLoginProof         Opcode = 0x1
	OpReconnectChallenge Opcode = 0x2
	OpReconnectProof     Opcode = 0x3
	OpRealmList          Opcode = 0x10
)

type RespCode byte

const (
	CodeSuccess            RespCode = 0x0
	CodeFailUnknownAccount RespCode = 0x4
)
