package app

import (
	"net/http"

	"example.com/main/server"
)

type App struct {
	server *server.Server
}

func NewApp() *App {
	return &App{
		server: server.NewServer(),
	}
}

func (a *App) Start() error {
	http.HandleFunc("/join", a.server.HandleWSConnection)
	err := http.ListenAndServe(":3000", nil)
	return err
}
