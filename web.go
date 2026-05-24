package main

import (
	"embed"
	"log/slog"
	"net/http"

	"g.tizu.dev/CCWSUI/components"
)

//go:embed static/*
var staticFS embed.FS

func (app *CCWSUI) Run() {
	mux := http.NewServeMux()
	mux.Handle("GET /{$}", Handler(app.handleRoot))
	mux.Handle("GET /static/", http.FileServer(http.FS(staticFS)))

	slog.Info("Listening!", "addr", app.addr)
	if err := http.ListenAndServe(app.addr, mux); err != nil {
		slog.Error("Failed to start server", "err", err)
	}
}

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *CCWSUI) handleRoot(w http.ResponseWriter, r *http.Request) error {
	return components.Root("hi").Render(w)
}
