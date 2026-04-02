package datasource

import (
	"context"
	"fmt"

	models "github.com/nstandage/f1-go-cli-app/model"
	"github.com/nstandage/f1-go-cli-app/service"
)

type DataSource interface {
	Start(meetingId string, ctx context.Context, out chan<- string)
}

type HistoricalSource struct {
	Service *service.OpenF1Service
}

func (hs *HistoricalSource) Start(ctx context.Context, sessionKey string) (*models.SessionData, error) {
	session, err := hs.Service.FetchIntervals(ctx, sessionKey)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Session: %+v\n", session)

	sessionData := models.SessionData{
		Intervals: session,
	}
	return &sessionData, nil
}

func (hs *HistoricalSource) Load(sessionData *models.SessionData, out chan<- models.Event) {

}
