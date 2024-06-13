package packets

import "github.com/kangaroux/gomaggus/internal/srp"

// https://gtker.com/wow_messages/docs/cmd_auth_logon_challenge_server.html#protocol-version-8
// FIELD ORDER MATTERS, DO NOT REORDER
type ServerLoginChallenge struct {
	Opcode          byte
	ProtocolVersion byte
	ErrorCode       byte
	PublicKey       [srp.KeySize]byte
	GeneratorSize   byte
	Generator       byte
	LargePrimeSize  byte
	LargePrime      [srp.LargePrimeSize]byte
	Salt            [srp.SaltSize]byte
	CrcHash         [16]byte

	// Using any flags would require additional fields but this is set to zero for now
	SecurityFlags byte
}
