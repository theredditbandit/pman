package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/theredditbandit/pman/pkg/db"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the status of a project",
	Long: `Set the status of a project to a specified status
    Usage:
    pman set <project_name> <status>

    Common statuses: Indexed (default), Idea, Started, Paused, Completed, Aborted, Ongoing, Not Started
    `,
	RunE: func(cmd *cobra.Command, args []string) error {
		interactiveFlag, _ := cmd.Flags().GetBool("i") // TODO: Implement this
		if interactiveFlag {
			cmd.SilenceUsage = true
			return errors.New("Not implemented yet")
		}
		if len(args) != 2 {
			return errors.New("Please provide a directory name")
		}
		var pname string
		alias := args[0]
		status := args[1]
		project, err := db.GetRecord(alias, ProjectAliasBucket)
		if err == nil {
			pname = project
		} else {
			pname = alias
		}
		err = db.UpdateRec(pname, status, StatusBucket)
		if err != nil {
			return fmt.Errorf("Error updating record : %w", err)
		}
		fmt.Printf("Project %s set to status %s\n", pname, status)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().Bool("i", false, "Set the status of projects interactively")
}
