package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/theredditbandit/pman/pkg/db"
)

var delCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Deletes a project from the index database. This does not delete the project from the filesystem",
	Aliases: []string{"del", "d"},
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Please provide a project name")
		}
		projName := args[0]
		_, err := db.GetRecord(projName, StatusBucket)
		if err != nil {
			fmt.Printf("%s is not a valid project to be deleted\n", projName)
			fmt.Println("Note : projects cannot be deleted using their alias")

			// not a real error here
			return nil
		}
		err = db.DeleteFromDb(projName, ProjectPathBucket)
		if err != nil {
			return err
		}
		err = db.DeleteFromDb(projName, StatusBucket)
		if err != nil {
			return err
		}
		alias, err := db.GetRecord(projName, ProjectAliasBucket)
		if err == nil {
			err = db.DeleteFromDb(alias, ProjectAliasBucket)
			if err != nil {
				return err
			}
			err = db.DeleteFromDb(projName, ProjectAliasBucket)
			if err != nil {
				return err
			}
		}
		err = nil
		if err != nil {
			return err
		}
		fmt.Printf("Successfully deleted %s from the db \n", projName)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(delCmd)
}
