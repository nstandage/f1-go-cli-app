package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type DataSource interface {
	Start(ctx context.Context, out chan<- string)
}

type HistoricalSource struct{}

func (hs *HistoricalSource) Start(ctx context.Context, out chan<- string) {
	res, err := http.Get("https://api.openf1.org/v1/drivers?driver_number=1&session_key=9158")
	if err != nil {
		fmt.Printf("Error 1: %v", err)
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error 2: %v", err)
	}

	out <- string(data)
	out <- "Test String"
	close(out)
}
