package datasource

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nstandage/f1-go-cli-app/model"
	"github.com/nstandage/f1-go-cli-app/service"
)

type (
	SessionType string
	SessionName string
)

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
	Fetch(context.Context, string, string) error
}

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

func (hs *HistoricalSource) Start() (*model.RaceData, <-chan *model.Event) {
	replayEngine := ReplayEngine{EventData: hs.eventData}
	c := make(chan *model.Event)
	go replayEngine.Start(c, hs.raceData.SessionStart)
	return hs.raceData, c
}

func (hs *HistoricalSource) IsReplay() bool {
	return true
}

func (hs *HistoricalSource) Fetch(ctx context.Context, sessionKey string, meetingKey string) error {
	meetings, err := hs.service.FetchMeetings(ctx, meetingKey)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - meetings failed: %w", err)
	}

	if len(meetings) == 0 {
		return fmt.Errorf("HistoricalSource.Fetch - meetings is 0")
	}

		sessions, err := hs.service.FetchSessionsByMeeting(ctx, meetingKey)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - sessions failed %w", err)
	}

	raceSession, err := getSessionByTypeAndName(sessions, RaceType, Race)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - sessions - %w", err)
	}

	raceControls, err := hs.service.FetchRaceControls(ctx, sessionKey)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - raceControls failed %w", err)
	}

	startTime, err := hs.getStartTime(raceControls)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - race control get startTime returned nil %w", err)
	}

	qSession, err := getSessionByTypeAndName(sessions, QualifyingType, Qualifying)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - Qualifying Session - %w", err)
	}

	grid, err := hs.service.FetchStartingGrid(ctx, qSession.GetSessionKey())
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - Starting Grid failed %w", err)
	}

	drivers, err := hs.service.FetchDrivers(ctx, sessionKey)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - drivers failed %w", err)
	}

	positions, err := hs.service.FetchPositions(ctx, sessionKey)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - positions failed %w", err)
	}

	intervals, err := hs.service.FetchIntervals(ctx, sessionKey)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - intervals failed %w", err)
	}

	laps, err := hs.service.FetchLaps(ctx, sessionKey)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - laps failed %w", err)
	}

	stints, err := hs.service.FetchStints(ctx, sessionKey)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - stints failed %w", err)
	}

	pits, err := hs.service.FetchPits(ctx, sessionKey)
	if err != nil {
		return fmt.Errorf("HistoricalSource.Fetch - pits failed %w", err)
	}

	total := len(raceControls) + len(intervals) + len(positions) + len(laps) + len(pits)
	hs.eventData.EventModels = make([]model.EventModel, 0, total)

	for i := range raceControls {
		hs.eventData.EventModels = append(hs.eventData.EventModels, &raceControls[i])
	}

	for i := range intervals {
		hs.eventData.EventModels = append(hs.eventData.EventModels, &intervals[i])
	}

	for i := range positions {
		hs.eventData.EventModels = append(hs.eventData.EventModels, &positions[i])
	}

	for i := range laps {
		hs.eventData.EventModels = append(hs.eventData.EventModels, &laps[i])
	}

	for i := range pits {
		hs.eventData.EventModels = append(hs.eventData.EventModels, &pits[i])
	}

	hs.raceData.Meeting = &meetings[0]
	hs.raceData.Session = raceSession
	hs.raceData.TotalLaps = getLapCount(raceControls)
	hs.raceData.StartingGrid = grid
	hs.raceData.Drivers = drivers
	hs.raceData.SessionStart = *startTime
	hs.raceData.Stints = stints

	return nil
}

// MARK: helper

func getLapCount(rcs []model.RaceControl) uint {
	for _, rc := range rcs {
		if strings.ToLower(rc.Flag) == "chequered" {
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
	err := fmt.Errorf("couldn't find Session of type: %v and name: %v", st, sn)
	if len(ss) > 0 {
		meetingKey = ss[0].MeetingKey
		err = fmt.Errorf("%v for meeting_key: %v", err, meetingKey)
	}
	return nil, err
}

func (hs *HistoricalSource) getStartTime(rcs []model.RaceControl) (*time.Time, error) {
	for _, rc := range rcs {
		lowerMsg := strings.ToLower(rc.Message)
		if strings.Contains(lowerMsg, "session start") {
			return &rc.DateStart, nil
		}
	}
	return nil, fmt.Errorf("startTime Couldn't be found from RaceControl")
}
