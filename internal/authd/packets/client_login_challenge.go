package packets

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
