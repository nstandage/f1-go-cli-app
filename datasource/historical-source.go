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
	Service *service.OpenF1Service
}

// Fetches data from server all at once.
func (hs *HistoricalSource) Start(ctx context.Context, sessionKey string) (*model.SessionData, error) {
	session, err := hs.Service.FetchIntervals(ctx, sessionKey)
	if err != nil {
		return nil, err
	}
	// fmt.Printf("Session: %+v\n", session)
	sessionData := model.SessionData{
		Intervals: *session,
	}
	return &sessionData, nil
}

// Starts replay engine
func (hs *HistoricalSource) Load(eng *ReplayEngine) {
	eng.StartStream()
}
