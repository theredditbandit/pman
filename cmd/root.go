package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	StatusBucket       = "projects"
	ProjectPathBucket  = "projectPaths"
	ProjectAliasBucket = "projectAliases"
	version            = "1.0"
)

var rootCmd = &cobra.Command{
	Use:     "pman",
	Short:   "A cli project manager",
	Version: version,
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
