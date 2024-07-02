package cmd

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	c "github.com/theredditbandit/pman/constants"
	"github.com/theredditbandit/pman/pkg/db"
)

var (
	ErrBadUsageDelCmd = errors.New("bad usage of delete command")
)

var delCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Deletes a project from the index database. This does not delete the project from the filesystem",
	Aliases: []string{"del", "d"},
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			fmt.Println("Usage : pman delete <projectName>")
			return ErrBadUsageDelCmd
		}
		projName := args[0]
		_, err := db.GetRecord(db.DBName, projName, c.StatusBucket)
		if err != nil {
			fmt.Printf("%s is not a valid project to be deleted\n", projName)
			fmt.Println("Note : projects cannot be deleted using their alias")
			return err
		}
		err = db.DeleteFromDb(db.DBName, projName, c.ProjectPaths)
		if err != nil {
			return err
		}
		err = db.DeleteFromDb(db.DBName, projName, c.StatusBucket)
		if err != nil {
			return err
		}
		err = db.DeleteFromDb(db.DBName, projName, c.LastUpdatedBucket)
		if err != nil {
			return err
		}
		alias, err := db.GetRecord(db.DBName, projName, c.ProjectAliasBucket)
		if err == nil {
			err = db.DeleteFromDb(db.DBName, alias, c.ProjectAliasBucket)
			if err != nil {
				return err
			}
			err = db.DeleteFromDb(db.DBName, projName, c.ProjectAliasBucket)
			if err != nil {
				return err
			}
		}
		err = nil
		if err != nil {
			return err
		}
		lastEdit := make(map[string]string)
		lastEdit["lastWrite"] = fmt.Sprint(time.Now().Format("02 Jan 06 15:04"))
		err = db.WriteToDB(db.DBName, lastEdit, c.ConfigBucket)
		if err != nil {
			log.Print(err)
			return err
		}
		fmt.Printf("Successfully deleted %s from the db \n", projName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(delCmd)
}
