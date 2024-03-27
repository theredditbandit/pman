package cmd

import (
	"fmt"
	"pman/pkg"

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
		if len(args) != 2 {
			fmt.Println("Please provide a directory name")
			return
		}
		project := args[0]
		status := args[1]
		err := pkg.UpdateRec(project, status, StatusBucket)
		if err != nil {
			fmt.Println("Error updating record")
			return
		}
        fmt.Printf("Project %s set to status %s\n", project, status)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
