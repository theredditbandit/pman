package cmd

import (
	"fmt"
	"log"
	"pman/pkg"

	"github.com/spf13/cobra"
)

// TODO : Pretty print this in the future
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all indexed projects along with their status",
	Long: `List all indexed projects along with their status
    Usage : pman ls
    `,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := pkg.GetAllRecords(StatusBucket)
		if err != nil {
			log.Fatal(err)
		}
		for k, v := range data {
			fmt.Println(pkg.Title(v), " : ", k)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
	// TODO: add -f filter flag to only list projects with a certain status
	// to implement a bucket will be created that maps a status to all the projects with that status
	// will have to refactor the WriteToDB function to accept an intrface as the data parameter
}
