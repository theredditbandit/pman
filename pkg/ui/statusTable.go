package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"os"
	"pman/pkg"
	"pman/pkg/db"
	"sort"
)

// RenderTable: renders the given data as a table
func RenderTable(data map[string]string) error {
	var TableData [][]string
	for p, status := range data {
		alias, err := db.GetRecord(p, pkg.ProjectAliasBucket)
		if err == nil {
			pname := fmt.Sprintf("%s (%s)", p, alias)
			row := []string{pkg.TitleCase(status), pname} // Status | prjectName (alias)
			TableData = append(TableData, row)
		} else {
			row := []string{pkg.TitleCase(status), p} // Status | prjectName
			TableData = append(TableData, row)
		}
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
	headers := []string{"Project Name", "Status"}
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
		Headers(headers...).
		Width(80).
		Rows(TableData...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}
			color, ok := statusColors[fmt.Sprint(TableData[row-1][0])]
			if ok {
				return baseStyle.Copy().Foreground(color)
			} else {
				color := statusColors["Default"]
				return baseStyle.Copy().Foreground(color)
			}
		})
	fmt.Println(t)
	return nil
}
