//go:build wasm && js

package main

import (
	"context"
	"encoding/json"
	"log/slog"

	"g.tizu.dev/CCWSUI/components"
	"g.tizu.dev/CCWSUI/web/webmsg"
)

func (j *jsapi) startReceiveLoop() {
	for {
		_, msg, err := j.ws.Read(context.Background())
		if err != nil {
			slog.Error("WS read error", "err", err)
			return
		}

		var e webmsg.Envelope
		if err := json.Unmarshal(msg, &e); err != nil {
			slog.Error("Failed to unmarshal message", "err", err)
			continue
		}

		switch e.Type {
		case webmsg.TypeUpdate:
			var u webmsg.Update
			if err := json.Unmarshal(e.Data, &u); err != nil {
				slog.Error("Failed to unmarshal update", "err", err)
				continue
			}
			newroot, err := components.FromWire(u.Root)
			if err != nil {
				slog.Error("Failed to unmarshal root", "err", err)
				continue
			}
			j.Root = newroot
			j.TotalRerender("Server-sent Tree Update")

		case webmsg.TypeTexture:
			var t webmsg.Texture
			if err := json.Unmarshal(e.Data, &t); err != nil {
				slog.Error("Failed to unmarshal texture", "err", err)
				continue
			}
			j.wrap.Get("userTextures").Set(t.ID, t.Data)
			j.TotalRerender("Server-sent Texture Update")
		}
	}
}

func (j *jsapi) SendEvent(id string, v any) {
	if err := webmsg.SendMsg(j.ws, webmsg.TypeEvent,
		webmsg.Event{ID: id, Event: mustMarshal(v)}); err != nil {
		slog.Error("Failed to send event", "err", err)
	}
}

func mustMarshal(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
