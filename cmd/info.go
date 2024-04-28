package cmd

import (
	"fmt"
	"github.com/theredditbandit/pman/pkg/db"
	"os"
	"path/filepath"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "The info command pretty prints the README.md file present at the root of the specified project.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide a project name")
			return
		}
		projectName := args[0]
		actualName, err := db.GetRecord(projectName, ProjectAliasBucket)
		if err == nil {
			projectName = actualName
		}
		path, err := db.GetRecord(projectName, ProjectPathBucket)
		if err != nil {
			fmt.Printf("project: %v not a valid project\n", projectName)
			return
		}
		infoPath := filepath.Join(path, "README.md")
		infoData, err := os.ReadFile(infoPath)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		r, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(120),
			glamour.WithAutoStyle(),
		)
		out, _ := r.Render(string(infoData))
		fmt.Print(out)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
