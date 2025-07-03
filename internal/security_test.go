package internal

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	input := "Text to encrypt"
	key := RandomKey(16)
	fmt.Printf("Random key: %s\n", key)
	output, err := Encrypt([]byte(input), []byte(key))
	if err != nil {
		t.Errorf("Error to encrypt: %s", err)
	}

	fmt.Println(output)

	output, err = Decrypt(output, key)
	if err != nil {
		t.Errorf("Error to decrypt: %s", err)
	}
	fmt.Println(output)
}
