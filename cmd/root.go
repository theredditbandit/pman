package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/theredditbandit/pman/pkg/utils"
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
	Version:      utils.GetVersion(),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			err := cmd.Help()
			if err != nil {
				return errors.Join(err, ErrNoArgs)
			}
			return ErrNoArgs
		}
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
