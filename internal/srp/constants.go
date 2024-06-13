package srp

import (
	"math/big"

	"github.com/kangaroux/gomaggus/internal"
)

const (
	ProofDataSize  = 16
	ProofSize      = 20
	SaltSize       = 32
	VerifierSize   = 32
	KeySize        = 32
	LargePrimeSize = 32
	SessionKeySize = 40

	Generator = 7
)

var (
	largeSafePrime = LargePrime()

	xorHash = []byte{
		0xDD, 0x7B, 0xB0, 0x3A, 0x38, 0xAC, 0x73, 0x11,
		0x03, 0x98, 0x7C, 0x5A, 0x50, 0x6F, 0xCA, 0x96,
		0x6C, 0x7B, 0xC2, 0xA7,
	}

	n = BytesToInt(largeSafePrime)
	g = big.NewInt(Generator)
	k = big.NewInt(3)
)

// LargePrime returns a little endian byte array.
func LargePrime() []byte {
	return internal.Reverse([]byte{
		0x89, 0x4B, 0x64, 0x5E, 0x89, 0xE1, 0x53, 0x5B,
		0xBD, 0xAD, 0x5B, 0x8B, 0x29, 0x06, 0x50, 0x53,
		0x08, 0x01, 0xB1, 0x8E, 0xBF, 0xBF, 0x5E, 0x8F,
		0xAB, 0x3C, 0x82, 0x87, 0x2A, 0x3E, 0x9B, 0xB7,
	})
}
