package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"g.tizu.dev/CCWSUI/components"
	"g.tizu.dev/CCWSUI/web/webmsg"
	"github.com/coder/websocket"
	"github.com/google/uuid"
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
		if r.Header.Get("Upgrade") != "websocket" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		slog.Error("Failed to handle request", "method", r.Method, "path", r.URL.Path, "err", err)
	}
}

func (app *CCWSUI) handleRoot(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Upgrade") == "websocket" {
		return app.handleHost(w, r)
	}
	http.Redirect(w, r, "/r/home", http.StatusPermanentRedirect)
	return nil
}

func (app *CCWSUI) handleRoom(w http.ResponseWriter, r *http.Request) error {
	room := app.room(r.PathValue("room"))
	if room == nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return nil
	}

	return components.Root(fmt.Sprintf("/r/%s/service",
		url.PathEscape(room.ID)), room.Title, w)
}

func (app *CCWSUI) handleRoomService(w http.ResponseWriter, r *http.Request) error {
	room := app.room(r.PathValue("room"))
	if room == nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := websocket.Accept(w, r,
		&websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		return err
	}
	defer conn.CloseNow()

	userid, err := uuid.Parse(r.URL.Query().Get("user"))
	if err != nil {
		return err
	}

	id, err := room.Add(conn, userid)
	if err != nil {
		return err
	}
	defer room.Remove(id)

	for {
		_, msg, err := conn.Read(ctx)
		if err != nil {
			break
		}

		var e webmsg.Envelope
		if err := json.Unmarshal(msg, &e); err != nil {
			return err
		}

		switch e.Type {
		case webmsg.TypeEvent:
			var event webmsg.Event
			if err := json.Unmarshal(e.Data, &event); err != nil {
				return err
			}
			room.HandleEvent(id, event)
		default:
			return err
		}
	}

	return nil
}
