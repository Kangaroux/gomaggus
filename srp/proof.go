package srp

import (
	"crypto/sha1"
	"strings"
)

func CalculateClientProof(
	username string,
	salt,
	clientPublicKey,
	serverPublicKey,
	sessionKey []byte,
) []byte {
	hUsername := sha1.Sum([]byte(strings.ToUpper(username)))
	h := sha1.New()
	h.Write(xorHash)
	h.Write(hUsername[:])
	h.Write(salt)
	h.Write(clientPublicKey)
	h.Write(serverPublicKey)
	h.Write(sessionKey)
	return h.Sum(nil)
}

func CalculateServerProof(clientPublicKey, clientProof, sessionKey []byte) []byte {
	h := sha1.New()
	h.Write(clientPublicKey)
	h.Write(clientProof)
	h.Write(sessionKey)
	return h.Sum(nil)
}

func CalculateReconnectProof(username string, clientData, serverData, sessionKey []byte) []byte {
	h := sha1.New()
	h.Write([]byte(strings.ToUpper(username)))
	h.Write(clientData)
	h.Write(serverData)
	h.Write(sessionKey)
	return h.Sum(nil)
}

func CalculateWorldProof(username string, clientSeed, serverSeed, sessionKey []byte) []byte {
	h := sha1.New()
	h.Write([]byte(strings.ToUpper(username)))
	h.Write([]byte{0, 0, 0, 0})
	h.Write(clientSeed)
	h.Write(serverSeed)
	h.Write(sessionKey)
	return h.Sum(nil)
}
