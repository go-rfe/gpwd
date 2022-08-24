package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func GetCrypto(key []byte) (func([]byte) ([]byte, error), func([]byte) ([]byte, error), error) {
	cryptoCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(cryptoCipher)
	if err != nil {
		return nil, nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce := make([]byte, nonceSize)

	encrypt := func(plaintext []byte) ([]byte, error) {
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return nil, err
		}

		return gcm.Seal(nonce, nonce, plaintext, nil), nil
	}

	decrypt := func(ciphertext []byte) ([]byte, error) {
		return gcm.Open(nil, ciphertext[:nonceSize], ciphertext[nonceSize:], nil)
	}

	return encrypt, decrypt, nil
}
