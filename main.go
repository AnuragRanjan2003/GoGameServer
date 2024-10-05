package main

import (
	"fmt"
	"log"

	"example.com/main/app"
)

func main() {
	fmt.Println("starting app")
	app := app.NewApp()
	err := app.Start()
	if err != nil {
		log.Fatal("server start failed: ", err)
	}
	fmt.Println("server starte on 3000")

}
