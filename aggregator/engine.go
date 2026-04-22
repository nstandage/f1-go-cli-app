package aggregator

import (
	"fmt"
	"math"

	"github.com/nstandage/f1-go-cli-app/model"
)

var numberOfRaceControls = 6

type Engine struct {
	Datasource *Datasource
}

func (eng *Engine) Start(out chan *model.Event) { // Drivers, laps, pits, stint
	for event := range out {
		eng.handle(event)
	}
}

func (eng *Engine) handle(e *model.Event) {
	switch m := e.Model.(type) {
	case *model.Interval:
		eng.updateInterval(m)
	case *model.Lap:
		eng.updateLap(m)
	case *model.Location:
		eng.updateLocation(m)
	case *model.Position:
		eng.updatePosition(m)
	case *model.RaceControl:
		eng.updateRaceControl(m)
	}
}

func (e *Engine) updateInterval(data *model.Interval) {
	// e.Program.Send(data)
}

func (e *Engine) updateLap(data *model.Lap) {

}

func (e *Engine) updateLocation(data *model.Location) {

}

func (e *Engine) updateMeeting(data *model.Meeting) {

}

func (e *Engine) updatePosition(data *model.Position) {

}

func (e *Engine) updateRaceControl(rc *model.RaceControl) {
	e.Datasource.RaceControl = appendCapped(e.Datasource.RaceControl, *rc, numberOfRaceControls)
}

func (e *Engine) updateSesion(data *model.Session) {

}

func (e *Engine) updateStartingGrid(data []model.StartingGrid) {

}

func (e *Engine) GetSnapshot(offset uint) Snapshot {
	sessionBar := SessionBarSnapShot{
		EventName:        e.Datasource.Meeting.MeetingOfficialName,
		EventType:        e.Datasource.Session.SessionType,
		CurrentLap:       0,
		FastestLapNumber: 11,
		TotalLaps:        e.Datasource.TotalLaps,
		IsReplay:         e.Datasource.IsReplay,
		EventDate:        e.Datasource.Session.DateStart,
	}
	return Snapshot{
		SessionBar:      sessionBar,
		RaceControlMsgs: e.getRaceControlMessages(),
		DriverNames:     e.getDriverNames(),
		LastLap:         e.getLastLap(),
	}
}

func (e *Engine) getRaceControlMessages() []string {
	strs := []string{}
	for _, rc := range e.Datasource.RaceControl {
		strs = append(strs, rc.Message)
	}

	return strs
}

func (e *Engine) HistoryLen() int {
	return len(e.Datasource.history)
}

func appendCapped[T any](s []T, val T, max int) []T {
	s = append(s, val)
	if len(s) > max {
		s = s[1:]
	}
	return s
}

func (e *Engine) getDriverNames() []string {
	strs := make([]string, len(e.Datasource.Drivers))
	for _, d := range e.Datasource.Drivers {
		strs[d.Position-1] = d.Info.BroadcastName
	}

	return strs
}

func (e *Engine) getLastLap() []string {
	strs := make([]string, len(e.Datasource.Drivers))
	for _, d := range e.Datasource.Drivers {
		// lastLapInSec := d.LastLap
		// times := time.Duration(lastLapInSec * float64(time.Second))
		// str := times.Minutes()
		// strs[d.Position-1] = strconv.FormatFloat(str, 'f', 2, 64)
		strs[d.Position-1] = formatLapTime(d.LastLap)
	}
	return strs
}

func formatLapTime(seconds float64) string {
	if seconds <= 0 {
		return "--:--.---"
	}

	totalMs := int(math.Round(seconds * 1000))
	totalSecs := totalMs / 1000
	ms := totalMs % 1000
	secs := totalSecs % 60
	mins := totalSecs / 60

	return fmt.Sprintf("%d:%02d.%03d", mins, secs, ms)
}
