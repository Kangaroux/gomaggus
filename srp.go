/*
This file is based off the SRP6 impl from:
https://gtker.com/implementation-guide-for-the-world-of-warcraft-flavor-of-srp6/
*/
package main

import (
	"crypto/sha1"
	"fmt"
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

// Calculates the server's public key. Returns a little endian byte array.
func calcServerPublicKey(verifier *ByteArray, serverPrivateKey *ByteArray) *ByteArray {
	result := big.NewInt(0)

	// k * v
	result.Mul(bigK(), verifier.BigInt())
	// k * v + (g^b % N)
	result.Add(result, big.NewInt(0).Exp(bigG(), serverPrivateKey.BigInt(), bigN()))
	// (k * v + (g^b % N)) % N
	result.Mod(result, bigN())

	return NewByteArray(result.Bytes(), true).LittleEndian()
}

// Calculates the client's S key, a value that is used to generate the client's session key. Returns
// a little endian byte array.
func calcClientSKey(clientPrivateKey *ByteArray, serverPublicKey *ByteArray, x *ByteArray, u *ByteArray) *ByteArray {
	bigX := x.BigInt()

	// u * x
	exponent := big.NewInt(0).Mul(u.BigInt(), bigX)
	// a + u * x
	exponent.Add(exponent, clientPrivateKey.BigInt())

	// g^x % N
	S := big.NewInt(0).Exp(bigG(), bigX, bigN())
	// k * (g^x % N)
	S.Mul(S, bigK())
	// B - (k * (g^x % N))
	S.Sub(serverPublicKey.BigInt(), S)
	// (B - (k * (g^x % N)))^(a + u * x) % N
	S.Exp(S, exponent, bigN())

	return NewByteArray(S.Bytes(), true).LittleEndian()
}

// Calculates the server's S key, a value that is used to generate the server's session key. Returns
// a little endian byte array.
func calcServerSKey(clientPublicKey *ByteArray, verifier *ByteArray, u *ByteArray, serverPrivateKey *ByteArray) *ByteArray {
	// v^u % N
	S := big.NewInt(0).Exp(verifier.BigInt(), u.BigInt(), bigN())
	// A * (v^u % N)
	S.Mul(S, clientPublicKey.BigInt())
	// (A * (v^u % N))^b % N
	S.Exp(S, serverPrivateKey.BigInt(), bigN())

	return NewByteArray(S.Bytes(), true).LittleEndian()
}

// Calculates U, a hash used for generating session keys. Returns a BIG endian byte array.
func calcU(clientPublicKey *ByteArray, serverPublicKey *ByteArray) *ByteArray {
	u := sha1.New()

	// SHA1(clientPublicKey | serverPublicKey)
	u.Write(clientPublicKey.LittleEndian().Bytes())
	u.Write(serverPublicKey.LittleEndian().Bytes())

	return NewByteArray(u.Sum(nil), true)
}

// Prepares the S key to be interleaved. Returns a raw little endian byte array.
func prepareInterleave(S *ByteArray) []byte {
	result := S.LittleEndian().Bytes()

	// Trim the two LSB while the LSB is zero
	for len(result) > 0 && result[0] == 0x0 {
		result = result[2:]
	}

	return result
}

// Interleaves the S
func calcInterleave(S *ByteArray) *ByteArray {
	preparedS := prepareInterleave(S)
	halfSLen := len(preparedS) / 2
	even := make([]byte, halfSLen)
	odd := make([]byte, halfSLen)

	for i := 0; i < halfSLen; i++ {
		even[i] = preparedS[i*2]
		odd[i] = preparedS[i*2+1]
	}

	hEven := sha1.Sum(even)
	hOdd := sha1.Sum(odd)

	result := make([]byte, 40)
	for i := 0; i < 20; i++ {
		result[i*2] = hEven[i]
		result[i*2+1] = hOdd[i]
	}

	return NewByteArray(result, false)
}

// Little endian args + return
func calcServerSessionKey(clientPublicKey *ByteArray, serverPublicKey *ByteArray, verifier *ByteArray, serverPrivateKey *ByteArray) *ByteArray {
	u := calcU(clientPublicKey, serverPublicKey)
	fmt.Printf("u: %x\n", u.Bytes())
	S := calcServerSKey(clientPublicKey, verifier, u, serverPrivateKey)
	fmt.Printf("S: %x\n", S.Bytes())
	return calcInterleave(S)
}
