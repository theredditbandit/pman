/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "path to your project directory",
	Long: `This command will initialize a directory as a project directory.
    It will index any folder which contains a README.md as a project directory.
    
    Running projman init <dirname> is the same as running: projman add <dirname>/*
    `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
