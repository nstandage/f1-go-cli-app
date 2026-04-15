package tui

import (
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
		// log.Printf("Unable to read msg of type: %T", msg)
	}
	return m, nil
}

func (m Model) View() tea.View {
	// var content string

	// if m.width == 0 {
	// 	content = "Loading..."
	// } else {
	// 	content = lipgloss.Place(m.width,
	// 		m.height-6,
	// 		lipgloss.Center,
	// 		lipgloss.Center,
	// 		m.text,
	// 	)
	// }
	barData := view.GetTestSessionBarData()
	sessionBar := view.SessionBar(&barData)
	legendBar := view.LegendBar()
	positionColumn := view.PositionsColumn()
	var lapSectors = [][]int{
		[]int{2049, 2049, 2049, 2051, 2049, 2051, 2049, 2049},
		[]int{2049, 2049, 2049, 2049, 2049, 2049, 2049, 2049},
		[]int{2048, 2048, 2048, 2048, 2048, 2064, 2064, 2064},
	}
	var lapSectorCount = []int{
		len(lapSectors[0]),
		len(lapSectors[1]),
		len(lapSectors[2]),
	}
	topBar := view.Topbar(lapSectorCount)

	var driverNames = []string{
		"VER", "NOR", "LEC", "PIA", "PER", "HAM", "ANT", "RUS", "HAD", "SAI",
	}

	var intervals = []string{
		"----", "0.23", "0.85", "1.04", "3.22", "0.98", "0.12", "1.01", "+1 Lap", "+1 Lap",
	}

	var gapToLeaders = []string{
		"----", "0.23", "1.85", "2.04", "3.22", "4.98", "5.12", "6.01", "26.79", "1.23.54",
	}

	var lastLap = []string{
		"1.29.54", "1.29.64", "1.30.85", "1.30.04", "1.30.22", "1.31.98", "1.35.12", "1.36.01", "1.46.79", "1.23.54",
	}

	var pits = []string{
		"1", "1", "1", "1", "0", "0", "2", "1", "0", "4",
	}

	var tires = []string{
		"MEDIUM", "HARD", "SOFT", "MEDIUM", "MEDIUM", "SOFT", "SOFT", "INT", "WET", "SOFT",
	}

	var tireAge = []string{
		"23", "22", "10", "17", "0", "1", "30", "29", "1", "2",
	}

	var raceControlMessages = []string{
		"CAR 43 (COL) TIME 1:43.165 DELETED - TRACK LIMITS AT TURN 13 LAP 21 14:47:52",
		"SAFETY CAR DEPLOYED",
		"MEDICAL CAR DEPLOYED",
		"TURN 13 INCIDENT INVOLVING CARS 43 (COL) AND 87 (BEA) NOTED",
	}

	var pitStops = []float64{
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

	combined := lipgloss.JoinVertical(lipgloss.Top, sessionBar, topBar,
		lipgloss.JoinHorizontal(lipgloss.Top,
			positionColumn,
			driverColumn,
			intervalColumn,
			gapToLeaderColumn,
			lastLapColumn,
			pitColumn,
			tiresColumn,
			tireAgeColumn,
			laps,
			lipgloss.JoinVertical(
				lipgloss.Top,
				raceControl,
				pitStopView,
			),
		),
		legendBar,
	)
	v := tea.NewView(combined)
	v.AltScreen = true
	return v
}
