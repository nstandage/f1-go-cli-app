package view

import (
	"time"
)

//Contains event name/type lap count, fastest lap, whether it's historical or live, date of the event.

type TopBarData struct {
	EventName        string
	EventType        string
	CurrentLap       int
	TotalLaps        int
	FastestLapTime   time.Time
	FastestLapDriver int
	FastestLapNumber int
	StreamType       string
	EventDate        time.Time
}

func TopBar(d *TopBarData) string {
	// var style = lipgloss.NewStyle().
	// 	Bold(true).
	// 	Foreground(lipgloss.Color("#FAFAFA")).
	// 	Background(lipgloss.Color("#7D56F4")).
	// 	PaddingTop(2).
	// 	PaddingLeft(4).
	// 	Width(22)
	return d.EventName
}
