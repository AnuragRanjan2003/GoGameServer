package main

import (
	"context"
	"example.com/main/app"
	"example.com/main/internal/logs"
	"fmt"
	"log"
)

func main() {
	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()
	fmt.Println("starting app")
	app := app.NewApp(ctx)
	port := ":3000"
	err := app.Start(port)
	if err != nil {
		log.Fatal("server start failed: ", err)
	}
	fmt.Println("server started on ", port)
	logger := logs.NewLogger(ctx)

	logger.Start()

}
