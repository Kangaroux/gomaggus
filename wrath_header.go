package main

import (
	"crypto"
	"crypto/hmac"
)

var (
	fixedDecryptKey = NewByteArray([]byte{
		0xC2, 0xB3, 0x72, 0x3C, 0xC6, 0xAE, 0xD9, 0xB5,
		0x34, 0x3C, 0x53, 0xEE, 0x2F, 0x43, 0x67, 0xCE,
	}, 16, false)
	fixedEncryptKey = NewByteArray([]byte{
		0xCC, 0x98, 0xAE, 0x04, 0xE8, 0x97, 0xEA, 0xCA,
		0x12, 0xDD, 0xC0, 0x93, 0x42, 0x91, 0x53, 0x57,
	}, 16, false)
)

type WrathHeaderCrypto struct {
	decryptKey []byte
	encryptKey []byte
	sessionKey *ByteArray
}

func NewWrathHeaderCrypto(sessionKey *ByteArray) *WrathHeaderCrypto {
	ret := &WrathHeaderCrypto{
		sessionKey: sessionKey,
	}

	ret.decryptKey = ret.generateKey(fixedDecryptKey)
	ret.encryptKey = ret.generateKey(fixedEncryptKey)

	return ret
}

func (h *WrathHeaderCrypto) generateKey(fixedKey *ByteArray) []byte {
	hash := hmac.New(crypto.SHA1.New, fixedKey.LittleEndian().Bytes())
	hash.Write(h.sessionKey.LittleEndian().Bytes())
	return hash.Sum(nil)
}
