package ui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/theredditbandit/pman/pkg"
	"github.com/theredditbandit/pman/pkg/db"
	pgr "github.com/theredditbandit/pman/pkg/ui/pager"
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
			if strings.Contains(project, ")") { // project is of the form a-long-project-name (alias)
				projectAliasArr := strings.Split(project, " ")
				project = projectAliasArr[0]
			}
			err := pgr.LaunchRenderer(project)
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

func RenderInteractiveTable(data map[string]string, refreshLastEditedTime bool) error {
	var rows []table.Row
	var lastEdited string
	var timestamp int64

	col := []table.Column{
		{Title: "Status", Width: 20},
		{Title: "Project", Width: 40},
		{Title: "Last Edited", Width: 20},
	}

	if refreshLastEditedTime {
		err := utils.UpdateLastEditedTime()
		if err != nil {
			return err
		}
	} else {
		rec, err := db.GetRecord(db.DBName, "lastRefreshTime", pkg.ConfigBucket)
		if err != nil { // lastRefreshTime key does not exist in db
			refreshLastEditedTime = true
			err := utils.UpdateLastEditedTime()
			if err != nil {
				return err
			}
		}
		if utils.DayPassed(rec) { // lastEdited values are more than a day old. need to refresh them
			refreshLastEditedTime = true
			err := utils.UpdateLastEditedTime()
			if err != nil {
				return err
			}
		}
	}

	for proj, status := range data {
		alias, err := db.GetRecord(db.DBName, proj, pkg.ProjectAliasBucket)
		if refreshLastEditedTime {
			lastEdited, timestamp = utils.GetLastModifiedTime(db.DBName, proj)
			rec := map[string]string{proj: fmt.Sprintf("%s-%d", lastEdited, timestamp)}
			err := db.WriteToDB(db.DBName, rec, pkg.LastUpdatedBucket)
			if err != nil {
				return err
			}
		} else {
			lE, err := db.GetRecord(db.DBName, proj, pkg.LastUpdatedBucket)
			if err != nil {
				return err
			}
			out := strings.Split(lE, "-")
			lastEdited = out[0]
			tstmp, err := strconv.ParseInt(out[1], 10, 64)
			if err != nil {
				return err
			}
			timestamp = tstmp
		}
		if err == nil {
			pname := fmt.Sprintf("%s (%s)", proj, alias)
			row := []string{utils.TitleCase(status), pname, lastEdited, fmt.Sprint(timestamp)} // Status | projectName (alias) | lastEdited | timestamp
			rows = append(rows, row)
		} else {
			row := []string{utils.TitleCase(status), proj, lastEdited, fmt.Sprint(timestamp)} // Status | projectName | lastEdited | timestamp
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
		valI, _ := strconv.ParseInt(rows[i][3], 10, 64)
		valJ, _ := strconv.ParseInt(rows[j][3], 10, 64)
		return valI > valJ
	})
	cleanUp := func(r []table.Row) []table.Row {
		result := make([]table.Row, len(r))
		for i, inner := range r {
			n := len(inner)
			result[i] = inner[:n-1]
		}
		return result
	}
	rows = cleanUp(rows)
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
