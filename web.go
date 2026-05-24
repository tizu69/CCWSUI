package main

import (
	"embed"
	"log/slog"
	"net/http"

	"g.tizu.dev/CCWSUI/components"
	"github.com/coder/websocket"
)

//go:embed static/*
var staticFS embed.FS

func (app *CCWSUI) Run() {
	mux := http.NewServeMux()
	mux.Handle("/{$}", Handler(app.handleRoot))
	mux.Handle("GET /r/{room}", Handler(app.handleRoom))
	mux.Handle("GET /r/{room}/service", Handler(app.handleRoomService))
	mux.Handle("GET /host", Handler(app.handleHost))
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
		slog.Error("Failed to handle request", "method", r.Method, "path", r.URL.Path, "err", err)
	}
}

func (app *CCWSUI) handleRoot(w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, "/r/home", http.StatusPermanentRedirect)
	return nil
}

func (app *CCWSUI) handleRoom(w http.ResponseWriter, r *http.Request) error {
	roomid := r.PathValue("room")
	room := app.room(roomid)
	if room == nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return nil
	}

	return components.Root(roomid, room.Title, room.Root).Render(w)
}

func (app *CCWSUI) handleRoomService(w http.ResponseWriter, r *http.Request) error {
	room := app.room(r.PathValue("room"))
	if room == nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return nil
	}

	conn, err := websocket.Accept(w, r,
		&websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		return err
	}
	defer conn.Close(websocket.StatusInternalError, "Internal Server Error")

	room.Add(conn)
	defer room.Remove(conn)

	return nil
}
