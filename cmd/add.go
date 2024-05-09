package cmd

import (
	"github.com/spf13/cobra"

	"github.com/theredditbandit/pman/pkg"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a project directory to the index",
	Long: `This command will add a directory to the project database.
    The directory will not be added if it does not contain a README.md.
    `,
	RunE: func(cmd *cobra.Command, args []string) error {
		return pkg.InitDirs(args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
