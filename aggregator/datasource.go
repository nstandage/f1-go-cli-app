package aggregator

import (
	"time"

	"github.com/nstandage/f1-go-cli-app/model"
)

type Datasource struct {
	Meeting     model.Meeting
	Session     model.Session
	history     []Snapshot
	drivers     map[uint]Driver // mapped to DriverNumber
	RaceControl []model.RaceControl
	Pitstops    []model.Pit
	TotalLaps   uint
	IsReplay    bool
}

type Driver struct {
	Number           uint
	Info             *model.Driver
	Position         uint
	StartingPosition uint
	IsOut            bool
	Interval         float32
	ToLeader         float32
	LastLap          time.Time
	Stint            *model.Stint
}

type StartingInfo struct {
	Session      *model.Session
	Meeting      *model.Meeting
	StartingGrid []model.StartingGrid
	drivers      []model.Driver
}
