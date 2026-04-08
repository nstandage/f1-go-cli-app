package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nstandage/f1-go-cli-app/tui"
)

func main() {

	// var sessionKey = "11253"
	// service := service.OpenF1HTTP{}

	// hs := datasource.HistoricalSource{
	// 	Service: &service,
	// }

	// raceData, eventData, err := hs.Fetch(context.Background(), sessionKey)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// replay := datasource.ReplayEngine{EventData: eventData}
	// ag := aggregator.Engine{RaceData: raceData}

	// c := make(chan *model.Event)
	// go ag.Start(c)
	// replay.Start(c)

	quesions := []tui.Question{
		tui.NewQuestion("What is your name?"),
		tui.NewQuestion("What is your favorite editor?"),
	}
	m := tui.New(quesions)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
