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
	"strings"
)

type BigInteger = *big.Int

func bigN() BigInteger {
	return bytesToBig([]byte{
		0x89, 0x4b, 0x64, 0x5e, 0x89, 0xe1, 0x53, 0x5b,
		0xbd, 0xad, 0x5b, 0x8b, 0x29, 0x06, 0x50, 0x53,
		0x08, 0x01, 0xb1, 0x8e, 0xbf, 0xbf, 0x5e, 0x8f,
		0xab, 0x3c, 0x82, 0x87, 0x2a, 0x3e, 0x9b, 0xb7,
	})
}

func bigG() BigInteger {
	return big.NewInt(7)
}

func bigK() BigInteger {
	return big.NewInt(3)
}

func bytesToBig(data []byte) BigInteger {
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

// Calculates x, a hash used by the password verifier. Returns a big integer.
func calcX(username string, password string, salt *ByteArray) BigInteger {
	h1 := sha1.New()
	h2 := sha1.New()

	// SHA1(username | ":" | password)
	io.WriteString(h1, strings.ToUpper(username)+":"+strings.ToUpper(password))

	// SHA1(salt | SHA1(username | ":" | password))
	h2.Write(salt.LittleEndian().Bytes())
	h2.Write(h1.Sum(nil))

	return NewByteArray(h2.Sum(nil), 20, true).LittleEndian().BigInt()
}

// Calculates the password verifier. Returns a big integer.
func passVerify(username string, password string, salt *ByteArray) BigInteger {
	// g^x % N
	x := calcX(username, password, salt)
	verifier := big.NewInt(0).Exp(bigG(), x, bigN())
	return NewByteArray(verifier.Bytes(), 32, true).LittleEndian().BigInt()
}

// Calculates the server's public key. The arguments must be little endian. Returns a big integer.
func calcServerPublicKey(verifier BigInteger, serverPrivateKey BigInteger) BigInteger {
	result := big.NewInt(0)

	// k * v
	result.Mul(bigK(), verifier)
	// k * v + (g^b % N)
	result.Add(result, big.NewInt(0).Exp(bigG(), serverPrivateKey, bigN()))
	// (k * v + (g^b % N)) % N
	result.Mod(result, bigN())

	return NewByteArray(result.Bytes(), 32, true).LittleEndian().BigInt()
}

// Calculates the client's S key, a value that is used to generate the client's session key. The
// arguments must be little endian. Returns a little endian byte array.
func calcClientSKey(clientPrivateKey BigInteger, serverPublicKey BigInteger, x BigInteger, u BigInteger) *ByteArray {
	// u * x
	exponent := big.NewInt(0).Mul(u, x)
	// a + u * x
	exponent.Add(exponent, clientPrivateKey)

	// g^x % N
	S := big.NewInt(0).Exp(bigG(), x, bigN())
	// k * (g^x % N)
	S.Mul(S, bigK())
	// B - (k * (g^x % N))
	S.Sub(serverPublicKey, S)
	// (B - (k * (g^x % N)))^(a + u * x) % N
	S.Exp(S, exponent, bigN())

	return NewByteArray(S.Bytes(), 32, true).LittleEndian()
}

// Calculates the server's S key, a value that is used to generate the server's session key. The
// arguments must be little endian. Returns a little endian byte array.
func calcServerSKey(clientPublicKey BigInteger, verifier BigInteger, u BigInteger, serverPrivateKey BigInteger) *ByteArray {
	// v^u % N
	S := big.NewInt(0).Exp(verifier, u, bigN())
	// A * (v^u % N)
	S.Mul(S, clientPublicKey)
	// (A * (v^u % N))^b % N
	S.Exp(S, serverPrivateKey, bigN())

	return NewByteArray(S.Bytes(), 32, true).LittleEndian()
}

// Calculates U, a hash used for generating session keys. The arguments must be little endian.
// Returns a big integer.
func calcU(clientPublicKey BigInteger, serverPublicKey BigInteger) BigInteger {
	u := sha1.New()

	// SHA1(clientPublicKey | serverPublicKey)
	u.Write(clientPublicKey.Bytes())
	u.Write(serverPublicKey.Bytes())

	return NewByteArray(u.Sum(nil), 20, true).LittleEndian().BigInt()
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

// Interleaves the S key which generates the session key. Returns a BIG endian byte array.
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

	return NewByteArray(result, 40, true)
}

// Calculates the server's session key. The arguments must be little endian. Returns a little endian
// byte array.
func calcServerSessionKey(
	clientPublicKey BigInteger,
	serverPublicKey BigInteger,
	verifier BigInteger,
	serverPrivateKey BigInteger,
) *ByteArray {
	u := calcU(clientPublicKey, serverPublicKey)
	fmt.Printf("u: %x\n", u.Bytes())
	S := calcServerSKey(clientPublicKey, verifier, u, serverPrivateKey)
	fmt.Printf("S: %x\n", S.Bytes())
	return calcInterleave(S)
}

func calcClientSessionKey(
	username string,
	password string,
	serverPublicKey BigInteger,
	clientPrivateKey BigInteger,
	clientPublicKey BigInteger,
	salt *ByteArray,
) *ByteArray {
	x := calcX(username, password, salt)
	u := calcU(clientPublicKey, serverPublicKey)
	S := calcClientSKey(clientPrivateKey, serverPublicKey, x, u)
	return calcInterleave(S)
}
