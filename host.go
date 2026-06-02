package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"g.tizu.dev/CCWSUI/components"
	"github.com/coder/websocket"
	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (app *CCWSUI) handleHost(w http.ResponseWriter, r *http.Request) error {
	conn, err := websocket.Accept(w, r,
		&websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		return err
	}
	defer conn.Close(websocket.StatusInternalError, "Internal Server Error")

	room := gonanoid.Must(6)
	roominst, err := NewRoom(conn, room, "Home",
		components.LiteralOf("Waiting for host..."))
	if err != nil {
		return err
	}
	app.roomsmu.Lock()
	app.rooms[room] = roominst
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

// Host to Server
const (
	HostMessageTypeUpdate HostMessageType = iota + 1
)

// Server to Host
const (
	HostMessageTypeReady HostMessageType = iota + 1
	HostMessageTypeHello
)

type HostMessageEnvelope struct {
	Type HostMessageType `json:"t"`
	Data json.RawMessage `json:"d"`
}

type HostMessageUpdate struct {
	Client uuid.UUID           `json:"client"`
	Root   components.WireNode `json:"root"`
}

type HostMessageReady struct {
	URL string `json:"url"`
}

type HostMessageHello struct {
	Client uuid.UUID `json:"client"`
}

func (room *Room) handleHostMessage(conn *websocket.Conn, msg HostMessageEnvelope) error {
	switch msg.Type {
	case HostMessageTypeUpdate:
		var data HostMessageUpdate
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			return err
		}
		if c, ok := room.Get(data.Client); ok {
			if err := c.Update(data.Root); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("invalid type")
	}
	return nil
}
