package tui

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/nstandage/f1-go-cli-app/model"
	"github.com/nstandage/f1-go-cli-app/tui/view"
)

type Model struct {
	width  int
	height int
	text   string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case *model.Interval:
		m.text = msg.DateStart.String()
		return m, nil
	default:
		log.Printf("Unable to read msg of type: %T", msg)
	}
	return m, nil
}

func (m Model) View() tea.View {
	var content string

	if m.width == 0 {
		content = "Loading..."
	} else {
		content = lipgloss.Place(m.width,
			m.height-6,
			lipgloss.Center,
			lipgloss.Center,
			m.text,
		)
	}
	barData := view.GetTestSessionBarData()
	topBar := view.SessionBar(&barData)
	legendBar := view.LegendBar()
	combined := lipgloss.JoinVertical(lipgloss.Top, topBar, content, legendBar)
	v := tea.NewView(combined)
	v.AltScreen = true
	return v
}
