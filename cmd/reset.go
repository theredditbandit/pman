package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/theredditbandit/pman/pkg/db"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Deletes the current indexed projects, run pman init to reindex the projects",
	RunE: func(_ *cobra.Command, _ []string) error {
		err := db.DeleteDb()
		if err != nil {
			return err
		}

		fmt.Println("Successfully reset the database, run pman init to reindex the projects")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
