package cmd

import (
	"github.com/spf13/cobra"

	"github.com/theredditbandit/pman/pkg"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Takes exactly 1 argument, a directory name, and initializes it as a project directory.",
	Long: `This command will initialize a directory as a project directory.

    It will index any folder which contains a README.md as a project directory.

    Running pman init <dirname> is the same as running: pman add <dirname>/*
    `,
	Run: func(cmd *cobra.Command, args []string) {
		pkg.InitDirs(args)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
