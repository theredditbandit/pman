package cmd

import (
	"fmt"
	"os"
	"projman/pkg"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Takes exactly 1 argument, a directory name, and initializes it as a project directory.",
	Long: `This command will initialize a directory as a project directory.

    It will index any folder which contains a README.md as a project directory.

    Running projman init <dirname> is the same as running: projman add <dirname>/*
    `,
	Run: func(cmd *cobra.Command, args []string) {
		projIdentifier := "README.md" // the file which identifies a project directory
		if len(args) != 1 {
			fmt.Println("Please provide a directory name")
			return
		}
		dirname := args[0]
		if stat, err := os.Stat(dirname); os.IsNotExist(err) {
			fmt.Printf("%s is not a directory \n", dirname)
			return
		} else if !stat.IsDir() {
			fmt.Printf("%s is a file and not a directory \n", dirname)
			return
		}
		fmt.Printf("Indexing %s. . .\n", dirname)
		projDirs, err := pkg.IndexDir(dirname, projIdentifier)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("projDirs: %v\n", projDirs)
        for k, v := range projDirs {
            fmt.Println("Project Name: ", k)
            fmt.Println("Project Path: ", v)
            fmt.Println("=====================================")
        }

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
