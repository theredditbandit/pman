package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all projects along with their status",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ls called")
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
