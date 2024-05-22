package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

const (
	StatusBucket       = "projects"
	ProjectPathBucket  = "projectPaths"
	ProjectAliasBucket = "projectAliases"
	ConfigBucket       = "config"
	version            = "1.0"
)

var (
	ErrNoArgs = errors.New("this command has no argument")
)

var rootCmd = &cobra.Command{
	Use:          "pman",
	Short:        "A cli project manager",
	Version:      version,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			return ErrNoArgs
		}
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
