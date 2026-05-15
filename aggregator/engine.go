package aggregator

import (
	"fmt"
	"math"
	"strconv"
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

func (e *Engine) Start() {
	raceData, c := e.datasource.Start()
	e.setUpInitialStore(raceData)
	e.listen(c)
}

func (e *Engine) setUpInitialStore(rd *model.RaceData) {
	e.store.Meeting = rd.Meeting
	e.store.Session = rd.Session
	e.store.TotalLaps = rd.TotalLaps
	e.store.StartingGrid = rd.StartingGrid
	e.store.IsReplay = e.datasource.IsReplay()
	e.store.Drivers = convertDrivers(rd.Drivers)
	e.store.Stints = make(map[uint][]model.Stint)
	for _, s := range rd.Stints {
		e.store.Stints[s.DriverNumber] = append(e.store.Stints[s.DriverNumber], s)
	}

	for _, sg := range e.store.StartingGrid {
		driver, ok := e.store.Drivers[sg.DriverNumber]
		if ok {
			driver.StartingPosition = sg.Position
			driver.Position = sg.Position
			driver.LastLap = sg.LapDuration
			driver.LapsOnTire = 0
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
func (e *Engine) listen(c <-chan *model.Event) {
	for event := range c {
		e.handle(event)
	}
}

func (e *Engine) handle(event *model.Event) {
	switch m := event.Model.(type) {
	case *model.Interval:
		e.updateInterval(m)
	case *model.Lap:
		e.updateLap(m)
	case *model.Location:
		e.updateLocation(m)
	case *model.Position:
		e.updatePosition(m)
	case *model.RaceControl:
		e.store.updateRaceControl(m)
	case *model.Pit:
		e.updatePit(m)
	}
}

func (e *Engine) updateInterval(data *model.Interval) {
	e.store.updateInterval(data)
}

func (e *Engine) updateLap(data *model.Lap) {
	e.store.updateLap(data)
}

func (e *Engine) updateLocation(data *model.Location) {
}

func (e *Engine) updateMeeting(data *model.Meeting) {
}

func (e *Engine) updatePosition(data *model.Position) {
	e.store.updatePosition(data)
}

func (e *Engine) updatePit(data *model.Pit) {
	e.store.updatePit(data)
}

func (e *Engine) updateSesion(data *model.Session) {
}

func (e *Engine) updateStartingGrid(data []model.StartingGrid) {
}

// MARK: Snapshot Functions

func (e *Engine) GetSnapshot(offset uint) *model.Snapshot {
	sessionBar := model.SessionBarSnapShot{
		EventName:  e.store.Meeting.MeetingOfficialName,
		EventType:  e.store.Session.SessionType,
		CurrentLap: e.store.CurrentLap,
		FastestLap: e.getFastestLap(),
		TotalLaps:  e.store.TotalLaps,
		IsReplay:   e.store.IsReplay,
		EventDate:  e.store.Session.DateStart,
	}
	lastLap, lastLapIsPitOut := e.getLastLap()
	snapshot := model.Snapshot{
		SessionBar:      &sessionBar,
		RaceControlMsgs: e.getRaceControlMessages(),
		Drivers:         e.getDriverNames(),
		LastLap:         lastLap,
		LastLapIsPitOut: lastLapIsPitOut,
		Intervals:       e.getIntervals(),
		GapsToLeaders:   e.getGapToLeader(),
		PitCounts:       e.getPitCounts(),
		TireCompounds:   e.getTireCompounds(),
		TireAges:        e.getTireAges(),
	}

	e.store.updateHistory(&snapshot)

	return &snapshot
}

func (e *Engine) getFastestLap() *model.FastestLapSnapshot {
	if e.store.FastestLap == nil {
		return &model.FastestLapSnapshot{
			LapTime: formatLapTime(0.0),
			Driver: "---",
			LapNumber: "---",
		}
	}

	driverNumber := e.store.FastestLap.DriverNumber
	driverName := e.store.Drivers[driverNumber].Info.NameAcronym
	lapNumber := e.store.FastestLap.LapNumber
	return &model.FastestLapSnapshot{
		LapTime:   formatLapTime(e.store.FastestLap.LapTime),
		Driver:    driverName,
		LapNumber: strconv.FormatUint(uint64(lapNumber), 10),
	}
}

func (e *Engine) getRaceControlMessages() []string {
	strs := []string{}
	for _, rc := range e.store.RaceControl {
		strs = append(strs, rc.Message)
	}

	return strs
}

func (e *Engine) getDriverNames() []model.DriverSnapshot {
	drivers := make([]model.DriverSnapshot, len(e.store.Drivers))
	for _, d := range e.store.Drivers {
		name := d.Info.NameAcronym
		isFastest := false
		if e.store.FastestLap != nil {
			isFastest = e.store.FastestLap.DriverNumber == d.Number
		}
		drivers[d.Position-1] = model.DriverSnapshot{Name: name, IsFastestLap: isFastest}
	}

	return drivers
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

func (e *Engine) getPitCounts() []string {
	counts := make([]string, len(e.store.Drivers))
	for _, d := range e.store.Drivers {
		counts[d.Position-1] = strconv.Itoa(int(d.PitCount))
	}
	return counts
}

func (e *Engine) getTireCompounds() []string {
	compounds := make([]string, len(e.store.Drivers))
	for _, d := range e.store.Drivers {
		stint := getActiveStint(e.store.Stints[d.Number], d.CurrentLap)
		if stint == nil {
			compounds[d.Position-1] = "---"
		} else {
			compounds[d.Position-1] = normalizeCompound(stint.Compound)
		}
	}
	return compounds
}

func (e *Engine) getTireAges() []string {
	ages := make([]string, len(e.store.Drivers))
	for _, d := range e.store.Drivers {
		stint := getActiveStint(e.store.Stints[d.Number], d.CurrentLap)
		if stint == nil {
			ages[d.Position-1] = "--"
		} else {
			ages[d.Position-1] = strconv.FormatUint(uint64(tireAge(stint, d.CurrentLap)), 10)
		}
	}
	return ages
}
