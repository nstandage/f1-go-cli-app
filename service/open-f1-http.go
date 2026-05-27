package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nstandage/f1-go-cli-app/model"
)

type OpenF1HTTP struct{
	client *http.Client
	baseUrl string
}

var baseURL string = "https://api.openf1.org/v1"

func NewOpenF1HTTP() OpenF1HTTP {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	return OpenF1HTTP{
		client: &client,
	}
}

func (s *OpenF1HTTP) FetchSessions(ctx context.Context, sessionKey string) ([]model.Session, error) {
	url := fmt.Sprintf("%v/sessions?session_key=%v", baseURL, sessionKey)
	session, err := fetchData[[]model.Session](ctx, s.client, url)
	if session == nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchSessions sessions == nil %w", err)
	}
	return *session, err
}

func (s *OpenF1HTTP) FetchSessionsByMeeting(ctx context.Context, meetingKey string) ([]model.Session, error) {
	url := fmt.Sprintf("%v/sessions?meeting_key=%v", baseURL, meetingKey)
	session, err := fetchData[[]model.Session](ctx, s.client, url)
	if session == nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchMeetingSessions sessions == nil %w", err)
	}
	return *session, err
}

func (s *OpenF1HTTP) FetchMeetings(ctx context.Context, meetingKey string) ([]model.Meeting, error) {
	url := fmt.Sprintf("%v/meetings?meeting_key=%v", baseURL, meetingKey)
	meeting, err := fetchData[[]model.Meeting](ctx, s.client, url)
	if meeting == nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchMeetings meetings == nil %w", err)
	}
	return *meeting, err
}

func (s *OpenF1HTTP) FetchDrivers(ctx context.Context, sessionKey string) ([]model.Driver, error) {
	url := fmt.Sprintf("%v/drivers?session_key=%v", baseURL, sessionKey)
	drivers, err := fetchData[[]model.Driver](ctx, s.client, url)
	if drivers == nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchDrivers drivers == nil %w", err)
	}
	return *drivers, err
}

func (s *OpenF1HTTP) FetchIntervals(ctx context.Context, sessionKey string) ([]model.Interval, error) {
	url := fmt.Sprintf("%v/intervals?session_key=%v", baseURL, sessionKey)
	intervals, err := fetchData[[]model.Interval](ctx, s.client, url)
	if intervals == nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchIntervals ints == nil %w", err)
	}
	return *intervals, err
}

func (s *OpenF1HTTP) FetchLaps(ctx context.Context, sessionKey string) ([]model.Lap, error) {
	url := fmt.Sprintf("%v/laps?session_key=%v", baseURL, sessionKey)
	laps, err := fetchData[[]model.Lap](ctx, s.client, url)
	if laps == nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchLaps laps == nil %w", err)
	}
	return *laps, err
}

func (s *OpenF1HTTP) FetchLocations(ctx context.Context, sessionKey string, driverNumber uint) ([]model.Location, error) {
	url := fmt.Sprintf("%v/location?session_key=%v&driver_number=%v", baseURL, sessionKey, driverNumber)
	locations, err := fetchData[[]model.Location](ctx, s.client, url)
	if locations == nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchLocations loc == nil %w", err)
	}
	return *locations, err
}

func (s *OpenF1HTTP) FetchPits(ctx context.Context, sessionKey string) ([]model.Pit, error) {
	url := fmt.Sprintf("%v/pit?session_key=%v", baseURL, sessionKey)
	pits, err := fetchData[[]model.Pit](ctx, s.client, url)
	if pits == nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchPits pits == nil %w", err)
	}
	return *pits, err
}

func (s *OpenF1HTTP) FetchPositions(ctx context.Context, sessionKey string) ([]model.Position, error) {
	url := fmt.Sprintf("%v/position?session_key=%v", baseURL, sessionKey)
	positions, err := fetchData[[]model.Position](ctx, s.client, url)
	if positions == nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchPositions pos == nil %w", err)
	}
	return *positions, err
}

func (s *OpenF1HTTP) FetchRaceControls(ctx context.Context, sessionKey string) ([]model.RaceControl, error) {
	url := fmt.Sprintf("%v/race_control?session_key=%v", baseURL, sessionKey)
	raceControl, err := fetchData[[]model.RaceControl](ctx, s.client, url)
	if raceControl == nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchRaceControl rc == nil %w", err)
	}
	return *raceControl, err
}

func (s *OpenF1HTTP) FetchStints(ctx context.Context, sessionKey string) ([]model.Stint, error) {
	url := fmt.Sprintf("%v/stints?session_key=%v", baseURL, sessionKey)
	stints, err := fetchData[[]model.Stint](ctx, s.client, url)
	if stints == nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchStints stints == nil %w", err)
	}
	return *stints, err
}

func (s *OpenF1HTTP) FetchStartingGrid(ctx context.Context, sessionKey string) ([]model.StartingGrid, error) {
	url := fmt.Sprintf("%v/starting_grid?session_key=%v", baseURL, sessionKey)
	grid, err := fetchData[[]model.StartingGrid](ctx, s.client, url)
	if grid == nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchStartingGrid startingGrids == nil %w", err)
	}
	return *grid, err
}

func fetchData[T any](ctx context.Context, client *http.Client, url string) (*T, error) {
	fmt.Printf("URL: %v\n", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		err = fmt.Errorf("Unexpected status code: %d, url: %v", res.StatusCode, url)
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil
}
