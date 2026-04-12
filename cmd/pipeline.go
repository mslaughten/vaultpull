package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/sync"
)

func init() {
	var stages []string
	var inputFile string
	var outputFormat string

	pipelineCmd := &cobra.Command{
		Use:   "pipeline <path>",
		Short: "Run a named transformation pipeline against a Vault secret path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := os.ReadFile(inputFile)
			if err != nil {
				return fmt.Errorf("read input: %w", err)
			}
			var secrets map[string]string
			if err := json.Unmarshal(raw, &secrets); err != nil {
				return fmt.Errorf("parse input JSON: %w", err)
			}

			p := sync.NewPipeline()
			for _, s := range stages {
				name := strings.TrimSpace(s)
				switch name {
				case "upper":
					p.AddStage(sync.PipelineStage{
						Name: "upper",
						Apply: func(m map[string]string) (map[string]string, error) {
							out := make(map[string]string, len(m))
							for k, v := range m {
								out[k] = strings.ToUpper(v)
							}
							return out, nil
						},
					})
				default:
					return fmt.Errorf("unknown stage %q", name)
				}
			}

			out, err := p.Run(secrets)
			if err != nil {
				return err
			}

			switch outputFormat {
			case "json":
				return json.NewEncoder(os.Stdout).Encode(out)
			default:
				for k, v := range out {
					fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
				}
			}
			return nil
		},
	}

	pipelineCmd.Flags().StringSliceVar(&stages, "stage", nil, "ordered list of pipeline stages (e.g. upper)")
	pipelineCmd.Flags().StringVar(&inputFile, "input", "", "path to JSON file containing secrets (required)")
	pipelineCmd.Flags().StringVar(&outputFormat, "output", "dotenv", "output format: dotenv or json")
	_ = pipelineCmd.MarkFlagRequired("input")

	rootCmd.AddCommand(pipelineCmd)
}
