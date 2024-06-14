package authd

type Opcode byte

const (
	OpcodeLoginChallenge     Opcode = 0x0
	OpcodeLoginProof         Opcode = 0x1
	OpcodeReconnectChallenge Opcode = 0x2
	OpcodeReconnectProof     Opcode = 0x3
	OpcodeRealmList          Opcode = 0x10
)

type RespCode byte

const (
	Success        RespCode = 0x0
	UnknownAccount RespCode = 0x4
)
