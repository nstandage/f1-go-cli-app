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

type OpenF1HTTP struct{}

var baseUrl string = "https://api.openf1.org/v1"

func (s *OpenF1HTTP) FetchSessions(ctx context.Context, sessionKey string) (*[]model.Session, error) {
	url := fmt.Sprintf("%v/sessions?session_key=%v", baseUrl, sessionKey)
	session, err := fetchData[[]model.Session](ctx, url)
	return session, err
}

func (s *OpenF1HTTP) FetchMeetings(ctx context.Context, sessionKey string) (*[]model.Meeting, error) {
	url := fmt.Sprintf("%v/meetings?session_key=%v", baseUrl, sessionKey)
	meeting, err := fetchData[[]model.Meeting](ctx, url)
	return meeting, err
}

func (s *OpenF1HTTP) FetchDrivers(ctx context.Context, sessionKey string) (*[]model.Driver, error) {
	url := fmt.Sprintf("%v/drivers?session_key=%v", baseUrl, sessionKey)
	drivers, err := fetchData[[]model.Driver](ctx, url)
	return drivers, err
}

func (s *OpenF1HTTP) FetchIntervals(ctx context.Context, sessionKey string) (*[]model.Interval, error) {
	url := fmt.Sprintf("%v/intervals?session_key=%v&interval<0.005", baseUrl, sessionKey)
	interval, err := fetchData[[]model.Interval](ctx, url)
	return interval, err
}

func (s *OpenF1HTTP) FetchLaps(ctx context.Context, sessionKey string) (*[]model.Lap, error) {
	url := fmt.Sprintf("%v/laps?session_key=%v", baseUrl, sessionKey)
	laps, err := fetchData[[]model.Lap](ctx, url)
	return laps, err
}

func (s *OpenF1HTTP) FetchLocations(ctx context.Context, sessionKey string, driverNumber uint) (*[]model.Location, error) {
	url := fmt.Sprintf("%v/location?session_key=%v&driver_number=%v", baseUrl, sessionKey, driverNumber)
	locations, err := fetchData[[]model.Location](ctx, url)
	return locations, err
}

func (s *OpenF1HTTP) FetchPits(ctx context.Context, sessionKey string) (*[]model.Pit, error) {
	url := fmt.Sprintf("%v/pit?session_key=%v", baseUrl, sessionKey)
	pits, err := fetchData[[]model.Pit](ctx, url)
	return pits, err
}

func (s *OpenF1HTTP) FetchPositions(ctx context.Context, sessionKey string) (*[]model.Position, error) {
	url := fmt.Sprintf("%v/pit?session_key=%v", baseUrl, sessionKey)
	positions, err := fetchData[[]model.Position](ctx, url)
	return positions, err
}

func (s *OpenF1HTTP) FetchRaceControls(ctx context.Context, sessionKey string) (*[]model.RaceControl, error) {
	url := fmt.Sprintf("%v/race_control?session_key=%v", baseUrl, sessionKey)
	raceControl, err := fetchData[[]model.RaceControl](ctx, url)
	return raceControl, err
}

func (s *OpenF1HTTP) FetchStint(ctx context.Context, sessionKey string) (*[]model.Stint, error) {
	url := fmt.Sprintf("%v/stints?session_key=%v", baseUrl, sessionKey)
	stints, err := fetchData[[]model.Stint](ctx, url)
	return stints, err
}

func fetchData[T any](ctx context.Context, url string) (*T, error) {
	fmt.Printf("URL: %v\n", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{
		Timeout: time.Second * 10,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

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
