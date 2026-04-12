// Package sync provides synchronisation utilities for vaultpull.
//
// # Encrypt
//
// The Encryptor type provides AES-256-GCM symmetric encryption for secret
// values before they are written to disk or transmitted over the wire.
//
// Usage:
//
//	enc, err := sync.NewEncryptor(os.Getenv("VAULTPULL_ENCRYPT_KEY"))
//	if err != nil { ... }
//
//	encrypted, err := enc.Encrypt("my-secret")
//	plain, err := enc.Decrypt(encrypted)
//
//	// Encrypt an entire secret map:
//	encryptedMap, err := enc.ApplyToMap(secrets)
//
// Keys must be 32 random bytes encoded as base64. Generate one with:
//
//	openssl rand -base64 32
package sync
