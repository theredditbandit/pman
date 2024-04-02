package cmd

import (
	"fmt"
	"pman/pkg/db"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the status of a project",
	Long: `Set the status of a project to a specified status
    Usage:
    pman set <project_name> <status>

    Common statuses: Indexed (default) , Started , Paused , Completed , Aborted , Deleted , Ongoing , Not Started
    `,
	Run: func(cmd *cobra.Command, args []string) {
		interactiveFlag, _ := cmd.Flags().GetBool("i")
		if interactiveFlag {
			fmt.Println("Not implemented yet")
			return
		}
		if len(args) != 2 {
			fmt.Println("Please provide a directory name")
			return
		}
		project := args[0]
		status := args[1]
		err := db.UpdateRec(project, status, StatusBucket)
		if err != nil {
			fmt.Println("Error updating record")
			return
		}
		fmt.Printf("Project %s set to status %s\n", project, status)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().Bool("i", false, "Set the status of projects interactively")
}
