package datasource

import (
	"context"
	"fmt"

	"github.com/nstandage/f1-go-cli-app/model"
	"github.com/nstandage/f1-go-cli-app/service"
)

type DataSource interface {
	Start(meetingId string, ctx context.Context, out chan<- model.Event)
	IsReplay() bool
}

// Historical Source handles fetching all data from selected session at once and returns a sessionData object.
type HistoricalSource struct {
	Service *service.OpenF1HTTP
}

// Fetches data from server all at once.
func (hs *HistoricalSource) Fetch(ctx context.Context, sessionKey string, meetingKey string) (*model.RaceData, *model.EventData, error) {
	rl := service.NewRateLimiter(2)
	defer rl.Stop()
	var raceData model.RaceData
	var eventData model.EventData

	meetings, err := hs.getMeetings(ctx, rl, meetingKey)
	if err != nil {
		return nil, nil, fmt.Errorf("HistoricaSource.Fetch - meetings failed: %w", err)
	}

	if len(meetings) == 0 {
		return nil, nil, fmt.Errorf("HistoricalSource.Fetch - meetings is 0 %w", err)
	}

	sessions, err := hs.getSessions(ctx, rl, sessionKey)
	if err != nil {
		return nil, nil, fmt.Errorf("HistoricalSource.Fetch - sessions failed %w", err)
	}

	if len(sessions) == 0 {
		return nil, nil, fmt.Errorf("HistoricalSource.Fetch - sessions is 0 %w", err)
	}

	raceControls, err := hs.getRaceControls(ctx, rl, sessionKey)
	if err != nil {
		return nil, nil, fmt.Errorf("HistoricalSource.Fetch - raceControls failed %w", err)
	}

	if len(raceControls) == 0 {
		return nil, nil, fmt.Errorf("HistoricalSource.Fetch - raceControls is 0 %w", err)
	}

	raceData.Meeting = &meetings[0]
	raceData.Session = &sessions[0]
	raceData.TotalLaps = getLapCount(raceControls)

	for _, rc := range raceControls {
		eventData.EventModels = append(eventData.EventModels, &rc)
	}

	return &raceData, &eventData, nil
}

func getLapCount(rcs []model.RaceControl) uint {
	for _, rc := range rcs {
		if rc.Flag == "CHEQUERED" {
			return rc.LapNumber
		}
	}
	return getLapCountByNumber(rcs)
}

func getLapCountByNumber(rcs []model.RaceControl) uint {
	var count uint = 0
	for _, rc := range rcs {
		if rc.LapNumber > count {
			count = rc.LapNumber
		}
	}
	return count
}

func (hs *HistoricalSource) getMeetings(ctx context.Context, rl *service.RateLimiter, meetingKey string) ([]model.Meeting, error) {
	rl.Wait()
	return hs.Service.FetchMeetings(ctx, meetingKey)
}

func (hs *HistoricalSource) getSessions(ctx context.Context, rl *service.RateLimiter, sessionKey string) ([]model.Session, error) {
	rl.Wait()
	return hs.Service.FetchSessions(ctx, sessionKey)
}

func (hs *HistoricalSource) getRaceControls(ctx context.Context, rl *service.RateLimiter, sessionKey string) ([]model.RaceControl, error) {
	rl.Wait()
	return hs.Service.FetchRaceControls(ctx, sessionKey)
}
