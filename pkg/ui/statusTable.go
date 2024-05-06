package ui

import (
	"fmt"
	"os"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/theredditbandit/pman/pkg"
	"github.com/theredditbandit/pman/pkg/db"
	"github.com/theredditbandit/pman/pkg/utils"
)

// RenderTable: renders the given data as a table
func RenderTable(data map[string]string) error {
	var TableData [][]string

	for p, status := range data {
		alias, err := db.GetRecord(p, pkg.ProjectAliasBucket)
		lastEdited := utils.GetLastModifiedTime(p)
		if err == nil {
			pname := fmt.Sprintf("%s (%s)", p, alias)
			row := []string{utils.TitleCase(status), pname, lastEdited} // Status | projectName (alias) | lastEdited
			TableData = append(TableData, row)
		} else {
			row := []string{utils.TitleCase(status), p, lastEdited} // Status | projectName | lastEdited
			TableData = append(TableData, row)
		}
	}
	if len(TableData) == 0 {
		fmt.Printf("No projects found in the database\n\n")
		fmt.Printf("Add projects to the database using \n\n")
		fmt.Println("pman init .")
		fmt.Println("or")
		fmt.Println("pman add <projectDir>")
		return nil
	}
	sort.Slice(TableData, func(i, j int) bool {
		valI := TableData[i][1]
		valJ := TableData[j][1]
		return valI < valJ
	})
	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Copy().Foreground(lipgloss.Color("252")).Bold(true)
	// selectedStyle := baseStyle.Copy().Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))
	statusColors := map[string]lipgloss.Color{
		"Idea":        lipgloss.Color("#FF87D7"),
		"Indexed":     lipgloss.Color("#727272"),
		"Not Started": lipgloss.Color("#D7FF87"),
		"Ongoing":     lipgloss.Color("#00E2C7"),
		"Started":     lipgloss.Color("#00E2C7"),
		"Paused":      lipgloss.Color("#7D5AFC"),
		"Completed":   lipgloss.Color("#75FBAB"),
		"Aborted":     lipgloss.Color("#FF875F"),
		"Default":     lipgloss.Color("#929292"),
	}
	headers := []string{"Status", "Project Name", "Last Edited"}
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(headers...).
		Width(90).
		Rows(TableData...).
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}
			color, ok := statusColors[fmt.Sprint(TableData[row-1][0])]
			if ok {
				return baseStyle.Copy().Foreground(color)
			}

			color = statusColors["Default"]
			return baseStyle.Copy().Foreground(color)
		})
	fmt.Println(t)
	return nil
}
