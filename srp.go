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

var (
	_bigN = bytesToBig([]byte{
		0x89, 0x4b, 0x64, 0x5e, 0x89, 0xe1, 0x53, 0x5b,
		0xbd, 0xad, 0x5b, 0x8b, 0x29, 0x06, 0x50, 0x53,
		0x08, 0x01, 0xb1, 0x8e, 0xbf, 0xbf, 0x5e, 0x8f,
		0xab, 0x3c, 0x82, 0x87, 0x2a, 0x3e, 0x9b, 0xb7,
	})
)

func bigN() *big.Int {
	return cloneBig(_bigN)
}

func bigG() *big.Int {
	return big.NewInt(7)
}

func bigK() *big.Int {
	return big.NewInt(3)
}

func cloneBig(v *big.Int) *big.Int {
	return big.NewInt(0).Set(v)
}

func bytesToBig(data []byte) *big.Int {
	return big.NewInt(0).SetBytes(data)
}

func ReverseBytes(data []byte) []byte {
	n := len(data)
	for i := 0; i < n/2; i++ {
		data[i], data[n-i-1] = data[n-i-1], data[i]
	}
	return data
}

// Returns little endian
func passVerify(username string, password string, salt []byte) []byte {
	x := bytesToBig(ReverseBytes(calcX(username, password, salt)))
	verifier := big.NewInt(0).Exp(bigG(), x, bigN())
	return ReverseBytes(verifier.Bytes())
}

func calcX(username string, password string, salt []byte) []byte {
	h1 := sha1.New()
	h2 := sha1.New()

	// sha1(username | ":" | password)
	io.WriteString(h1, username+":"+password)

	// sha1(salt | sha1(username | ":" | password))
	h2.Write(salt)
	h2.Write(h1.Sum(nil))

	return h2.Sum(nil)
}

// Big endian args + return
func calcServerPublicKey(verifier []byte, serverPrivateKey []byte) []byte {
	result := big.NewInt(0)
	result.Mul(bigK(), bytesToBig(verifier))
	result.Add(result, big.NewInt(0).Exp(bigG(), bytesToBig(serverPrivateKey), bigN()))
	result.Mod(result, bigN())
	return result.Bytes()
}
