package packets

import (
	"bytes"
	"encoding/binary"
)

// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_proof_client.html
// FIELD ORDER MATTERS, DO NOT REORDER
type ClientReconnectProof struct {
	Opcode         byte
	ProofData      [16]byte
	ClientProof    [20]byte
	ClientChecksum [20]byte
	KeyCount       byte
}

func (p *ClientReconnectProof) Read(data []byte) error {
	reader := bytes.NewReader(data)
	return binary.Read(reader, binary.LittleEndian, p)
}
