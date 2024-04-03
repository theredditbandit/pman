package ui

import (
	"fmt"
	"pman/pkg"
	"pman/pkg/db"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// RenderTable: renders the given data as a pretty table
func RenderTable(data map[string]string) error { // FIX : see if you can always print the data in the same order
	var rows [][]string
	for k, v := range data {
		alias, err := db.GetRecord(k, pkg.ProjectAliasBucket)
		if err == nil {
			pname := fmt.Sprintf("%s (%s)", k, alias)
			row := []string{pkg.Title(v), pname} // Status | prjectName (alias)
			rows = append(rows, row)
		} else {
			row := []string{pkg.Title(v), k} // Status | prjectName
			rows = append(rows, row)
		}
	}
	CompletedStyle := lipgloss.NewStyle().Background(lipgloss.Color("#cb85c"))
	defaultStyle := lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Foreground(lipgloss.Color("#FFFFFF"))
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		Headers("Status", "Project ").Width(60).
		Rows(rows...).StyleFunc(func(row, col int) lipgloss.Style {
		if row < len(rows) && row != 0 {
			status := rows[row][col]
			switch {
			case status == "Completed":
				return CompletedStyle
			}
		}
		return defaultStyle
	})
	fmt.Println(t)

	return nil
}
