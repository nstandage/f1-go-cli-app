package aggregator

import (
	"time"

	"github.com/nstandage/f1-go-cli-app/model"
)

var (
	historyMaxLength     = 180
	raceControlMaxLength = 6
)

type Store struct {
	history              []model.Snapshot
	Drivers              map[uint]*Driver // mapped to DriverNumber
	RaceControl          []model.RaceControl
	Pitstops             []model.Pit
	TotalLaps            uint
	CurrentLap           uint
	IsReplay             bool
	Session              *model.Session
	Meeting              *model.Meeting
	StartingGrid         []model.StartingGrid
	FastestLap			 *FastestLap
}

type Driver struct {
	Number           uint
	Info             *model.Driver
	Position         uint
	StartingPosition uint
	IsOut            bool
	Interval         string
	ToLeader         string
	LastLap          float64
	LastLapIsPitOut  bool
	Stint            *model.Stint
}

type FastestLap struct {
	LapTime           float64
	DriverNumber 		 uint
	LapNumber     uint
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

func (s *Store) updateLap(data *model.Lap) {
	if data.LapDuration <= 0 {
		return
	}

	if data.LapNumber > s.CurrentLap {
		s.CurrentLap = data.LapNumber
	}
	go s.sleepForLapDuration(data)
}

func (s *Store) sleepForLapDuration(data *model.Lap) {
	duration := time.Duration(data.LapDuration * float64(time.Second))
	time.Sleep(duration)

	driver, ok := s.Drivers[data.DriverNumber]
	if !ok {
		return
	}

	driver.LastLap = data.LapDuration
	driver.LastLapIsPitOut = data.IsPitOutLap
	if s.FastestLap == nil || data.LapDuration < s.FastestLap.LapTime {
		s.FastestLap = &FastestLap{
			LapTime: data.LapDuration,
			DriverNumber: data.DriverNumber,
			LapNumber: data.LapNumber,
		}

	}
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
