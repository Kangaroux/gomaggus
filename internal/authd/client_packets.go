package authd

import (
	"bytes"
	"encoding/binary"
	"strings"
)

// https://gtker.com/wow_messages/docs/cmd_auth_logon_challenge_client.html
// FIELD ORDER MATTERS, DO NOT REORDER
type loginChallengeFixed struct {
	Opcode         byte
	Error          byte
	Size           uint16
	GameName       [4]byte
	Version        [3]byte
	Build          uint16
	OSArch         [4]byte
	OS             [4]byte
	Locale         [4]byte
	TimezoneBias   uint32
	IP             [4]byte
	UsernameLength uint8
}

type ClientLoginChallenge struct {
	loginChallengeFixed
	Username string
}

func (p *ClientLoginChallenge) Read(data []byte) error {
	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.LittleEndian, &p.loginChallengeFixed); err != nil {
		return err
	}

	usernameBytes := make([]byte, p.UsernameLength)
	if _, err := reader.Read(usernameBytes); err != nil {
		return err
	}

	if reader.Len() != 0 {
		return &ErrPacketUnreadBytes{What: "LoginChallengePacket", Count: reader.Len()}
	}

	p.Username = strings.ToUpper(string(usernameBytes))
	return nil
}

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
