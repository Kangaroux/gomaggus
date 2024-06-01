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

	MOCK_REALMS = []Realm{
		{
			Type:            REALMTYPE_PVE,
			Locked:          false,
			Flags:           REALMFLAG_NONE,
			Name:            "Test Realm\x00",
			Host:            "localhost:8085\x00",
			Population:      0.01,
			NumCharsOnRealm: 0,
			Region:          REALMREGION_US,
			Id:              0,
			// Version:         RealmVersion{Major: 4, Minor: 3, Patch: 6, Build: 12340},
		},
		// {
		// 	Type:            REALMTYPE_PVP,
		// 	Locked:          false,
		// 	Flags:           REALMFLAG_NONE,
		// 	Name:            "Test Realm1\x00",
		// 	Host:            "localhost:8085\x00",
		// 	Population:      0,
		// 	NumCharsOnRealm: 0,
		// 	Region:          REALMREGION_US,
		// 	Id:              1,
		// 	// Version:         RealmVersion{Major: 4, Minor: 3, Patch: 6, Build: 12340},
		// },
	}
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
