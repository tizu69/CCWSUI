package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"g.tizu.dev/CCWSUI/components"
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

func (app *CCWSUI) handleHost(w http.ResponseWriter, r *http.Request) error {
	conn, err := websocket.Accept(w, r,
		&websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		return err
	}
	defer conn.Close(websocket.StatusInternalError, "Internal Server Error")

	room := uuid.New().String()
	app.roomsmu.Lock()
	app.rooms[room] = NewRoom(room, "Home",
		components.LiteralOf("Waiting for host..."))
	app.roomsmu.Unlock()

	for {
		_, b, err := conn.Read(r.Context())
		if err != nil {
			break
		}

		var msg HostMessageEnvelope
		if err := json.Unmarshal(b, &msg); err != nil {
			return conn.Close(websocket.StatusUnsupportedData, "Bad JSON")
		}

		app.roomsmu.RLock()
		if err := app.rooms[room].handleHostMessage(conn, msg); err != nil {
			return conn.Close(websocket.StatusUnsupportedData, "Invalid Message")
		}
		app.roomsmu.RUnlock()
	}

	return nil
}

type HostMessageType int

const (
	HostMessageTypeTitle HostMessageType = iota + 1
)

type HostMessageEnvelope struct {
	Type HostMessageType `json:"t"`
	Data json.RawMessage `json:"d"`
}

type HostMessageTitle struct {
	Title string `json:"title"`
}

func (room *Room) handleHostMessage(conn *websocket.Conn, msg HostMessageEnvelope) error {
	switch msg.Type {
	case HostMessageTypeTitle:
		var data HostMessageTitle
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			return err
		}
		room.Title = data.Title
	default:
		return fmt.Errorf("invalid type")
	}
	return nil
}
