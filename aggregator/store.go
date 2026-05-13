package aggregator

import (
	"github.com/nstandage/f1-go-cli-app/model"
)

var historyMaxLength = 180
var raceControlMaxLength = 6

type Store struct {
	history      []model.Snapshot
	Drivers      map[uint]*Driver // mapped to DriverNumber
	RaceControl  []model.RaceControl
	Pitstops     []model.Pit
	TotalLaps    uint
	IsReplay     bool
	Session      *model.Session
	Meeting      *model.Meeting
	StartingGrid []model.StartingGrid
}

type Driver struct {
	Number           uint
	Info             *model.Driver
	Position         uint
	StartingPosition uint
	IsOut            bool
	Interval         string
	ToLeader         string
	LastLap         float64
	LastLapIsPitOut bool
	Stint           *model.Stint
}

func (s *Store) updateHistory(h *model.Snapshot) {
	s.history = appendCapped(s.history, *h, historyMaxLength)
}

func (s *Store) updateRaceControl(rc *model.RaceControl) {
	s.RaceControl = appendCapped(s.RaceControl, *rc, raceControlMaxLength)
}

func (s *Store) updateInterval(i *model.Interval) {
	driver := s.Drivers[i.DriverNumber]
	driver.Interval = string(i.Interval)
	driver.ToLeader = string(i.GapToLeader)
}

func (s *Store) updatePosition(p *model.Position) {
	driver := s.Drivers[p.DriverNumber]
	driver.Position = p.Position
}

func appendCapped[T any](s []T, val T, max int) []T {
	s = append(s, val)
	if len(s) > max {
		s = s[1:]
	}
	return s
}
