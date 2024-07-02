package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/theredditbandit/pman/pkg/db"
	"github.com/theredditbandit/pman/pkg/utils"
)

var ErrPathDoesNotExist = errors.New("Project path does not exist")

var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "Quicky change directory to the project name",
	Long:  "pman cd <project Name> can be used to change the current working directory to that of the project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			fmt.Println("Please provide a project name")
			return ErrBadUsageInfoCmd
		}
		projectName := args[0]
		pPath, err := utils.GetProjectPath(db.DBName, projectName)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("path", pPath, "does not exist")
				return ErrPathDoesNotExist
			}
			return err
		}
		cddir := filepath.Dir(pPath)
		fmt.Printf("cddir: %v\n", cddir)
		err = os.Chdir(cddir) // this approach does
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("this approach does not work because when pman is executed in a shell it is started in a new process and any directories changed withtin that process will not affect the directory of the process that was used to call pman cd in the first place, the only way to do this right is with a shell wrapper but i am not looking to add that to pman at the moment to it remains unfinished ")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cdCmd)
}
