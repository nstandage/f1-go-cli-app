package main

import (
	"context"
	"fmt"

	datasource "github.com/nstandage/f1-go-cli-app/data-sources"
	"github.com/nstandage/f1-go-cli-app/service"
)

func main() {

	var sessionKey uint = 11253
	service := service.OpenF1Service{}

	hs := datasource.HistoricalSource{
		Service: &service,
	}

	_, err := hs.Start(sessionKey, context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\nDone")

	// go hs.Start(context.Background(), channel)
	// fmt.Println("1")
	// for msg := range channel {
	// 	fmt.Println("2")
	// 	fmt.Println(msg)
	// }

	// fmt.Println("3")
	// close(channel)
}
