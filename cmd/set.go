package cmd

import (
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
	Run: func(cmd *cobra.Command, args []string) {
		interactiveFlag, _ := cmd.Flags().GetBool("i") // TODO: Implement this
		if interactiveFlag {
			fmt.Println("Not implemented yet")
			return
		}
		if len(args) != 2 {
			fmt.Println("Please provide a directory name")
			return
		}
		var pname string
		alias := args[0]
		status := args[1]
		project, err := db.GetRecord(db.DBName, alias, ProjectAliasBucket)
		if err == nil {
			pname = project
		} else {
			pname = alias
		}
		err = db.UpdateRec(db.DBName, pname, status, StatusBucket)
		if err != nil {
			fmt.Println("Error updating record : ", err)
			return
		}
		fmt.Printf("Project %s set to status %s\n", pname, status)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().Bool("i", false, "Set the status of projects interactively")
}
