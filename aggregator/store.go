package aggregator

import (
	"context"
	"time"

	"github.com/nstandage/f1-go-cli-app/model"
)

var (
	historyMaxLength     = 180
	raceControlMaxLength = 6
	recentPitsMaxLength  = 8
)

type Store struct {
	history      []model.Snapshot
	Drivers      map[uint]*Driver // mapped to DriverNumber
	RaceControl  []model.RaceControl
	TotalLaps    uint
	CurrentLap   uint
	IsReplay     bool
	Session      *model.Session
	Meeting      *model.Meeting
	StartingGrid []model.StartingGrid
	FastestLap   *FastestLap
	Stints       map[uint][]model.Stint
	RecentPits   []model.PitStopEntry
	SectorCounts [3]int
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
	LapsOnTire       uint
	PitCount         uint
	CurrentLap       uint
	Sectors          [3][]uint
	cancelSectors    context.CancelFunc
}

type FastestLap struct {
	LapTime      float64
	DriverNumber uint
	LapNumber    uint
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

	if s.SectorCounts[0] == 0 && len(data.SegmentsSector1) > 0 {
		s.SectorCounts = [3]int{
			len(data.SegmentsSector1),
			len(data.SegmentsSector2),
			len(data.SegmentsSector3),
		}
	}

	if driver, ok := s.Drivers[data.DriverNumber]; ok {
		if data.LapNumber > driver.CurrentLap {
			driver.CurrentLap = data.LapNumber
		}
		if driver.cancelSectors != nil {
			driver.cancelSectors()
		}
		driver.Sectors = [3][]uint{}
		ctx, cancel := context.WithCancel(context.Background())
		driver.cancelSectors = cancel
		go s.animateSectors(ctx, driver, data)
	}

	go s.sleepForLapDuration(data)
}

func (s *Store) animateSectors(ctx context.Context, driver *Driver, data *model.Lap) {
	segments := [3][]uint{data.SegmentsSector1, data.SegmentsSector2, data.SegmentsSector3}
	durations := [3]float64{data.DurationSector1, data.DurationSector2, data.DurationSector3}

	for i, segs := range segments {
		if len(segs) == 0 {
			continue
		}
		delay := time.Duration(durations[i] / float64(len(segs)) * float64(time.Second))
		for _, seg := range segs {
			time.Sleep(delay)
			select {
			case <-ctx.Done():
				return
			default:
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
			driver.Sectors[i] = append(driver.Sectors[i], seg)
		}
	}
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
	driver.LapsOnTire++
	if s.FastestLap == nil || data.LapDuration < s.FastestLap.LapTime {
		s.FastestLap = &FastestLap{
			LapTime:      data.LapDuration,
			DriverNumber: data.DriverNumber,
			LapNumber:    data.LapNumber,
		}
	}
}

func (s *Store) updatePosition(p *model.Position) {
	driver := s.Drivers[p.DriverNumber]
	driver.Position = p.Position
}

func (s *Store) updatePit(p *model.Pit) {
	driver, ok := s.Drivers[p.DriverNumber]
	if !ok {
		return
	}
	driver.PitCount++
	entry := model.PitStopEntry{
		DriverAcronym: driver.Info.NameAcronym,
		StopDuration:  p.StopDuration,
	}
	s.RecentPits = append([]model.PitStopEntry{entry}, s.RecentPits...)
	if len(s.RecentPits) > recentPitsMaxLength {
		s.RecentPits = s.RecentPits[:recentPitsMaxLength]
	}
}

func appendCapped[T any](s []T, val T, max int) []T {
	s = append(s, val)
	if len(s) > max {
		s = s[1:]
	}
	return s
}
