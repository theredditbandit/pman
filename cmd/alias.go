package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// aliasCmd represents the alias command
var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("alias called")
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
