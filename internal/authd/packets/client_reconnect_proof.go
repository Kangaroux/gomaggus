package packets

import (
	"bytes"
	"encoding/binary"
)

// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_proof_client.html
type ClientReconnectProof struct {
	Opcode         byte // 0x3
	ProofData      [16]byte
	ClientProof    [20]byte
	ClientChecksum [20]byte // unused
	KeyCount       byte     // unused
}

func (p *ClientReconnectProof) Read(data []byte) error {
	reader := bytes.NewReader(data)
	return binary.Read(reader, binary.LittleEndian, p)
}
