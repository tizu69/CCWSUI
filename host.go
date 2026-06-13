package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/coder/websocket"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (app *CCWSUI) handleHost(w http.ResponseWriter, r *http.Request) error {
	conn, err := websocket.Accept(w, r,
		&websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		return err
	}
	defer conn.Close(websocket.StatusInternalError, "Internal Server Error")

	room := NewRemoteRoom(conn)
	defer func() {
		if room.ID != "" {
			app.roomsmu.Lock()
			delete(app.rooms, room.ID)
			for _, c := range room.clients {
				c.Close()
			}
			app.roomsmu.Unlock()
		}
	}()

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			if err := conn.Write(ctx, websocket.MessageText, []byte{}); err != nil {
				return
			}
			time.Sleep(20 * time.Second)
		}
	}()

	for {
		_, b, err := conn.Read(r.Context())
		if err != nil {
			return err
		}

		var env HostEnvelope
		time.Sleep(time.Duration(app.latency) * time.Millisecond)
		if err := json.Unmarshal(b, &env); err != nil {
			conn.Close(websocket.StatusUnsupportedData, "Bad JSON")
			return err
		}

		if err := app.handleHostMsg(room, env); err != nil {
			conn.Close(websocket.StatusUnsupportedData, err.Error())
			return err
		}
	}
}

func (app *CCWSUI) handleHostMsg(room *Room, env HostEnvelope) error {
	switch env.Type {
	case HostMsgUpdate:
		if !room.Frozen() {
			return fmt.Errorf("room not yet frozen")
		}
		var data HostUpdatePayload
		if err := json.Unmarshal(env.Data, &data); err != nil {
			return err
		}
		c, ok := room.Get(data.Client)
		if !ok {
			return nil
		}
		return c.Update(data.Root)

	case HostMsgMetadata:
		if !room.Frozen() {
			return fmt.Errorf("room not yet frozen")
		}
		var data HostMetadataPayload
		if err := json.Unmarshal(env.Data, &data); err != nil {
			return err
		}
		c, ok := room.Get(data.Client)
		if !ok {
			return nil
		}
		return c.UpdateMetadata(data)

	case HostMsgWantSlug:
		if room.Frozen() {
			return fmt.Errorf("room is frozen")
		}
		var data HostWantSlugPayload
		if err := json.Unmarshal(env.Data, &data); err != nil {
			return err
		}
		room.wantedSlug = data.Slug
		return nil

	case HostMsgTexture:
		if room.Frozen() {
			return fmt.Errorf("room is frozen")
		}
		var data HostTexturePayload
		if err := json.Unmarshal(env.Data, &data); err != nil {
			return err
		}
		b, err := stripColorProfile(data.Data)
		if err != nil {
			return err
		}
		room.textures[data.ID] = b
		return nil

	case HostMsgFreeze:
		if room.Frozen() {
			return fmt.Errorf("room already frozen")
		}
		room.frozen = true

		id := room.wantedSlug
		for id == "" || app.slugTaken(id) {
			id = gonanoid.Must(6)
		}
		room.ID = id

		app.roomsmu.Lock()
		app.rooms[id] = room
		app.roomsmu.Unlock()

		slog.Info("Room frozen", "id", id)
		return sendHostMsg(room.hostConn, HostMsgReady, HostReadyPayload{
			URL: room.UserURL(),
		})

	default:
		return fmt.Errorf("unknown host message type: %d", env.Type)
	}
}
