package cmd

import (
	"fmt"
	"log"
	"pman/pkg"
	"pman/pkg/db"
	"pman/pkg/ui"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all indexed projects along with their status",
	Long: `List all indexed projects along with their status
    Usage : pman ls
    `,
	Run: func(cmd *cobra.Command, args []string) {
		filterFlag, _ := cmd.Flags().GetString("f")
		data, err := db.GetAllRecords(StatusBucket)
		if err != nil {
			log.Fatal(err)
		}
		if filterFlag != "" {
			fmt.Println("Filtering by status : ", filterFlag)
			data := pkg.FilterByStatus(data, filterFlag)
			_ = data
			ui.RenderTable(data)
			// ui.Pikachu()
			return
		}
		ui.RenderTable(data)
		// ui.Pikachu()

	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().String("f", "", "Filter projects by status. Usage : pman ls -f <status>")
}
