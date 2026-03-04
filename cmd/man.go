package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// NewManCmd creates a new instance of the man command.
func NewManCmd() *cobra.Command {
	manCmd := &cobra.Command{
		Use:    "man",
		Short:  "Generate man page to stdout",
		Hidden: true,
		RunE:   runManCmd,
	}

	return manCmd
}

func runManCmd(cmd *cobra.Command, _ []string) error {
	header := &doc.GenManHeader{
		Date:    nil,
		Title:   "LECH-LINTER",
		Section: "1",
		Source:  "",
		Manual:  "",
	}

	return doc.GenMan(NewRootCmd(), header, cmd.OutOrStdout())
}
