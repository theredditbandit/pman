package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "Quicky change directory to the project name",
	Long:  "pman cd <project Name> can be used to change the current working directory to that of the project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cd called")
	},
}

func init() {
	rootCmd.AddCommand(cdCmd)
}
