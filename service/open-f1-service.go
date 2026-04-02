package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nstandage/f1-go-cli-app/models"
)

type OpenF1Service struct{}

var baseUrl string = "https://api.openf1.org/v1"

func (s *OpenF1Service) FetchSession(ctx context.Context, key uint) (*[]models.Session, error) {
	url := fmt.Sprintf("%v/sessions?session_key=%v", baseUrl, key)
	session, err := fetchData[[]models.Session](ctx, url)
	return session, err
}

func (s *OpenF1Service) FetchMeeting(ctx context.Context, key uint) (*[]models.Meeting, error) {
	url := fmt.Sprintf("%v/meetings?meeting_key=%v", baseUrl, key)
	meeting, err := fetchData[[]models.Meeting](ctx, url)
	return meeting, err
}

func fetchData[T any](ctx context.Context, url string) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
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
