package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

func TestEncrypt_Decrypt_Success(t *testing.T) {
	plaintexts := []string{
		"hello world",
		"",
		"1234567890!@#$%^&*()_+-=",
		strings.Repeat("a", 100),
		"测试中文字符",
	}

	for _, pt := range plaintexts {
		encrypted, err := Encrypt(pt)
		fmt.Printf("Encrypt(%q) = %q\n", pt, encrypted)
		if err != nil {
			t.Errorf("Encrypt(%q) returned error: %v", pt, err)
			continue
		}
		if encrypted == "" {
			t.Errorf("Encrypt(%q) returned empty string", pt)
		}
		// Check if encrypted is valid base64
		if _, err := base64.URLEncoding.DecodeString(encrypted); err != nil {
			t.Errorf("Encrypt(%q) returned invalid base64: %v", pt, err)
		}

		// Decrypt and check round-trip
		decrypted, err := Decrypt(encrypted)
		fmt.Printf("Decrypt(Encrypt(%q)) = %q\n", pt, decrypted)
		if err != nil {
			t.Errorf("Decrypt(Encrypt(%q)) returned error: %v", pt, err)
		}
		if decrypted != pt {
			t.Errorf("Decrypt(Encrypt(%q)) = %q, want %q", pt, decrypted, pt)
		}
	}
}

func TestEncrypt_ErrorOnInvalidKey(t *testing.T) {
	// Backup and replace secretKey with invalid length
	origKey := secretKey
	defer func() { secretKey = origKey }()
	secretKey = []byte("shortkey")

	_, err := Encrypt("test")
	if err == nil {
		t.Error("Encrypt should return error with invalid key length")
	}
}
