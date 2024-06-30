package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	c "github.com/theredditbandit/pman/constants"
	"github.com/theredditbandit/pman/pkg/db"
	"github.com/theredditbandit/pman/pkg/ui"
	"github.com/theredditbandit/pman/pkg/utils"
)

// iCmd represents the interactive command
var iCmd = &cobra.Command{
	Use:     "i",
	Short:   "Launches pman in interactive mode",
	Aliases: []string{"interactive", "iteractive"},
	RunE: func(cmd *cobra.Command, _ []string) error {
		filterFlag, _ := cmd.Flags().GetString("f")
		refreshLastEditTime, _ := cmd.Flags().GetBool("r")
		data, err := db.GetAllRecords(db.DBName, c.StatusBucket)
		if err != nil {
			return err
		}
		if filterFlag != "" {
			fmt.Println("Filtering by status : ", filterFlag)
			data = utils.FilterByStatuses(data, strings.Split(filterFlag, ","))
		}
		return ui.RenderInteractiveTable(data, refreshLastEditTime)
	},
}

func init() {
	rootCmd.AddCommand(iCmd)
	iCmd.Flags().String("f", "", "Filter projects by status. Usage : pman ls --f <status1[,status2]>")
	iCmd.Flags().Bool("r", false, "Refresh Last Edited time: pman ls --r")
}
