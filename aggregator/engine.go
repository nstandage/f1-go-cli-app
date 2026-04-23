package aggregator

import (
	"fmt"
	"math"

	"github.com/nstandage/f1-go-cli-app/datasource"
	"github.com/nstandage/f1-go-cli-app/model"
)

var numberOfRaceControls = 6

type Engine struct {
	store      *Store
	datasource datasource.DataSource
}

func NewEngine(ds datasource.DataSource) *Engine {
	return &Engine{
		store:      &Store{},
		datasource: ds,
	}
}

func (eng *Engine) Start() {
	raceData, c := eng.datasource.Start()
	eng.setUpInitialStore(raceData)
	eng.listen(c)
}

func (eng *Engine) listen(c <-chan *model.Event) {
	for event := range c {
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
	e.store.RaceControl = appendCapped(e.store.RaceControl, *rc, numberOfRaceControls)
}

func (e *Engine) updateSesion(data *model.Session) {

}

func (e *Engine) updateStartingGrid(data []model.StartingGrid) {

}

func (e *Engine) GetSnapshot(offset uint) *model.Snapshot {
	sessionBar := model.SessionBarSnapShot{
		EventName:        e.store.Meeting.MeetingOfficialName,
		EventType:        e.store.Session.SessionType,
		CurrentLap:       0,
		FastestLapNumber: 11,
		TotalLaps:        e.store.TotalLaps,
		IsReplay:         e.store.IsReplay,
		EventDate:        e.store.Session.DateStart,
	}
	return &model.Snapshot{
		SessionBar:      &sessionBar,
		RaceControlMsgs: e.getRaceControlMessages(),
		DriverNames:     e.getDriverNames(),
		LastLap:         e.getLastLap(),
	}
}

func (e *Engine) getRaceControlMessages() []string {
	strs := []string{}
	for _, rc := range e.store.RaceControl {
		strs = append(strs, rc.Message)
	}

	return strs
}

func (e *Engine) HistoryLen() int {
	return len(e.store.history)
}

func appendCapped[T any](s []T, val T, max int) []T {
	s = append(s, val)
	if len(s) > max {
		s = s[1:]
	}
	return s
}

func (e *Engine) getDriverNames() []string {
	strs := make([]string, len(e.store.Drivers))
	for _, d := range e.store.Drivers {
		strs[d.Position-1] = d.Info.BroadcastName
	}

	return strs
}

func (e *Engine) getLastLap() []string {
	strs := make([]string, len(e.store.Drivers))
	for _, d := range e.store.Drivers {
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

func (eng *Engine) setUpInitialStore(rd *model.RaceData) {
	eng.store.Meeting = rd.Meeting
	eng.store.Session = rd.Session
	eng.store.TotalLaps = rd.TotalLaps
	eng.store.StartingGrid = rd.StartingGrid
	eng.store.IsReplay = eng.datasource.IsReplay()
	eng.store.Drivers = convertDrivers(rd.Drivers)

	for _, sg := range eng.store.StartingGrid {
		driver, ok := eng.store.Drivers[sg.DriverNumber]
		if ok {
			driver.StartingPosition = sg.Position
			driver.Position = sg.Position
			driver.LastLap = sg.LapDuration
		}
	}
}

func convertDrivers(ds []model.Driver) map[uint]*Driver {
	drivers := make(map[uint]*Driver, len(ds))
	for _, d := range ds {
		driver := Driver{
			Number: d.DriverNumber,
			Info:   &d,
		}
		drivers[d.DriverNumber] = &driver
	}

	return drivers
}
