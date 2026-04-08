package aggregator

import (
	"fmt"

	"github.com/nstandage/f1-go-cli-app/model"
)

type Engine struct {
	RaceData *model.RaceData
}

func (eng *Engine) Start(out chan *model.Event) { // Drivers, laps, pits, stint
	eng.updateSesion(eng.RaceData.Session)
	eng.updateMeeting(eng.RaceData.Meeting)
	for event := range out {
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
