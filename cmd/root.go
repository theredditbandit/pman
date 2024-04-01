package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const StatusBucket = "projects"
const ProjectPathBucket = "projectPaths"

var rootCmd = &cobra.Command{
	Use:   "pman",
	Short: "The final project manager",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
