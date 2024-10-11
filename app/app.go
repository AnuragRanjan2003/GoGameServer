package app

import (
	"context"
	"net/http"

	"example.com/main/server"
)

type App struct {
	server *server.Server
}

func NewApp(ctx context.Context) *App {
	return &App{
		server: server.NewServer(ctx),
	}
}

func (a *App) Start(port string) error {
	http.HandleFunc("/join", setCSP(a.server.HandleWSConnection))
	http.HandleFunc("/info", a.server.HandleInfoRequest)
	err := http.ListenAndServe(port, nil)
	return err
}

func setCSP(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Content-Security-Policy header
		w.Header().Set("Content-Security-Policy", "default-src 'self'; connect-src https: ws://localhost:3000 wss://localhost:3000")

		// Call the next handler
		next(w, r)
	})
}
