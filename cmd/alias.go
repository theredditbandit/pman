package cmd

import (
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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Usage: pman alias <a-long-project-name> <alias>")
			return
		}
		pname := args[0]
		alias := args[1]
		_, err := db.GetRecord(pname, StatusBucket)
		if err != nil {
			fmt.Printf("%s project does not exist in db", pname)
			return
		}
		fmt.Printf("Aliasing %s to %s \n", pname, alias)
		data := map[string]string{alias: pname}
		revData := map[string]string{pname: alias}
		db.WriteToDB(data, ProjectAliasBucket)
		db.WriteToDB(revData, ProjectAliasBucket)
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
