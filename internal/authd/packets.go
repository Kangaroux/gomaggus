package authd

// https://gtker.com/wow_messages/docs/cmd_auth_logon_challenge_client.html
type LoginChallengePacket struct {
	Opcode         byte // 0x0 or 0x2 if reconnecting
	Error          byte // unused?
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

	// The username is a variable size and needs to be read manually
	// Username string
}

// https://gtker.com/wow_messages/docs/cmd_auth_logon_proof_client.html#protocol-version-8
type LoginProofPacket struct {
	Opcode           byte // 0x1
	ClientPublicKey  [32]byte
	ClientProof      [20]byte
	CRCHash          [20]byte // unused
	NumTelemetryKeys uint8    // unused
}

// https://gtker.com/wow_messages/docs/cmd_auth_reconnect_proof_client.html
type ReconnectProofPacket struct {
	Opcode         byte // 0x3
	ProofData      [16]byte
	ClientProof    [20]byte
	ClientChecksum [20]byte // unused
	KeyCount       byte     // unused
}
