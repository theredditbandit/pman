package cmd

import (
	"github.com/spf13/cobra"

	"github.com/theredditbandit/pman/pkg/ui"
)

// iCmd represents the interactive command
var iCmd = &cobra.Command{
	Use:     "i",
	Short:   "Launches pman in interactive mode",
	Aliases: []string{"interactive", "iteractive"},
	RunE: func(_ *cobra.Command, _ []string) error {
		return ui.Tui()
	},
}

func init() {
	rootCmd.AddCommand(iCmd)
}
