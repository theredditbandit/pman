package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/thoas/go-funk"

	"github.com/theredditbandit/pman/pkg"
	"github.com/theredditbandit/pman/pkg/db"
)

var (
	appStyle   = lipgloss.NewStyle().Padding(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	name   string
	status string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return i.status }
func (i item) FilterValue() string { return i.status }

type listKeyMap struct {
	setStatus        key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		setStatus: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "set status"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type model struct {
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
}

func newModel() (model, error) {
	var (
		listKeys     = newListKeyMap()
		delegateKeys = newDelegateKeyMap()
	)

	data, err := db.GetAllRecords(pkg.StatusBucket)
	if err != nil {
		return model{}, err
	}

	d := funk.Map(data, func(k, v string) item {
		return item{
			name:   k,
			status: v,
		}
	})

	formattedData, ok := d.([]item)
	if !ok {
		return model{}, fmt.Errorf("something went wrong while creating a list of items")
	}
	items := make([]list.Item, len(data))
	for i := 0; i < len(data); i++ {
		items[i] = formattedData[i]
	}

	delegate := newItemDelegate(delegateKeys)
	projectList := list.New(items, delegate, 0, 0)
	projectList.Title = "pman"
	projectList.Styles.Title = titleStyle

	projectList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.setStatus,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}
	return model{
		list:         projectList,
		keys:         listKeys,
		delegateKeys: delegateKeys,
	}, nil
}

func (model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.keys.setStatus):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.keys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.keys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.keys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.keys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
}

func Tui() error {
	model, err := newModel()
	if err != nil {
		return err
	}

	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		return fmt.Errorf("Error running program: %w", err)
	}

	return nil
}
