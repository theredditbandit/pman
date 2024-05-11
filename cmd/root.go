package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

const (
	StatusBucket       string = "projects"
	ProjectPathBucket  string = "projectPaths"
	ProjectAliasBucket string = "projectAliases"
	version            string = "1.0"
)

var rootCmd = &cobra.Command{
	Use:     "pman",
	Short:   "A cli project manager",
	Version: version,
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("this command has no argument")
		}
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
