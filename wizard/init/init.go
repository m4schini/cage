package init

import (
	"fmt"
	"os"
	"sync"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type styles struct {
	app           lipgloss.Style
	title         lipgloss.Style
	statusMessage lipgloss.Style
}

func newStyles(darkBG bool) styles {
	lightDark := lipgloss.LightDark(darkBG)

	return styles{
		app: lipgloss.NewStyle().
			Padding(1, 2),
		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1),
		statusMessage: lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("#04B575"), lipgloss.Color("#04B575"))),
	}
}

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type model struct {
	styles        styles
	darkBG        bool
	width, height int
	once          *sync.Once
	delegateKeys  *delegateKeyMap
	quitting      bool

	agentList     list.Model
	selectedAgent string

	baseList     list.Model
	selectedBase string
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		tea.RequestBackgroundColor,
	)
}

func (m *model) updateListProperties() {
	// Update agentList size.
	h, v := m.styles.app.GetFrameSize()
	m.agentList.SetSize(m.width-h, m.height-v)

	// Update the model and agentList styles.
	m.styles = newStyles(m.darkBG)
	m.agentList.Styles.Title = m.styles.title
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.darkBG = msg.IsDark()
		m.updateListProperties()
		return m, nil

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateListProperties()
		return m, nil
	}

	switch msg.(type) {
	case tea.KeyPressMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.agentList.FilterState() == list.Filtering {
			break
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.agentList.Update(msg)
	m.agentList = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) View() tea.View {
	if m.quitting {
		return tea.NewView("\n  See you later!\n\n")
	}

	if m.selectedAgent == "" {
		v := tea.NewView(m.styles.app.Render(m.agentList.View()))
		v.AltScreen = true
		return v
	} else {
		tea.NewView(m.selectedAgent)
	}
	if m.selectedBase == "" {
		v := tea.NewView(m.styles.app.Render(m.baseList.View()))
		v.AltScreen = true
		return v
	}

	return tea.NewView(fmt.Sprintf("agent=%v base=%v", m.selectedAgent, m.selectedBase))
}

func initialModel() *model {
	// Initialize the model and agentList.
	delegateKeys := newDelegateKeyMap()
	m := &model{}
	m.styles = newStyles(false) // default to dark background styles
	m.delegateKeys = delegateKeys

	m.agentList = NewAgentList(m, &m.styles, delegateKeys)
	m.baseList = NewBaseList(m, &m.styles, delegateKeys)

	return m
}

func Run() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
