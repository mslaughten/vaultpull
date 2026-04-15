package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/sync"
)

func init() {
	sampleCmd := &cobra.Command{
		Use:   "sample <env-file>",
		Short: "Sample a subset of keys from a .env file",
		Args:  cobra.ExactArgs(1),
		RunE:  runSample,
	}
	sampleCmd.Flags().Int("n", 5, "number of keys to sample")
	sampleCmd.Flags().String("strategy", "first", "sampling strategy: random|first|last")
	sampleCmd.Flags().Bool("dry-run", false, "print sampled keys without writing")
	root.AddCommand(sampleCmd)
}

func runSample(cmd *cobra.Command, args []string) error {
	path := args[0]
	n, _ := cmd.Flags().GetInt("n")
	strategyStr, _ := cmd.Flags().GetString("strategy")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	strategy, err := sync.SampleStrategyFromString(strategyStr)
	if err != nil {
		return err
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	sampler, err := sync.NewSampler(n, strategy, rng)
	if err != nil {
		return err
	}

	reader := envfile.NewReader(path)
	secrets, err := reader.Read()
	if err != nil {
		return fmt.Errorf("sample: read %s: %w", path, err)
	}

	out := sampler.Apply(secrets)

	if dryRun {
		for k, v := range out {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
		return nil
	}

	writer := envfile.NewWriter(path)
	return writer.Write(out)
}
