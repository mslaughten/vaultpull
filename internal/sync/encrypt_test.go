package sync

import (
	"encoding/base64"
	"strings"
	"testing"
)

func validBase64Key() string {
	return base64.StdEncoding.EncodeToString([]byte("12345678901234567890123456789012"))
}

func TestNewEncryptor_Valid(t *testing.T) {
	_, err := NewEncryptor(validBase64Key())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewEncryptor_EmptyKey(t *testing.T) {
	_, err := NewEncryptor("")
	if err == nil || !strings.Contains(err.Error(), "empty") {
		t.Fatalf("expected empty key error, got %v", err)
	}
}

func TestNewEncryptor_InvalidBase64(t *testing.T) {
	_, err := NewEncryptor("not-valid-base64!!!")
	if err == nil || !strings.Contains(err.Error(), "invalid base64") {
		t.Fatalf("expected base64 error, got %v", err)
	}
}

func TestNewEncryptor_WrongKeyLength(t *testing.T) {
	short := base64.StdEncoding.EncodeToString([]byte("tooshort"))
	_, err := NewEncryptor(short)
	if err == nil || !strings.Contains(err.Error(), "32 bytes") {
		t.Fatalf("expected key length error, got %v", err)
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	enc, err := NewEncryptor(validBase64Key())
	if err != nil {
		t.Fatal(err)
	}
	plaintext := "super-secret-value"
	ciphertext, err := enc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if ciphertext == plaintext {
		t.Fatal("ciphertext should differ from plaintext")
	}
	got, err := enc.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if got != plaintext {
		t.Fatalf("want %q, got %q", plaintext, got)
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	enc, _ := NewEncryptor(validBase64Key())
	_, err := enc.Decrypt("!!!notbase64")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDecrypt_TooShort(t *testing.T) {
	enc, _ := NewEncryptor(validBase64Key())
	_, err := enc.Decrypt(base64.StdEncoding.EncodeToString([]byte("tiny")))
	if err == nil || !strings.Contains(err.Error(), "too short") {
		t.Fatalf("expected too short error, got %v", err)
	}
}

func TestApplyToMap_EncryptsAllValues(t *testing.T) {
	enc, _ := NewEncryptor(validBase64Key())
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result, err := enc.ApplyToMap(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range result {
		if v == secrets[k] {
			t.Errorf("key %q: value was not encrypted", k)
		}
		decrypted, err := enc.Decrypt(v)
		if err != nil {
			t.Errorf("key %q: decrypt error: %v", k, err)
		}
		if decrypted != secrets[k] {
			t.Errorf("key %q: want %q got %q", k, secrets[k], decrypted)
		}
	}
}
