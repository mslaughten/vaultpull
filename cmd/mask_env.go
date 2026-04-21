package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/sync"
)

func init() {
	var mode string
	var reveal int
	var symbol string
	var keys []string

	cmd := &cobra.Command{
		Use:   "mask-env <file>",
		Short: "Print .env file with secret values masked",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMaskEnv(args[0], mode, reveal, symbol, keys)
		},
	}

	cmd.Flags().StringVar(&mode, "mode", "all", "masking mode: all, suffix, prefix")
	cmd.Flags().IntVar(&reveal, "reveal", 0, "number of characters to reveal")
	cmd.Flags().StringVar(&symbol, "symbol", "*", "mask symbol character")
	cmd.Flags().StringSliceVar(&keys, "keys", nil, "keys to mask (default: all)")

	rootCmd.AddCommand(cmd)
}

func runMaskEnv(file, mode string, reveal int, symbol string, keys []string) error {
	reader, err := envfile.NewReader(file)
	if err != nil {
		return fmt.Errorf("mask-env: %w", err)
	}
	secrets, err := reader.Read()
	if err != nil {
		return fmt.Errorf("mask-env: %w", err)
	}

	r, err := sync.NewMaskEnvRenderer(sync.MaskEnvOptions{
		Mode:        sync.MaskEnvMode(strings.ToLower(mode)),
		RevealChars: reveal,
		MaskSymbol:  symbol,
		Keys:        keys,
	}, os.Stdout)
	if err != nil {
		return fmt.Errorf("mask-env: %w", err)
	}
	return r.Render(secrets)
}
