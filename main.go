package main

import (
	"context"
	"log"

	"github.com/nstandage/f1-go-cli-app/aggregator"
	"github.com/nstandage/f1-go-cli-app/datasource"
	"github.com/nstandage/f1-go-cli-app/model"
	"github.com/nstandage/f1-go-cli-app/service"
)

func main() {

	var sessionKey = "11253"
	service := service.OpenF1HTTP{}

	hs := datasource.HistoricalSource{
		Service: &service,
	}

	raceData, eventData, err := hs.Fetch(context.Background(), sessionKey)
	if err != nil {
		log.Fatal(err)
	}

	replay := datasource.ReplayEngine{EventData: eventData}
	ag := aggregator.Engine{RaceData: raceData}

	c := make(chan *model.Event)
	go ag.Start(c)
	replay.Start(c)
}
