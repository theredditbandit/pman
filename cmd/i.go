package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// iCmd represents the interactive command
var iCmd = &cobra.Command{
	Use:     "i",
	Short:   "Launches pman in interactive mode",
	Aliases: []string{"interactive", "iteractive"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("not implemented")
	},
}

func init() {
	rootCmd.AddCommand(iCmd)
}
