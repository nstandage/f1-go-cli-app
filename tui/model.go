package tui

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/nstandage/f1-go-cli-app/aggregator"
	"github.com/nstandage/f1-go-cli-app/tui/view"
)

type Model struct {
	Window   Window
	Engine   *aggregator.Engine
	offset   uint
	isPaused bool
}

type Window struct {
	width  int
	height int
}

type TickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Window.width = msg.Width
		m.Window.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case TickMsg:
		return m, tick()
	default:
		// log.Printf("Unable to read msg of type: %T", msg)
	}
	return m, nil
}

func (m Model) View() tea.View {
	snapshot := m.Engine.GetSnapshot(0)
	sessionBar := view.SessionBar(snapshot.SessionBar)
	legendBar := view.LegendBar()
	positionColumn := view.PositionsColumn()
	lapSectors := [][]int{
		{2049, 2049, 2049, 2051, 2049, 2051, 2049, 2049},
		{2049, 2049, 2049, 2049, 2049, 2049, 2049, 2049},
		{2048, 2048, 2048, 2048, 2048, 2064, 2064, 2064},
	}
	lapSectorCount := []int{
		len(lapSectors[0]),
		len(lapSectors[1]),
		len(lapSectors[2]),
	}
	topBar := view.Topbar(lapSectorCount)

	driverNames := snapshot.DriverNames
	intervals := []string{
		"----", "0.23", "0.85", "1.04", "3.22", "0.98", "0.12", "1.01", "+1 Lap", "+1 Lap",
	}

	gapToLeaders := []string{
		"----", "0.23", "1.85", "2.04", "3.22", "4.98", "5.12", "6.01", "26.79", "1.23.54",
	}

	lastLap := snapshot.LastLap

	pits := []string{
		"1", "1", "1", "1", "0", "0", "2", "1", "0", "4",
	}

	tires := []string{
		"MEDIUM", "HARD", "SOFT", "MEDIUM", "MEDIUM", "SOFT", "SOFT", "INT", "WET", "SOFT",
	}

	tireAge := []string{
		"23", "22", "10", "17", "0", "1", "30", "29", "1", "2",
	}

	raceControlMessages := snapshot.RaceControlMsgs

	pitStops := []float64{
		3.0, 3.2, 3.8, 2.99, 3.12,
	}

	driverColumn := view.DefaultColumn(driverNames)
	intervalColumn := view.DefaultColumn(intervals)
	gapToLeaderColumn := view.DefaultColumn(gapToLeaders)
	lastLapColumn := view.LastLapColumn(lastLap)
	pitColumn := view.PitColumn(pits)
	tiresColumn := view.TireColumn(tires)
	tireAgeColumn := view.TireAgeColumn(tireAge)
	laps := view.Laps(lapSectors)
	raceControl := view.RaceControl(raceControlMessages)
	pitStopView := view.PitStops(pitStops)

	driverView := lipgloss.JoinHorizontal(
		lipgloss.Top,
		positionColumn,
		driverColumn,
		intervalColumn,
		gapToLeaderColumn,
		lastLapColumn,
		pitColumn,
		tiresColumn,
		tireAgeColumn,
		laps,
	)

	topView := lipgloss.JoinVertical(
		lipgloss.Top,
		sessionBar,
		topBar,
	)

	componentHeight := lipgloss.Height(driverView) +
		lipgloss.Height(topView) +
		lipgloss.Height(legendBar)

	var spacerSize uint = uint(m.Window.height - componentHeight)

	combined := lipgloss.JoinVertical(
		lipgloss.Top,
		topView,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.JoinVertical(
				lipgloss.Center,
				driverView,
				view.Spacer(spacerSize),
				legendBar,
			),
			lipgloss.JoinVertical(
				lipgloss.Top,
				raceControl,
				pitStopView,
			),
			// legendBar,
		),
	)
	v := tea.NewView(combined)
	v.AltScreen = true
	return v
}
