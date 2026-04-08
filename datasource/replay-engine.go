package datasource

import (
	"sort"

	"github.com/nstandage/f1-go-cli-app/model"
)

type ReplayEngine struct {
	EventData *model.EventData
}

func (eng *ReplayEngine) Start(out chan *model.Event) {
	eng.sortEventData()
	for _, i := range eng.EventData.EventModels {
		ev := model.Event{
			Model: i,
		}
		out <- &ev
	}
	close(out)
}

func (eng *ReplayEngine) sortEventData() {
	sort.Slice(eng.EventData.EventModels, func(i, j int) bool {
		return eng.EventData.EventModels[i].GetDateStart().Before(eng.EventData.EventModels[j].GetDateStart())
	})
}
