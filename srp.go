/*
This file is based off the SRP6 impl from:
https://gtker.com/implementation-guide-for-the-world-of-warcraft-flavor-of-srp6/
*/
package main

import (
	"crypto/sha1"
	"io"
	"math/big"
)

func bigN() *big.Int {
	return bytesToBig([]byte{
		0x89, 0x4b, 0x64, 0x5e, 0x89, 0xe1, 0x53, 0x5b,
		0xbd, 0xad, 0x5b, 0x8b, 0x29, 0x06, 0x50, 0x53,
		0x08, 0x01, 0xb1, 0x8e, 0xbf, 0xbf, 0x5e, 0x8f,
		0xab, 0x3c, 0x82, 0x87, 0x2a, 0x3e, 0x9b, 0xb7,
	})
}

func bigG() *big.Int {
	return big.NewInt(7)
}

func bigK() *big.Int {
	return big.NewInt(3)
}

func bytesToBig(data []byte) *big.Int {
	return big.NewInt(0).SetBytes(data)
}

func ReverseBytes(data []byte) []byte {
	n := len(data)
	result := make([]byte, n)
	for i := 0; i < n/2; i++ {
		result[i], result[n-i-1] = data[n-i-1], data[i]
	}
	return result
}

// Password verifier. Returns a little endian byte array.
func passVerify(username string, password string, salt *ByteArray) *ByteArray {
	x := calcX(username, password, salt).BigInt()

	// g^x % N
	verifier := big.NewInt(0).Exp(bigG(), x, bigN())

	return NewByteArray(verifier.Bytes(), true).LittleEndian()
}

// Calculates x, a hash used by the password verifier. Returns a little endian byte array.
func calcX(username string, password string, salt *ByteArray) *ByteArray {
	h1 := sha1.New()
	h2 := sha1.New()

	// SHA1(username | ":" | password)
	io.WriteString(h1, username+":"+password)

	// SHA1(salt | SHA1(username | ":" | password))
	h2.Write(salt.LittleEndian().Bytes())
	h2.Write(h1.Sum(nil))

	return NewByteArray(h2.Sum(nil), true).LittleEndian()
}

// Big endian args + return
func calcServerPublicKey(verifier []byte, serverPrivateKey []byte) []byte {
	result := big.NewInt(0)

	// k * v
	result.Mul(bigK(), bytesToBig(verifier))
	// k * v + (g^b % N)
	result.Add(result, big.NewInt(0).Exp(bigG(), bytesToBig(serverPrivateKey), bigN()))
	// (k * v + (g^b % N)) % N
	result.Mod(result, bigN())

	return result.Bytes()
}

// Big endian args + return
func calcClientSKey(clientPrivateKey []byte, serverPublicKey []byte, x []byte, u []byte) []byte {
	bigX := bytesToBig(x)

	// u * x
	exponent := big.NewInt(0).Mul(bytesToBig(u), bigX)
	// a + u * x
	exponent.Add(exponent, bytesToBig(clientPrivateKey))

	// g^x % N
	S := big.NewInt(0).Exp(bigG(), bigX, bigN())
	// k * (g^x % N)
	S.Mul(S, bigK())
	// B - (k * (g^x % N))
	S.Sub(bytesToBig(serverPublicKey), S)
	// (B - (k * (g^x % N)))^(a + u * x) % N
	S.Exp(S, exponent, bigN())

	return S.Bytes()
}

// Big endian args + return
func calcServerSKey(clientPublicKey []byte, verifier []byte, u []byte, serverPrivateKey []byte) []byte {
	// v^u % N
	S := big.NewInt(0).Exp(bytesToBig(verifier), bytesToBig(u), bigN())
	// A * (v^u % N)
	S.Mul(S, bytesToBig(clientPublicKey))
	// (A * (v^u % N))^b % N
	S.Exp(S, bytesToBig(serverPrivateKey), bigN())
	return S.Bytes()
}

// Little endian args + return
func calcU(clientPublicKey []byte, serverPublicKey []byte) []byte {
	u := sha1.New()

	// SHA1(clientPublicKey | serverPublicKey)
	u.Write(clientPublicKey)
	u.Write(serverPublicKey)

	return u.Sum(nil)
}

func splitSKey(S []byte) []byte {
	for len(S) > 0 && S[0] == 0x0 {
		S = S[2:]
	}
	return S
}

// Little endian args + return
func calcInterleave(S []byte) []byte {
	S = splitSKey(S)
	halfSLen := len(S) / 2
	even := make([]byte, halfSLen)
	odd := make([]byte, halfSLen)

	for i := 0; i < halfSLen; i++ {
		even[i] = S[i*2]
		odd[i] = S[i*2+1]
	}

	hEven := sha1.Sum(even)
	hOdd := sha1.Sum(odd)

	result := make([]byte, 40)
	for i := 0; i < 20; i++ {
		result[i*2] = hEven[i]
		result[i*2+1] = hOdd[i]
	}
	return result
}

// Little endian args + return
func calcServerSessionKey(clientPublicKey []byte, serverPublicKey []byte, verifier []byte, serverPrivateKey []byte) []byte {
	u := calcU(ReverseBytes(clientPublicKey), ReverseBytes(serverPublicKey))
	S := ReverseBytes(calcServerSKey(clientPublicKey, verifier, ReverseBytes(u), serverPrivateKey))
	return calcInterleave(S)
}
