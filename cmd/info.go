package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/theredditbandit/pman/pkg/utils"
)

var (
	ErrBadUsageInfoCmd error = errors.New("bad usage of info command")
)

var infoCmd = &cobra.Command{
	Use:     "info",
	Short:   "The info command pretty prints the README.md file present at the root of the specified project.",
	Aliases: []string{"ifo", "ifno", "ino"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			fmt.Println("Please provide a project name")
			return ErrBadUsageInfoCmd
		}
		projectName := args[0]
		infoData, err := utils.ReadREADME(projectName)
		if err != nil {
			return err
		}
		md, err := utils.BeautifyMD(infoData)
		if err != nil {
			return err
		}
		fmt.Print(md)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
