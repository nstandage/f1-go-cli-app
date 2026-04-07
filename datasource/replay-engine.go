package datasource

import (
	"fmt"

	"github.com/nstandage/f1-go-cli-app/model"
)

type ReplayEngine struct {
	SessionData *model.SessionData
	Channel     chan *model.Event
}

func (eng *ReplayEngine) StartStream() {
	fmt.Println("StartStream:")
	for _, i := range eng.SessionData.Intervals {
		ev := model.Event{
			Model: &i,
		}
		eng.Channel <- &ev
	}
	close(eng.Channel)
}
