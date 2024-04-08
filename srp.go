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

var SRP_N big.Int

const (
	// Generator
	SRP_G = 7
	// K-Value
	SRP_K = 3
)

func init() {
	// Large safe prime
	SRP_N.SetBytes([]byte{
		0xb7, 0x9b, 0x3e, 0x2a, 0x87, 0x82, 0x3c, 0xab,
		0x8f, 0x5e, 0xbf, 0xbf, 0x8e, 0xb1, 0x01, 0x08,
		0x53, 0x50, 0x06, 0x29, 0x8b, 0x5b, 0xad, 0xbd,
		0x5b, 0x53, 0xe1, 0x89, 0x5e, 0x64, 0x4b, 0x89,
	})
}

func passVerify(username string, password string, salt []byte) []byte {
	// x := calcX(username, password, salt)
	// f := big.NewInt(0)
	// f.Exp()
	return []byte{}
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
