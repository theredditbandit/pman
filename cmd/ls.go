package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"pman/pkg"
	"pman/pkg/db"
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
			pkg.PrintData(data)
			return
		}
		pkg.PrintData(data)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().String("f", "", "Filter projects by status. Usage : pman ls -f <status>")
}
