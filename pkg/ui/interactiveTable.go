package ui

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/theredditbandit/pman/pkg"
	"github.com/theredditbandit/pman/pkg/db"
	p "github.com/theredditbandit/pman/pkg/ui/pager"
	"github.com/theredditbandit/pman/pkg/utils"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type tableModel struct {
	table table.Model
}

func (tableModel) Init() tea.Cmd { return nil }

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if msg, ok := msg.(tea.KeyMsg); ok {
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
			project := m.table.SelectedRow()[1]
			err := p.LaunchRenderer(project)
			if err != nil {
				return m, tea.Quit
			}
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
	for proj, status := range data {
		alias, err := db.GetRecord(db.DBName, proj, pkg.ProjectAliasBucket)
		lastEdited := utils.GetLastModifiedTime(proj)
		if err == nil {
			pname := fmt.Sprintf("%s (%s)", proj, alias)
			row := []string{utils.TitleCase(status), pname, lastEdited} // Status | projectName (alias) | lastEdited
			rows = append(rows, row)
		} else {
			row := []string{utils.TitleCase(status), proj, lastEdited} // Status | projectName | lastEdited
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
		return fmt.Errorf("error running program: %w", err)
	}
	return nil
}
