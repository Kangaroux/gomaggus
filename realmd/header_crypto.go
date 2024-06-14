package realmd

import (
	"crypto"
	"crypto/hmac"
	"crypto/rc4"
	"errors"
)

const (
	// 23 bits + 1 bit for LARGE_HEADER_FLAG
	SizeFieldMaxValue = 0x7FFFFF

	// 15 bits (16th bit is reserved for LARGE_HEADER_FLAG)
	LargeHeaderThreshold = 0x7FFF

	// Set on MSB of size field (first header byte)
	LargeHeaderFlag = 0x80
)

var (
	fixedDecryptKey = []byte{
		0xC2, 0xB3, 0x72, 0x3C, 0xC6, 0xAE, 0xD9, 0xB5,
		0x34, 0x3C, 0x53, 0xEE, 0x2F, 0x43, 0x67, 0xCE,
	}
	fixedEncryptKey = []byte{
		0xCC, 0x98, 0xAE, 0x04, 0xE8, 0x97, 0xEA, 0xCA,
		0x12, 0xDD, 0xC0, 0x93, 0x42, 0x91, 0x53, 0x57,
	}

	ErrCryptoNotInitialized = errors.New("header crypto has not been initialized")
)

type HeaderCrypto struct {
	decryptCipher *rc4.Cipher
	encryptCipher *rc4.Cipher
	sessionKey    []byte
}

func NewHeaderCrypto(sessionKey []byte) *HeaderCrypto {
	return &HeaderCrypto{sessionKey: sessionKey}
}

func (h *HeaderCrypto) Init() error {
	return h.InitKeys(fixedDecryptKey, fixedEncryptKey)
}

func (h *HeaderCrypto) InitKeys(decryptKey, encryptKey []byte) error {
	decryptCipher, err := rc4.NewCipher(h.GenerateKey(decryptKey))
	if err != nil {
		return err
	}

	encryptCipher, err := rc4.NewCipher(h.GenerateKey(encryptKey))
	if err != nil {
		return err
	}

	h.decryptCipher = decryptCipher
	drop1024(h.decryptCipher)

	h.encryptCipher = encryptCipher
	drop1024(h.encryptCipher)

	return nil
}

func (h *HeaderCrypto) Decrypt(data []byte) error {
	if h.decryptCipher == nil {
		return ErrCryptoNotInitialized
	}

	h.decryptCipher.XORKeyStream(data, data)
	return nil
}

func (h *HeaderCrypto) Encrypt(data []byte) error {
	if h.encryptCipher == nil {
		return ErrCryptoNotInitialized
	}

	h.encryptCipher.XORKeyStream(data, data)
	return nil
}

func (h *HeaderCrypto) GenerateKey(fixedKey []byte) []byte {
	hash := hmac.New(crypto.SHA1.New, fixedKey)
	hash.Write(h.sessionKey)
	return hash.Sum(nil)
}

func drop1024(cipher *rc4.Cipher) {
	var drop1024 [1024]byte
	cipher.XORKeyStream(drop1024[:], drop1024[:])
}
