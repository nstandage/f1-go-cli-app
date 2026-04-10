package view

import (
	"fmt"
	"time"

	"charm.land/lipgloss/v2"
)

//Contains event name/type lap count, fastest lap, whether it's historical or live, date of the event.

type SessionBarData struct {
	EventName        string
	EventType        string
	CurrentLap       int
	TotalLaps        int
	FastestLapTime   time.Time
	FastestLapDriver string
	FastestLapNumber int
	StreamType       string
	EventDate        time.Time
}

func GetTestSessionBarData() SessionBarData {
	return SessionBarData{
		EventName:        "Japanese Grand Prix",
		EventType:        "Race",
		CurrentLap:       12,
		TotalLaps:        57,
		FastestLapTime:   time.Now(),
		FastestLapDriver: "NOR",
		FastestLapNumber: 10,
		StreamType:       "REPLAY",
		EventDate:        time.Now(),
	}
}

func SessionBar(d *SessionBarData) string {

	return defaultBorderStyle().Render(
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			defaultTextStyle(d.EventName, title1Color),
			defaultDivider,
			defaultTextStyle(d.EventType, title1Color),
			defaultDivider,
			defaultTextStyle(fmt.Sprintf("Lap %v/%v", d.CurrentLap, d.TotalLaps), title2Color),
			defaultDivider,
			defaultTextStyle(fmt.Sprintf("Fastest: %v %v L%v", d.FastestLapDriver, d.FastestLapTime.Format(time.TimeOnly), d.FastestLapNumber), title3Color),
			defaultDivider,
			defaultTextStyle(d.StreamType, titleDarkColor),
			defaultDivider,
			defaultTextStyle(d.EventDate.Format(time.DateOnly), titleDarkColor),
		))
}
