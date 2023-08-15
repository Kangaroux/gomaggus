package main

type LoginChallengePacket struct {
	Opcode       byte
	Err          byte
	Size         uint16
	GameName     [4]byte
	Version      [3]byte
	Build        uint16
	OSArch       [4]byte
	OS           [4]byte
	Locale       [4]byte
	TimezoneBias uint32
	Ip           uint32
	ILen         uint8
	I            byte
}

func reverseBytes(data []byte, n int) {
	for i := 0; i < n/2; i++ {
		data[i], data[n-i-1] = data[n-i-1], data[i]
	}
}
