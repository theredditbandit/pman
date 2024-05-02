package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/theredditbandit/pman/pkg/db"
)

var delCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Deletes a project from the index database. This does not delete the project from the filesystem",
	Aliases: []string{"del", "d"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage : pman delete <projectName>")
			return
		}
		projName := args[0]
		_, err := db.GetRecord(projName, StatusBucket)
		if err != nil {
			fmt.Printf("%s is not a valid project to be deleted\n", projName)
			fmt.Println("Note : projects cannot be deleted using their alias")
			return
		}
		err = db.DeleteFromDb(projName, ProjectPathBucket)
		if err != nil {
			log.Fatal(err)
		}
		err = db.DeleteFromDb(projName, StatusBucket)
		if err != nil {
			log.Fatal(err)
		}
		alias, err := db.GetRecord(projName, ProjectAliasBucket)
		if err == nil {
			err = db.DeleteFromDb(alias, ProjectAliasBucket)
			if err != nil {
				log.Fatal(err)
			}
			err = db.DeleteFromDb(projName, ProjectAliasBucket)
			if err != nil {
				log.Fatal(err)
			}
		}
		err = nil
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Successfully deleted %s from the db \n", projName)

	},
}

func init() {
	rootCmd.AddCommand(delCmd)
}
