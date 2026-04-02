package main

import (
	"context"
	"fmt"

	"github.com/nstandage/f1-go-cli-app/service"
)

func main() {

	var sessionKey = "11253"
	service := service.OpenF1Service{}

	// hs := datasource.HistoricalSource{
	// 	Service: &service,
	// }

	// _, err := hs.Start(sessionKey, context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	data, err := service.FetchStint(context.Background(), sessionKey)

	if err != nil {
		fmt.Printf("!!! ERROR: %v", err)
		return
	}

	fmt.Printf("!!! type: %T", data)
	// fmt.Printf("!!! Data: %v", data)
}
