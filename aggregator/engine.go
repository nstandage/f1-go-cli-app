package aggregator

import (
	"fmt"
	"math"
	"github.com/nstandage/f1-go-cli-app/datasource"
	"github.com/nstandage/f1-go-cli-app/model"
)

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

// MARK: Channel functions
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
		eng.store.updateRaceControl(m)
	}
}

func (e *Engine) updateInterval(data *model.Interval) {
	e.store.updateInterval(data)
}

func (e *Engine) updateLap(data *model.Lap) {
	if data.LapDuration <= 0 {
		return
	}
	driver, ok := e.store.Drivers[data.DriverNumber]
	if !ok {
		return
	}
	driver.LastLap = data.LapDuration
	driver.LastLapIsPitOut = data.IsPitOutLap
}

func (e *Engine) updateLocation(data *model.Location) {

}

func (e *Engine) updateMeeting(data *model.Meeting) {

}

func (e *Engine) updatePosition(data *model.Position) {
	e.store.updatePosition(data)
}

func (e *Engine) updateSesion(data *model.Session) {

}

func (e *Engine) updateStartingGrid(data []model.StartingGrid) {

}

// MARK: Snapshot Functions
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
	lastLap, lastLapIsPitOut := e.getLastLap()
	snapshot := model.Snapshot{
		SessionBar:      &sessionBar,
		RaceControlMsgs: e.getRaceControlMessages(),
		DriverNames:     e.getDriverNames(),
		LastLap:         lastLap,
		LastLapIsPitOut: lastLapIsPitOut,
		Intervals:       e.getIntervals(),
		GapsToLeaders:   e.getGapToLeader(),
	}

	e.store.updateHistory(&snapshot)

	return &snapshot
}

func (e *Engine) getRaceControlMessages() []string {
	strs := []string{}
	for _, rc := range e.store.RaceControl {
		strs = append(strs, rc.Message)
	}

	return strs
}

func (e *Engine) getDriverNames() []string {
	strs := make([]string, len(e.store.Drivers))
	for _, d := range e.store.Drivers {
		strs[d.Position-1] = d.Info.BroadcastName
	}

	return strs
}

func (e *Engine) getLastLap() ([]string, []bool) {
	strs := make([]string, len(e.store.Drivers))
	isPitOut := make([]bool, len(e.store.Drivers))
	for _, d := range e.store.Drivers {
		strs[d.Position-1] = formatLapTime(d.LastLap)
		isPitOut[d.Position-1] = d.LastLapIsPitOut
	}
	return strs, isPitOut
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

func (e *Engine) getIntervals() []string {
	strs := make([]string, len(e.store.Drivers))
	for _, d := range e.store.Drivers {
		strs[d.Position-1] = getTimingText(d.Interval)
	}
	return strs
}

func getTimingText(i string) string {
	if i == "" {
		return "----"
	}
	return i
}

func (e *Engine) getGapToLeader() []string {
	strs := make([]string, len(e.store.Drivers))
	for _, d := range e.store.Drivers {
		strs[d.Position-1] = getTimingText(d.ToLeader)
	}
	return strs
}
