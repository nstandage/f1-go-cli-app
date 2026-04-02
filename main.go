package main

import (
	"context"
	"fmt"

	datasources "github.com/nstandage/f1-go-cli-app/datasource"
	models "github.com/nstandage/f1-go-cli-app/model"
	"github.com/nstandage/f1-go-cli-app/service"
)

func main() {

	var sessionKey = "11253"
	service := service.OpenF1Service{}

	hs := datasources.HistoricalSource{
		Service: &service,
	}

	sessionData, err := hs.Start(context.Background(), sessionKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	c := make(chan models.Event)
	hs.Load(context.Background(), sessionData, c)

}
