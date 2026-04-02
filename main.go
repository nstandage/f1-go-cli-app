package main

import (
	"context"
	"fmt"

	"github.com/nstandage/f1-go-cli-app/service"
)

func main() {
	hs := service.HistoricalSource{}
	channel := make(chan string)

	go hs.Start(context.Background(), channel)
	fmt.Println("1")
	for msg := range channel {
		fmt.Println("2")
		fmt.Println(msg)
	}

	fmt.Println("3")
	// close(channel)
}
