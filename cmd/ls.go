package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/theredditbandit/pman/pkg/db"
	"github.com/theredditbandit/pman/pkg/ui"
	"github.com/theredditbandit/pman/pkg/utils"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all indexed projects along with their status",
	Long: `List all indexed projects along with their status
    Usage : pman ls
    `,
	RunE: func(cmd *cobra.Command, _ []string) error {
		filterFlag, _ := cmd.Flags().GetString("f")
		oldUI, _ := cmd.Flags().GetBool("o")
		data, err := db.GetAllRecords(db.DBName, StatusBucket)
		if err != nil {
			return err
		}
		if filterFlag != "" {
			fmt.Println("Filtering by status : ", filterFlag)
			data = utils.FilterByStatus(data, filterFlag)
		}
		if oldUI {
			return ui.RenderTable(data)
		}
		return ui.RenderInteractiveTable(data)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().String("f", "", "Filter projects by status. Usage : pman ls --f <status>")
	lsCmd.Flags().Bool("o", false, "list projects using the old ui. Usage : pman ls --o")
}
