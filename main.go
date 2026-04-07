package main

import (
	"context"
	"fmt"

	"github.com/nstandage/f1-go-cli-app/aggregator"
	"github.com/nstandage/f1-go-cli-app/datasource"
	"github.com/nstandage/f1-go-cli-app/model"
	"github.com/nstandage/f1-go-cli-app/service"
)

func main() {

	var sessionKey = "11253"
	service := service.OpenF1Service{}

	hs := datasource.HistoricalSource{
		Service: &service,
	}

	sessionData, err := hs.Start(context.Background(), sessionKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := make(chan *model.Event)
	eng := datasource.ReplayEngine{
		SessionData: sessionData,
		Channel:     c,
	}

	ag := aggregator.Engine{
		Channel: c,
	}

	go hs.Load(&eng)
	for event := range c {
		ag.Handle(event)
	}
}
