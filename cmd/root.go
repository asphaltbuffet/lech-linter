package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/asphaltbuffet/lech-linter/internal/checker"
	"github.com/asphaltbuffet/lech-linter/internal/version"
)

// NewRootCmd creates a new instance of the root command.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:          "lech-linter [file ...]",
		Short:        "Detect misspellings of Lechlitner",
		Version:      version.Version,
		SilenceUsage: true,
		Args:         cobra.ArbitraryArgs,
		RunE:         runRootCmd,
	}

	rootCmd.AddCommand(NewManCmd())

	return rootCmd
}

// Execute runs the root command.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func runRootCmd(cmd *cobra.Command, args []string) error {
	var findings []checker.Finding

	if len(args) == 0 {
		// Read from stdin.
		found, err := checker.Scan(os.Stdin, "<stdin>")
		if err != nil {
			return fmt.Errorf("reading stdin: %w", err)
		}

		findings = append(findings, found...)
	} else {
		for _, path := range args {
			f, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("opening %s: %w", path, err)
			}

			found, err := checker.Scan(f, path)
			// Possible edge case if checker.Scan panics, but better to close the file immediately rather than waiting for defer
			_ = f.Close()

			if err != nil {
				return fmt.Errorf("scanning %s: %w", path, err)
			}

			findings = append(findings, found...)
		}
	}

	for _, f := range findings {
		fmt.Fprintf(cmd.OutOrStdout(), "%s:%d:%d: possible misspelling %q (did you mean \"Lechlitner\"?)\n",
			f.File, f.Line, f.Col, f.Word)
	}

	if len(findings) > 0 {
		return fmt.Errorf("%d possible misspellings", len(findings))
	}

	return nil
}
