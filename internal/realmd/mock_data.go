package realmd

import (
	"crypto/rand"
	"log"

	"github.com/kangaroux/gomaggus/srp"
)

const (
	MOCK_USERNAME = "TEST"
	MOCK_PASSWORD = "PASSWORD"
)

var (
	MOCK_SALT        []byte
	MOCK_VERIFIER    []byte
	MOCK_PRIVATE_KEY []byte
	MOCK_PUBLIC_KEY  []byte
)

func init() {
	MOCK_SALT = make([]byte, 32)
	if _, err := rand.Read(MOCK_SALT); err != nil {
		log.Fatalf("error generating salt: %v\n", err)
	}

	MOCK_VERIFIER = srp.CalculateVerifier(MOCK_USERNAME, MOCK_PASSWORD, MOCK_SALT)
	MOCK_PRIVATE_KEY = srp.NewPrivateKey()
	MOCK_PUBLIC_KEY = srp.CalculateServerPublicKey(MOCK_VERIFIER, MOCK_PRIVATE_KEY)
}
