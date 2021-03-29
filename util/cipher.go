package util

// based on: https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-golang-application-crypto-packages/

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
)

type Cipher struct {
	key  string `json:"key"`
	hash string `json:"hash"`
}

//
func NewCipher(key string) (*Cipher, error) {
	c := &Cipher{key: key}
	c.hash = c.createHash(key)

	return c, nil
}

// Encryp/Decrypt Key must be 32 bytes
func (c *Cipher) createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))

	return hex.EncodeToString(hasher.Sum(nil))
}

//
func (c *Cipher) Encrypt(data []byte) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(c.hash))
	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

//
func (c *Cipher) Decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(c.hash))

	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

//
func (c *Cipher) EncryptFile(filename string, data []byte) error {
	f, _ := os.Create(filename)

	defer f.Close()

	cipherdata, err := c.Encrypt(data)

	if err != nil {
		return err
	}

	_, err = f.Write(cipherdata)

	return err
}

//
func (c *Cipher) DecryptFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return c.Decrypt(data)
}
