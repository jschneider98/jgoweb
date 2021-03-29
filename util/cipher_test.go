// +build unit

package util

import (
	"fmt"
	"testing"
)

//
func TestCipher(t *testing.T) {
	test := "Cipher test"
	c, err := NewCipher("test")

	ciphertext, err := c.Encrypt([]byte(test))

	if err != nil {
		t.Error(err)
	}

	plaintext, err := c.Decrypt(ciphertext)

	if err != nil {
		t.Error(err)
	}

	if test != string(plaintext) {
		t.Error(fmt.Sprintf("Cipher failed. Expected: %s, Got: %s\n", test, plaintext))
	}
}
