package srp

import (
	"crypto/sha1"
	"math/big"
	"strings"

	"github.com/kangaroux/gomaggus/internal"
)

func CalculateX(username, password string, salt []byte) []byte {
	h := sha1.New()
	inner := sha1.Sum([]byte(strings.ToUpper(username) + ":" + strings.ToUpper(password)))
	h.Write(salt)
	h.Write(inner[:])
	return h.Sum(nil)
}

func CalculateVerifier(username, password string, salt []byte) []byte {
	x := BytesToInt(CalculateX(username, password, salt))
	return IntToBytes(VerifierSize, big.NewInt(0).Exp(g, x, n))
}

func CalculateServerPublicKey(verifier []byte, serverPrivateKey []byte) []byte {
	publicKey := big.NewInt(0).Exp(g, BytesToInt(serverPrivateKey), n)
	kv := big.NewInt(0).Mul(k, BytesToInt(verifier))
	return IntToBytes(KeySize, publicKey.Add(publicKey, kv).Mod(publicKey, n))
}

func CalculateU(clientPublicKey, serverPublicKey []byte) []byte {
	h := sha1.New()
	h.Write(clientPublicKey)
	h.Write(serverPublicKey)
	return h.Sum(nil)
}

func CalculateServerSKey(clientPublicKey, verifier, u, serverPrivateKey []byte) []byte {
	S := big.NewInt(0).Exp(BytesToInt(verifier), BytesToInt(u), n)
	S.Mul(S, BytesToInt(clientPublicKey))
	S.Exp(S, BytesToInt(serverPrivateKey), n)
	return IntToBytes(KeySize, S)
}

func CalculateInterleave(S []byte) []byte {
	for len(S) > 0 && S[0] == 0 {
		S = S[2:]
	}

	lenS := len(S)
	even, odd := make([]byte, lenS/2), make([]byte, lenS/2)

	for i := 0; i < lenS/2; i++ {
		even[i] = S[i*2]
		odd[i] = S[i*2+1]
	}

	hEven := sha1.Sum(even)
	hOdd := sha1.Sum(odd)
	interleaved := make([]byte, 40)

	for i := 0; i < 20; i++ {
		interleaved[i*2] = hEven[i]
		interleaved[i*2+1] = hOdd[i]
	}

	return interleaved
}

func CalculateServerSessionKey(clientPublicKey, serverPublicKey, serverPrivateKey, verifier []byte) []byte {
	u := CalculateU(clientPublicKey, serverPublicKey)
	S := CalculateServerSKey(clientPublicKey, verifier, u, serverPrivateKey)
	return CalculateInterleave(S)
}

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

// BytesToInt returns a little endian big integer from a big endian byte array.
func BytesToInt(data []byte) *big.Int {
	return big.NewInt(0).SetBytes(internal.Reverse(data))
}

// IntToBytes returns a big endian byte array from a little endian big integer.
func IntToBytes(padding int, bi *big.Int) []byte {
	return internal.Reverse(internal.Pad(padding, bi.Bytes()))
}
