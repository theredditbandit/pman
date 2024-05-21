package cmd

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/theredditbandit/pman/pkg/db"
)

var (
	ErrBadUsageAliasCmd = errors.New("bad usage of alias command")
)

// aliasCmd represents the alias command
var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Sets the alias for a project, whose name might be too big",
	Long: `The idea is instead of having to type a-very-long-project-name-every-time you can alias it to
avlpn or something smaller and use that to query pman`,
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) != 2 {
			fmt.Println("Usage: pman alias <a-long-project-name> <alias>")
			return ErrBadUsageAliasCmd
		}
		pname := args[0]
		alias := args[1]
		_, err := db.GetRecord(db.DBName, pname, StatusBucket)
		if err != nil {
			fmt.Printf("%s project does not exist in db", pname)
			return err
		}
		fmt.Printf("Aliasing %s to %s \n", pname, alias)
		data := map[string]string{alias: pname}
		revData := map[string]string{pname: alias}
		err = db.WriteToDB(db.DBName, data, ProjectAliasBucket)
		if err != nil {
			return err
		}
		err = db.WriteToDB(db.DBName, revData, ProjectAliasBucket)
		if err != nil {
			return err
		}
		lastEdit := make(map[string]string)
		lastEdit["lastWrite"] = fmt.Sprint(time.Now().Format("02 Jan 06 15:04"))
		err = db.WriteToDB(db.DBName, lastEdit, ConfigBucket)
		if err != nil {
			log.Print(err)
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
