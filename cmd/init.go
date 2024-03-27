package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"pman/pkg"

	"github.com/spf13/cobra"
)

const StatusBucket = "projects"
const ProjectPaths = "projectPaths"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Takes exactly 1 argument, a directory name, and initializes it as a project directory.",
	Long: `This command will initialize a directory as a project directory.

    It will index any folder which contains a README.md as a project directory.

    Running pman init <dirname> is the same as running: pman add <dirname>/*
    `,
	Run: func(cmd *cobra.Command, args []string) {
		// the file which identifies a project directory
		projIdentifier := "README.md" // TODO : make this configurable using a flag
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
		projDirs, err := pkg.IndexDir(dirname, projIdentifier)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Indexed %d project directories . . .\n\n", len(projDirs))
		projectStatusMap := make(map[string]string)
		projectPathMap := make(map[string]string)
		for k, v := range projDirs { // k : full project path, v : project status ,
			projectStatusMap[filepath.Base(k)] = v // filepath.Base(k) : project name
			projectPathMap[filepath.Base(k)] = k
		}
		err = pkg.WriteToDB(projectStatusMap, StatusBucket)
		if err != nil {
			log.Fatal(err)
			return
		}
		err = pkg.WriteToDB(projectPathMap, ProjectPaths)
		if err != nil {
			log.Fatal(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
