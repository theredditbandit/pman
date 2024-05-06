package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/theredditbandit/pman/pkg/db"
	"github.com/theredditbandit/pman/pkg/utils"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of a project",
	Long:  `Query the database for the status of a project.`,
	RunE: func(_ *cobra.Command, args []string) error {
		var alias string
		if len(args) != 1 {
			return errors.New("Please provide a project name")
		}
		projName := args[0]
		actualName, err := db.GetRecord(projName, ProjectAliasBucket)
		if err == nil { // check if user has supplied an alias instead of actual project name
			alias = projName
			projName = actualName
		}
		status, err := db.GetRecord(projName, StatusBucket)
		if err != nil {
			return fmt.Errorf("%s is not a valid project name : Err -> %w", projName, err)
		}
		if alias != "" {
			fmt.Printf("Status of %s (%s) : %s\n", projName, alias, utils.TitleCase(status))
		} else {
			fmt.Printf("Status of %s  : %s\n", projName, utils.TitleCase(status))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
