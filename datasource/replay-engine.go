package datasource

import (
	"sort"
	"time"

	"github.com/nstandage/f1-go-cli-app/model"
)

type ReplayEngine struct {
	EventData *model.EventData
}

func (eng *ReplayEngine) Start(out chan<- *model.Event) {
	defer close(out)
	eng.sortEventData()

	for i, e := range eng.EventData.EventModels {
		if i > 0 {
			duration := e.GetDateStart().Sub(eng.EventData.EventModels[i-1].GetDateStart())
			time.Sleep(duration)
		}
		out <- &model.Event{Model: e}
	}
}

func (eng *ReplayEngine) sortEventData() {
	sort.Slice(eng.EventData.EventModels, func(i, j int) bool {
		return eng.EventData.EventModels[i].GetDateStart().Before(eng.EventData.EventModels[j].GetDateStart())
	})
}
