package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var delCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a project from the index database. This does not delete the project from the filesystem",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("del called")
	},
	Aliases: []string{"del", "d"},
}

func init() {
	rootCmd.AddCommand(delCmd)
}
