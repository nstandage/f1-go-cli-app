package aggregator

import (
	"github.com/nstandage/f1-go-cli-app/model"
)

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
	Interval         float32
	ToLeader         float32
	LastLap          float64
	Stint            *model.Stint
}
