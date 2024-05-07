package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/theredditbandit/pman/pkg/db"
)

// aliasCmd represents the alias command
var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Sets the alias for a project, whose name might be too big",
	Long: `The idea is instead of having to type a-very-long-project-name-every-time you can alias it to
avlpn or something smaller and use that to query pman`,
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("Please provide a project name and an alias")
		}
		pname := args[0]
		alias := args[1]
		_, err := db.GetRecord(pname, StatusBucket)
		if err != nil {
			return fmt.Errorf("%s project does not exist in db", pname)
		}
		fmt.Printf("Aliasing %s to %s \n", pname, alias)
		data := map[string]string{alias: pname}
		revData := map[string]string{pname: alias}
		err = db.WriteToDB(data, ProjectAliasBucket)
		if err != nil {
			return err
		}
		err = db.WriteToDB(revData, ProjectAliasBucket)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
