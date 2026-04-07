package datasource

import (
	"github.com/nstandage/f1-go-cli-app/model"
)

type ReplayEngine struct {
	SessionData *model.SessionData
}

func (eng *ReplayEngine) Start(out chan *model.Event) {
	for _, i := range eng.SessionData.Intervals {
		ev := model.Event{
			Model: &i,
		}
		out <- &ev
	}
	close(out)
}
