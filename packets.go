package main

type LoginChallengePacket struct {
	Opcode            byte // 0x0
	Error             byte // unused?
	Size              uint16
	GameName          [4]byte
	Version           [3]byte
	Build             uint16
	OSArch            [4]byte
	OS                [4]byte
	Locale            [4]byte
	TimezoneBias      uint32
	IP                [4]byte
	AccountNameLength uint8

	// The account name is a variable size and needs to be read manually
	// AccountName    []byte
}
