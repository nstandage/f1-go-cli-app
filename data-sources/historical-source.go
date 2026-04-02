package datasource

import (
	"context"
	"fmt"

	"github.com/nstandage/f1-go-cli-app/models"
	"github.com/nstandage/f1-go-cli-app/service"
)

type DataSource interface {
	Start(meetingId string, ctx context.Context, out chan<- string)
}

type HistoricalSource struct {
	Service *service.OpenF1Service
}

func (hs *HistoricalSource) Start(sessionKey string, ctx context.Context) (*models.SessionData, error) {
	session, err := hs.Service.FetchSessions(ctx, sessionKey)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Session: %+v\n", session)

	return nil, nil
}

func (hs *HistoricalSource) Load(ctx context.Context, out chan<- models.Event) {

}
