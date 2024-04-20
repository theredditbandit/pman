package cmd

import (
	"fmt"
	"pman/pkg"
	"pman/pkg/db"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of a project",
	Long:  `Query the database for the status of a project.`,
	Run: func(cmd *cobra.Command, args []string) {
		var alias string
		if len(args) != 1 {
			fmt.Println("Please provide a project name")
			return
		}
		projName := args[0]
		actualName, err := db.GetRecord(projName, ProjectAliasBucket)
		if err == nil { // check if user has supplied an alias instead of actual project name
			alias = projName
			projName = actualName
		}
		status, err := db.GetRecord(projName, StatusBucket)
		if err != nil {
			fmt.Printf("%s is not a valid project name : Err -> %s", projName, err)
			return
		}
		if alias != "" {
			fmt.Printf("Status of %s (%s) : %s\n", projName, alias, pkg.TitleCase(status))
		} else {
			fmt.Printf("Status of %s  : %s\n", projName, pkg.TitleCase(status))
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
