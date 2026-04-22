package aggregator

import (
	"github.com/nstandage/f1-go-cli-app/model"
)

type Datasource struct {
	history      []Snapshot
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

func ConvertDrivers(ds []model.Driver) map[uint]*Driver {
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

func (ds *Datasource) AddStartingGrid() {
	for _, sg := range ds.StartingGrid {
		driver, ok := ds.Drivers[sg.DriverNumber]
		if ok {
			driver.StartingPosition = sg.Position
			driver.Position = sg.Position
			driver.LastLap = sg.LapDuration
		}
	}
}
