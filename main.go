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
	var sessionKey = "9939"
	var meetingKey = "1265"
	service := service.OpenF1HTTP{}

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hs := datasource.HistoricalSource{
		Service: &service,
	}

	raceData, eventData, err := hs.Fetch(context.Background(), sessionKey, meetingKey)
	if err != nil {
		log.Fatal(err)
	}

	engine := datasource.ReplayEngine{EventData: eventData}
	datasource := aggregator.Datasource{
		Meeting:      raceData.Meeting,
		Session:      raceData.Session,
		TotalLaps:    raceData.TotalLaps,
		IsReplay:     engine.IsReplay(),
		StartingGrid: raceData.StartingGrid,
		Drivers:      aggregator.ConvertDrivers(raceData.Drivers),
	}
	datasource.AddStartingGrid()
	ag := aggregator.Engine{Datasource: &datasource}

	p := tea.NewProgram(tui.Model{
		Window: tui.Window{},
		Engine: &ag,
	})

	c := make(chan *model.Event)
	go ag.Start(c)
	go engine.Start(c)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
