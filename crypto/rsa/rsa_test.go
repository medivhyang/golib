package rsa

import (
	"testing"
)

func TestGenerateKeyPairFiles(t *testing.T) {
	if err := GenerateKeyFiles(2048, "public.pem", "private.pem"); err != nil {
		t.Error(err)
	}
}

func TestEncryptAndDecryptWithFile(t *testing.T) {
	if err := GenerateKeyFiles(2048, "public.pem", "private.pem"); err != nil {
		t.Error(err)
		return
	}

	plainText := "hello world"
	cipherText, err := EncryptWithFile([]byte(plainText), "public.pem")
	if err != nil {
		t.Error(err)
		return
	}

	result, err := DecryptWithFile(cipherText, "private.pem")
	if err != nil {
		t.Error(err)
		return
	}

	if string(result) != plainText {
		t.Errorf("want: %s, got: %s", plainText, result)
	}
}
