package packets

import (
	"bytes"
	"encoding/binary"
)

// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_client.html#protocol-version-8
// FIELD ORDER MATTERS, DO NOT REORDER
type ClientLoginProof struct {
	Opcode           byte
	ClientPublicKey  [32]byte
	ClientProof      [20]byte
	CRCHash          [20]byte
	NumTelemetryKeys uint8
}

func (p *ClientLoginProof) Read(data []byte) error {
	reader := bytes.NewReader(data)
	return binary.Read(reader, binary.LittleEndian, p)
}
