package cmd

import (
	"encoding/base64"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/sync"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt <secret-value>",
	Short: "Encrypt a secret value using AES-256-GCM",
	Args:  cobra.ExactArgs(1),
	RunE:  runEncrypt,
}

var decryptFlag bool
var genKeyFlag bool

func init() {
	encryptCmd.Flags().BoolVar(&decryptFlag, "decrypt", false, "decrypt the value instead of encrypting")
	encryptCmd.Flags().BoolVar(&genKeyFlag, "gen-key", false, "generate a new 32-byte base64 key and print it")
	rootCmd.AddCommand(encryptCmd)
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	if genKeyFlag {
		buf := make([]byte, 32)
		if _, err := rand.Read(buf); err != nil {
			return fmt.Errorf("key gen: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), base64.StdEncoding.EncodeToString(buf))
		return nil
	}

	keyB64 := os.Getenv("VAULTPULL_ENCRYPT_KEY")
	enc, err := sync.NewEncryptor(keyB64)
	if err != nil {
		return err
	}

	if decryptFlag {
		plain, err := enc.Decrypt(args[0])
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), plain)
		return nil
	}

	ciphertext, err := enc.Encrypt(args[0])
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), ciphertext)
	return nil
}
