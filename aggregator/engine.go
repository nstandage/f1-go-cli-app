package aggregator

import (
	"fmt"

	"github.com/nstandage/f1-go-cli-app/model"
)

type Engine struct{}

func (eng *Engine) Start(out chan *model.Event) { // Drivers, laps, pits, stint
	for event := range out {
		eng.Handle(event)
	}
}

func (eng *Engine) Handle(e *model.Event) {
	switch m := e.Model.(type) {
	case *model.Interval:
		eng.updateInterval(m)
	case *model.Lap:
		eng.updateLap(m)
	case *model.Location:
		eng.updateLocation(m)
	case *model.Meeting:
		eng.updateMeeting(m)
	case *model.Position:
		eng.updatePosition(m)
	case *model.RaceControl:
		eng.updateRaceControl(m)
	case *model.Session:
		eng.updateSesion(m)
	}
}

func (e *Engine) updateInterval(data *model.Interval) {
	fmt.Printf("Interval: %v\n", data.DateStart)
}

func (e *Engine) updateLap(data *model.Lap) {

}

func (e *Engine) updateLocation(data *model.Location) {

}

func (e *Engine) updateMeeting(data *model.Meeting) {

}

func (e *Engine) updatePosition(data *model.Position) {

}

func (e *Engine) updateRaceControl(data *model.RaceControl) {

}

func (e *Engine) updateSesion(data *model.Session) {

}
