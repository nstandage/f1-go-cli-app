package datasource

import (
	"context"
	"fmt"

	"github.com/nstandage/f1-go-cli-app/model"
	"github.com/nstandage/f1-go-cli-app/service"
)

type SessionType string
type SessionName string

const (
	FPType         SessionType = "Practice"
	QualifyingType SessionType = "Qualifying"
	RaceType       SessionType = "Race"
)

const (
	FP1              SessionName = "Practice 1"
	FP2              SessionName = "Practice 2"
	FP3              SessionName = "Practice 3"
	SprintQualifying SessionName = "Sprint Qualifying"
	Sprint           SessionName = "Sprint"
	Qualifying       SessionName = "Qualifying"
	Race             SessionName = "Race"
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

	sessions, err := hs.getMeetingSessions(ctx, rl, meetingKey)
	if err != nil {
		return nil, nil, fmt.Errorf("HistoricalSource.Fetch - sessions failed %w", err)
	}

	raceSession, err := getSessionByTypeAndName(sessions, RaceType, Race)
	if err != nil {
		return nil, nil, fmt.Errorf("HistoricalSource.Fetch - sessions is 0 %w", err)
	}

	raceControls, err := hs.getRaceControls(ctx, rl, sessionKey)
	if err != nil {
		return nil, nil, fmt.Errorf("HistoricalSource.Fetch - raceControls failed %w", err)
	}

	if len(raceControls) == 0 {
		return nil, nil, fmt.Errorf("HistoricalSource.Fetch - raceControls is 0 %w", err)
	}

	qSession, err := getSessionByTypeAndName(sessions, QualifyingType, Qualifying)
	if err != nil {
		return nil, nil, fmt.Errorf("HistoricalSource.Fetch - Qualifying Session is 0 %w", err)
	}

	grid, err := hs.getStartingGrid(ctx, rl, qSession.GetSessionKey())
	if err != nil {
		return nil, nil, fmt.Errorf("HistoricalSource.Fetch - Starting Grid failed %w", err)
	}

	drivers, err := hs.getDrivers(ctx, rl, sessionKey)
	if err != nil {
		return nil, nil, fmt.Errorf("HistoricalSource.Fetch - drivers failed %w", err)
	}

	raceData.Meeting = &meetings[0]
	raceData.Session = raceSession
	raceData.TotalLaps = getLapCount(raceControls)
	raceData.StartingGrid = grid
	raceData.Drivers = drivers
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

func getSessionByTypeAndName(ss []model.Session, st SessionType, sn SessionName) (*model.Session, error) {
	for _, s := range ss {
		if s.SessionType == string(st) && s.SessionName == string(sn) {
			return &s, nil
		}
	}

	var meetingKey uint = 0
	err := fmt.Errorf("Couldn't find Session of type: %v and name: %v", st, sn)
	if len(ss) > 0 {
		meetingKey = ss[0].MeetingKey
		err = fmt.Errorf("%v for meeting_key: %v", err, meetingKey)
	}
	return nil, err
}

func (hs *HistoricalSource) getMeetings(ctx context.Context, rl *service.RateLimiter, meetingKey string) ([]model.Meeting, error) {
	rl.Wait()
	return hs.Service.FetchMeetings(ctx, meetingKey)
}

func (hs *HistoricalSource) getSessions(ctx context.Context, rl *service.RateLimiter, sessionKey string) ([]model.Session, error) {
	rl.Wait()
	return hs.Service.FetchSessions(ctx, sessionKey)
}

func (hs *HistoricalSource) getMeetingSessions(ctx context.Context, rl *service.RateLimiter, meetingKey string) ([]model.Session, error) {
	rl.Wait()
	return hs.Service.FetchMeetingSessions(ctx, meetingKey)
}

func (hs *HistoricalSource) getRaceControls(ctx context.Context, rl *service.RateLimiter, sessionKey string) ([]model.RaceControl, error) {
	rl.Wait()
	return hs.Service.FetchRaceControls(ctx, sessionKey)
}

// API requires a Qualifying session_key
func (hs *HistoricalSource) getStartingGrid(ctx context.Context, rl *service.RateLimiter, sessionKey string) ([]model.StartingGrid, error) {
	rl.Wait()
	return hs.Service.FetchStartingGrid(ctx, sessionKey)
}

func (hs *HistoricalSource) getDrivers(ctx context.Context, rl *service.RateLimiter, sessionKey string) ([]model.Driver, error) {
	rl.Wait()
	return hs.Service.FetchDrivers(ctx, sessionKey)
}
