package datasource

import (
	"context"

	"github.com/nstandage/f1-go-cli-app/model"
	"github.com/nstandage/f1-go-cli-app/service"
)

type DataSource interface {
	Start(meetingId string, ctx context.Context, out chan<- string)
}

// Historical Source handles fetching all data from selected session at once and returns a sessionData object.
type HistoricalSource struct {
	Service *service.OpenF1HTTP
}

// Fetches data from server all at once.
func (hs *HistoricalSource) Fetch(ctx context.Context, sessionKey string) (*model.RaceData, *model.EventData, error) {
	var raceData model.RaceData
	var eventData model.EventData
	intervals, err := hs.Service.FetchIntervals(ctx, sessionKey)
	if err != nil {
		return nil, nil, err
	}

	for _, i := range *intervals {
		eventData.EventModels = append(eventData.EventModels, &i)
	}

	return &raceData, &eventData, nil
}
