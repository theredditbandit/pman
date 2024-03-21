package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new project database/index",
	Long: `Creates a new project database to write to.

    Creates a new index to write to. This is useful when you want to write to a new index and don't want to write to the default index.
    `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new called")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
