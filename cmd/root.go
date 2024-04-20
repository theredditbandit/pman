package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const StatusBucket = "projects"
const ProjectPathBucket = "projectPaths"
const ProjectAliasBucket = "projectAliases"

var rootCmd = &cobra.Command{
	Use:   "pman",
	Short: "The final project manager",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		interactiveFlag, _ := cmd.Flags().GetBool("i") // TODO: Implement this
		if interactiveFlag {
			fmt.Println("Not implemented")
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
	rootCmd.Flags().Bool("i", false, "Run pman interactively")
}
