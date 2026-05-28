package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"github.com/nstandage/f1-go-cli-app/model"
	"golang.org/x/time/rate"
)

type OpenF1HTTP struct{
	client *http.Client
	baseURL string
	rateLimiter *rate.Limiter
}

var baseURL string = "https://api.openf1.org/v1"

func NewOpenF1HTTP() *OpenF1HTTP {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	return &OpenF1HTTP{
		client: &client,
		rateLimiter: rate.NewLimiter(rate.Limit(2.9), 1),
	}
}

func (s *OpenF1HTTP) FetchSessions(ctx context.Context, sessionKey string) ([]model.Session, error) {
	reader,err := s.fetchData(ctx, fmt.Sprintf("/sessions?session_key=%s", sessionKey))
	if err != nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchSessions sessions == nil %w", err)
	}
	return parseJSON[[]model.Session](reader)
}

func (s *OpenF1HTTP) FetchSessionsByMeeting(ctx context.Context, meetingKey string) ([]model.Session, error) {
	reader,err := s.fetchData(ctx, fmt.Sprintf("/sessions?meeting_key=%s", meetingKey))
	if err != nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchMeetingSessions sessions == nil %w", err)
	}
	defer reader.Close()
	return parseJSON[[]model.Session](reader)
}

func (s *OpenF1HTTP) FetchMeetings(ctx context.Context, meetingKey string) ([]model.Meeting, error) {
	reader,err := s.fetchData(ctx, fmt.Sprintf("/meetings?meeting_key=%s", meetingKey))
	if err != nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchMeetings meetings == nil %w", err)
	}
	defer reader.Close()
	return parseJSON[[]model.Meeting](reader)
}

func (s *OpenF1HTTP) FetchDrivers(ctx context.Context, sessionKey string) ([]model.Driver, error) {
	reader,err := s.fetchData(ctx, fmt.Sprintf("/drivers?session_key=%s", sessionKey))
	if err != nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchDrivers drivers == nil %w", err)
	}
	defer reader.Close()
	return parseJSON[[]model.Driver](reader)
}

func (s *OpenF1HTTP) FetchIntervals(ctx context.Context, sessionKey string) ([]model.Interval, error) {
	reader,err := s.fetchData(ctx, fmt.Sprintf("/intervals?session_key=%s", sessionKey))
	if err != nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchIntervals ints == nil %w", err)
	}
	defer reader.Close()
	return parseJSON[[]model.Interval](reader)
}

func (s *OpenF1HTTP) FetchLaps(ctx context.Context, sessionKey string) ([]model.Lap, error) {
	reader,err := s.fetchData(ctx, fmt.Sprintf("/laps?session_key=%s", sessionKey))
	if err != nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchLaps laps == nil %w", err)
	}
	defer reader.Close()
	return parseJSON[[]model.Lap](reader)
}

func (s *OpenF1HTTP) FetchLocations(ctx context.Context, sessionKey string, driverNumber uint) ([]model.Location, error) {
	reader,err := s.fetchData(ctx, fmt.Sprintf("/location?session_key=%s&driver_number=%v", sessionKey, driverNumber))
	if err != nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchLocations loc == nil %w", err)
	}
	defer reader.Close()
	return parseJSON[[]model.Location](reader)
}

func (s *OpenF1HTTP) FetchPits(ctx context.Context, sessionKey string) ([]model.Pit, error) {
	reader,err := s.fetchData(ctx, fmt.Sprintf("/pit?session_key=%s", sessionKey))
	if err != nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchPits pits == nil %w", err)
	}
	defer reader.Close()
	return parseJSON[[]model.Pit](reader)
}

func (s *OpenF1HTTP) FetchPositions(ctx context.Context, sessionKey string) ([]model.Position, error) {
	reader,err := s.fetchData(ctx, fmt.Sprintf("/position?session_key=%s", sessionKey))
	if err != nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchPositions pos == nil %w", err)
	}
	defer reader.Close()
	return parseJSON[[]model.Position](reader)
}

func (s *OpenF1HTTP) FetchRaceControls(ctx context.Context, sessionKey string) ([]model.RaceControl, error) {
	reader, err := s.fetchData(ctx, fmt.Sprintf("/race_control?session_key=%v", sessionKey))
	if err != nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchRaceControl rc == nil %w", err)
	}
	defer reader.Close()
	return parseJSON[[]model.RaceControl](reader)
}

func (s *OpenF1HTTP) FetchStints(ctx context.Context, sessionKey string) ([]model.Stint, error) {
	reader, err := s.fetchData(ctx, fmt.Sprintf("/stints?session_key=%v", sessionKey))
	if err != nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchStints stints == nil %w", err)
	}
	defer reader.Close()
	return parseJSON[[]model.Stint](reader)
}

func (s *OpenF1HTTP) FetchStartingGrid(ctx context.Context, sessionKey string) ([]model.StartingGrid, error) {
	reader, err := s.fetchData(ctx, fmt.Sprintf("/starting_grid?session_key=%s", sessionKey))
	if err != nil {
		return nil, fmt.Errorf("OpenF1HTTP.FetchStartingGrid startingGrids == nil %w", err)
	}
	defer reader.Close()
	return parseJSON[[]model.StartingGrid](reader)
}

func (s *OpenF1HTTP) fetchData(ctx context.Context, endpoint string) (io.ReadCloser, error) {
	if err := s.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	fullURL := fmt.Sprintf("%s%s", baseURL, endpoint)
	fmt.Printf("URL: %v\n", fullURL)


	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status code: %d, endpoint: %v", res.StatusCode, endpoint)
		defer res.Body.Close()
		return nil, err
	}
	return res.Body, nil
}

func parseJSON[T any](r io.Reader) (T, error) {
	var result T
	if err := json.NewDecoder(r).Decode(&result); err != nil {
		return result, fmt.Errorf("json decode failed for type %T: %w", result, err)
	}
	return result, nil
}
