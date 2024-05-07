package ui

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/theredditbandit/pman/pkg"
	"github.com/theredditbandit/pman/pkg/db"
	"github.com/theredditbandit/pman/pkg/utils"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type tableModel struct {
	table table.Model
}

func (m tableModel) Init() tea.Cmd { return nil }

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func RenderInteractiveTable(data map[string]string) error {
	col := []table.Column{
		{Title: "Status", Width: 20},
		{Title: "Project", Width: 40},
		{Title: "Last Edited", Width: 20},
	}
	var rows []table.Row
	for p, status := range data {
		alias, err := db.GetRecord(p, pkg.ProjectAliasBucket)
		lastEdited := utils.GetLastModifiedTime(p)
		if err == nil {
			pname := fmt.Sprintf("%s (%s)", p, alias)
			row := []string{utils.TitleCase(status), pname, lastEdited} // Status | projectName (alias) | lastEdited
			rows = append(rows, row)
		} else {
			row := []string{utils.TitleCase(status), p, lastEdited} // Status | projectName | lastEdited
			rows = append(rows, row)
		}
	}

	if len(rows) == 0 {
		fmt.Printf("No projects found in the database\n\n")
		fmt.Printf("Add projects to the database using \n\n")
		fmt.Println("pman init .")
		fmt.Println("or")
		fmt.Println("pman add <projectDir>")
		return nil
	}
	sort.Slice(rows, func(i, j int) bool {
		valI := rows[i][1]
		valJ := rows[j][1]
		return valI < valJ
	})
	t := table.New(
		table.WithColumns(col),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
		table.WithWidth(90),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := tableModel{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		return fmt.Errorf("Error running program: %s", err)
	}
	return nil
}
