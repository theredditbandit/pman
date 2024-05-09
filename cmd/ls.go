package cmd

import (
	"fmt"
	"log"

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
	Run: func(cmd *cobra.Command, args []string) {
		filterFlag, _ := cmd.Flags().GetString("f")
		data, err := db.GetAllRecords(db.DBName, StatusBucket)
		if err != nil {
			log.Fatal(err)
		}
		if filterFlag != "" {
			fmt.Println("Filtering by status : ", filterFlag)
			data := utils.FilterByStatus(data, filterFlag)
			ui.RenderTable(data)
			return
		}
		ui.RenderTable(data)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().String("f", "", "Filter projects by status. Usage : pman ls -f <status>")
}
