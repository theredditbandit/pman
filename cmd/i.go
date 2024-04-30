package cmd

import (
	"github.com/spf13/cobra"
	"github.com/theredditbandit/pman/pkg/ui"
	_ "github.com/theredditbandit/pman/pkg/ui/example"
)

// iCmd represents the interactive command
var iCmd = &cobra.Command{
	Use:     "i",
	Short:   "Launches pman in interactive mode",
	Aliases: []string{"interactive", "iteractive"},
	Run: func(cmd *cobra.Command, args []string) {
		// example.RunExampleList()
		ui.Tui()
	},
}

func init() {
	rootCmd.AddCommand(iCmd)
}
