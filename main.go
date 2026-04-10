package main

import (
	"context"
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/nstandage/f1-go-cli-app/aggregator"
	"github.com/nstandage/f1-go-cli-app/datasource"
	"github.com/nstandage/f1-go-cli-app/model"
	"github.com/nstandage/f1-go-cli-app/service"
	"github.com/nstandage/f1-go-cli-app/tui"
)

func main() {

	var sessionKey = "11253"
	service := service.OpenF1HTTP{}

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hs := datasource.HistoricalSource{
		Service: &service,
	}

	raceData, eventData, err := hs.Fetch(context.Background(), sessionKey)
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(tui.Model{})

	replay := datasource.ReplayEngine{EventData: eventData}
	ag := aggregator.Engine{RaceData: raceData, Program: p}

	c := make(chan *model.Event)
	go ag.Start(c)
	go replay.Start(c)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
