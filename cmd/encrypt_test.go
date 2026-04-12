package cmd

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

func TestEncryptCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Use == "encrypt <secret-value>" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("encrypt command not registered on root")
	}
}

func TestEncryptCmd_RequiresOneArg(t *testing.T) {
	rootCmd.SetArgs([]string{"encrypt"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing arg")
	}
}

func TestEncryptCmd_GenKey_PrintsBase64(t *testing.T) {
	var buf bytes.Buffer
	encryptCmd.SetOut(&buf)
	encryptCmd.SetArgs([]string{"--gen-key", "placeholder"})
	// reset flags
	genKeyFlag = true
	if err := runEncrypt(encryptCmd, []string{"placeholder"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := strings.TrimSpace(buf.String())
	decoded, err := base64.StdEncoding.DecodeString(output)
	if err != nil {
		t.Fatalf("output is not valid base64: %v", err)
	}
	if len(decoded) != 32 {
		t.Fatalf("expected 32 bytes, got %d", len(decoded))
	}
	genKeyFlag = false
}

func TestEncryptCmd_DefaultFlags(t *testing.T) {
	if encryptCmd.Flag("decrypt") == nil {
		t.Error("missing --decrypt flag")
	}
	if encryptCmd.Flag("gen-key") == nil {
		t.Error("missing --gen-key flag")
	}
}

func TestEncryptCmd_MissingEnvKey_ReturnsError(t *testing.T) {
	t.Setenv("VAULTPULL_ENCRYPT_KEY", "")
	err := runEncrypt(encryptCmd, []string{"some-value"})
	if err == nil {
		t.Fatal("expected error when VAULTPULL_ENCRYPT_KEY is empty")
	}
}
