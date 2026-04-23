package view

import (
	"fmt"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/nstandage/f1-go-cli-app/model"
)

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
			defaultTextStyle(fmt.Sprintf("Fastest: %v %v L%v", s.FastestLapDriver, s.FastestLapTime.Format(time.TimeOnly), s.FastestLapNumber), title3Color),
			defaultDivider,
			defaultTextStyle(replayType, titleDarkColor),
			defaultDivider,
			defaultTextStyle(s.EventDate.Format("Mon, January 02, 2006"), titleDarkColor),
		))
}
