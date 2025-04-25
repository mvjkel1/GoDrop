package transfer_test

import (
	"bytes"
	"testing"

	"godrop/cmd/internal/transfer"
)

func TestEncryptDecrypt(t *testing.T) {
	original := []byte("Hello, GoDrop!")
	pass := "test-pass"

	nonce, encrypted, err := transfer.Encrypt(original, pass)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := transfer.Decrypt(nonce, encrypted, pass)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(original, decrypted) {
		t.Errorf("Decrypted data does not match original.\nExpected: %s\nGot: %s", original, decrypted)
	}
}

func TestDecryptWrongPassphrase(t *testing.T) {
	data := []byte("Secret data")
	rightPass := "correct123"
	wrongPass := "wrong456"

	nonce, encrypted, err := transfer.Encrypt(data, rightPass)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := transfer.Decrypt(nonce, encrypted, wrongPass)
	if err == nil && bytes.Equal(data, decrypted) {
		t.Error("Decryption should fail with wrong passphrase")
	}
}
