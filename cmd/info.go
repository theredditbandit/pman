package cmd

import (
	"fmt"

	"github.com/theredditbandit/pman/pkg/utils"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:     "info",
	Short:   "The info command pretty prints the README.md file present at the root of the specified project.",
	Aliases: []string{"ifo", "ifno", "ino"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide a project name")
			return
		}
		projectName := args[0]
		infoData := utils.ReadREADME(projectName)
		md := utils.BeautifyMD(infoData)
		fmt.Print(md)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
