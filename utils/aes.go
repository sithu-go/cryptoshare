package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"cryptoshare/conf"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

var (
	key = []byte(conf.AESKey)
)

func EncryptAES(plaintext string) (string, error) {

	if !IsPrivateKey(plaintext) {
		return "", errors.New("invalid private key")
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())

	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}
	outData := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return hex.EncodeToString(outData), nil
}

func DecryptAES(ct string) (string, error) {
	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("invalid NonceSize")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
