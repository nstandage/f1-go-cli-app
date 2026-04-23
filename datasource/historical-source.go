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
	Start() (*model.RaceData, <-chan *model.Event)
	IsReplay() bool
}

// Historical Source handles fetching all data from selected session at once and returns a sessionData object.
type HistoricalSource struct {
	service   *service.OpenF1HTTP
	raceData  *model.RaceData
	eventData *model.EventData
}

func NewHistoricalSource(s *service.OpenF1HTTP) *HistoricalSource {
	return &HistoricalSource{
		service:   s,
		raceData:  &model.RaceData{},
		eventData: &model.EventData{},
	}
}

// Fetches data from server all at once.
func (hs *HistoricalSource) Fetch(ctx context.Context, sessionKey string, meetingKey string) error {
	rl := NewRateLimiter(2)
	defer rl.Stop()

	meetings, err := hs.getMeetings(ctx, rl, meetingKey)
	if err != nil {
		return fmt.Errorf("HistoricaSource.Fetch - meetings failed: %w", err)
	}

	if len(meetings) == 0 {
		return fmt.Errorf("HistoricalSource.Fetch - meetings is 0 %w", err)
	}

	sessions, err := hs.getMeetingSessions(ctx, rl, meetingKey)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - sessions failed %w", err)
	}

	raceSession, err := getSessionByTypeAndName(sessions, RaceType, Race)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - sessions is 0 %w", err)
	}

	raceControls, err := hs.getRaceControls(ctx, rl, sessionKey)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - raceControls failed %w", err)
	}

	if len(raceControls) == 0 {
		return fmt.Errorf("HistoricalSource.Fetch - raceControls is 0 %w", err)
	}

	qSession, err := getSessionByTypeAndName(sessions, QualifyingType, Qualifying)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - Qualifying Session is 0 %w", err)
	}

	grid, err := hs.getStartingGrid(ctx, rl, qSession.GetSessionKey())
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - Starting Grid failed %w", err)
	}

	drivers, err := hs.getDrivers(ctx, rl, sessionKey)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - drivers failed %w", err)
	}

	hs.raceData.Meeting = &meetings[0]
	hs.raceData.Session = raceSession
	hs.raceData.TotalLaps = getLapCount(raceControls)
	hs.raceData.StartingGrid = grid
	hs.raceData.Drivers = drivers
	for _, rc := range raceControls {
		hs.eventData.EventModels = append(hs.eventData.EventModels, &rc)
	}
	return nil
}

func (hs *HistoricalSource) Start() (*model.RaceData, <-chan *model.Event) {
	replayEngine := ReplayEngine{EventData: hs.eventData}
	c := make(chan *model.Event)
	go replayEngine.Start(c)
	return hs.raceData, c
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

func (hs *HistoricalSource) IsReplay() bool {
	return true
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

func (hs *HistoricalSource) getMeetings(ctx context.Context, rl *RateLimiter, meetingKey string) ([]model.Meeting, error) {
	rl.Wait()
	return hs.service.FetchMeetings(ctx, meetingKey)
}

func (hs *HistoricalSource) getSessions(ctx context.Context, rl *RateLimiter, sessionKey string) ([]model.Session, error) {
	rl.Wait()
	return hs.service.FetchSessions(ctx, sessionKey)
}

func (hs *HistoricalSource) getMeetingSessions(ctx context.Context, rl *RateLimiter, meetingKey string) ([]model.Session, error) {
	rl.Wait()
	return hs.service.FetchMeetingSessions(ctx, meetingKey)
}

func (hs *HistoricalSource) getRaceControls(ctx context.Context, rl *RateLimiter, sessionKey string) ([]model.RaceControl, error) {
	rl.Wait()
	return hs.service.FetchRaceControls(ctx, sessionKey)
}

// API requires a Qualifying session_key
func (hs *HistoricalSource) getStartingGrid(ctx context.Context, rl *RateLimiter, sessionKey string) ([]model.StartingGrid, error) {
	rl.Wait()
	return hs.service.FetchStartingGrid(ctx, sessionKey)
}

func (hs *HistoricalSource) getDrivers(ctx context.Context, rl *RateLimiter, sessionKey string) ([]model.Driver, error) {
	rl.Wait()
	return hs.service.FetchDrivers(ctx, sessionKey)
}
