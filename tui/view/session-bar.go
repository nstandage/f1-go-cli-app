package view

import (
	"fmt"
	"time"
	"charm.land/lipgloss/v2"
	"github.com/nstandage/f1-go-cli-app/model"
)

func formatElapsed(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%d:%02d:%02d", h, m, s)
}

// Contains event name/type lap count, fastest lap, whether it's historical or live, date of the event.
func SessionBar(s *model.SessionBarSnapShot) string {
	var replayType string
	if s.IsReplay {
		replayType = "REPLAY"
	} else {
		replayType = "LIVE"
	}
	return defaultBorderStyle().Render(
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			defaultTextStyle(s.EventName, title1Color),
			defaultDivider,
			defaultTextStyle(s.EventType, title1Color),
			defaultDivider,
			defaultTextStyle(fmt.Sprintf("Lap %v/%v", s.CurrentLap, s.TotalLaps), title2Color),
			defaultDivider,
			defaultTextStyle(fmt.Sprintf("Fastest: %v %v Lap %v", s.FastestLap.Driver, s.FastestLap.LapTime, s.FastestLap.LapNumber), title3Color),
			defaultDivider,
			defaultTextStyle(replayType, titleDarkColor),
			defaultDivider,
			defaultTextStyle(s.EventDate.Format("Mon, January 02, 2006"), titleDarkColor),
			defaultDivider,
			defaultTextStyle(formatElapsed(s.RaceElapsed), title2Color),
		))
}
