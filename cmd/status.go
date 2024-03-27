package cmd

import (
	"fmt"
	"pman/pkg"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of a project",
	Long:  `Query the database for the status of a project.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide a project name")
			return
		}
		projName := args[0]
		status, err := pkg.GetRecord(projName, StatusBucket)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Status of %s : %s\n", projName, pkg.Title(status))
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
