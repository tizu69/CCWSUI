//go:build wasm && js

package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"g.tizu.dev/CCWSUI/components"
	"g.tizu.dev/CCWSUI/web/webmsg"
	"github.com/coder/websocket"
)

func (j *jsapi) startReceiveLoop() {
	for {
		_, msg, err := j.ws.Read(context.Background())
		if err != nil {
			slog.Error("WS read error", "err", err)
			return
		}

		if len(msg) == 0 {
			continue // ping
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

		case webmsg.TypeMetadata:
			var d webmsg.Metadata
			if err := json.Unmarshal(e.Data, &d); err != nil {
				slog.Error("Failed to unmarshal metadata", "err", err)
				continue
			}
			j.wrap.Call("updateMetadata", string(mustMarshal(d)))

		case webmsg.TypeTexture:
			var t webmsg.Texture
			if err := json.Unmarshal(e.Data, &t); err != nil {
				slog.Error("Failed to unmarshal texture", "err", err)
				continue
			}
			j.wrap.Get("userTextures").Set(t.ID, base64.StdEncoding.EncodeToString(t.Data))
			j.wrap.Set("textures", map[string]any{})
			j.wrap.Set("textureBorders", map[string]any{})
			j.TotalRerender("Server-sent Texture Update")

		case webmsg.TypeRedirect:
			var d webmsg.Redirect
			if err := json.Unmarshal(e.Data, &d); err != nil {
				slog.Error("Failed to unmarshal redirect", "err", err)
				continue
			}
			j.wrap.Call("redirect", d.URL)
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

func (j *jsapi) run() {
	const maxBackoff = 30 * time.Second
	backoff := time.Duration(0)
	j.wrap.Call("setConnectionStatus", "connecting")
	for {
		j.Root = nil
		j.context = make(map[string]any)
		j.texBordersCache = make(map[string][4]int)

		if r, _ := http.Get(j.ValidateSocketURL()); r != nil && r.StatusCode == 404 {
			j.wrap.Call("setConnectionStatus", "notfound")
		}

		ws, _, err := websocket.Dial(context.Background(), j.SocketURL(), &websocket.DialOptions{})
		if err != nil {
			slog.Error("Failed to connect to gateway", "err", err)
			if backoff == 0 {
				backoff = time.Second
			}
			time.Sleep(backoff)
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}
		ws.SetReadLimit(1024 * 1024)

		old := j.ws
		j.ws = ws
		if old != nil {
			old.Close(websocket.StatusNormalClosure, "reconnecting")
		}

		j.wrap.Call("setConnectionStatus", "connected")
		backoff = 0
		slog.Info("Connected to gateway!", "url", j.SocketURL())
		j.startReceiveLoop()
		slog.Warn("Gateway disconnected, reconnecting...")
		j.wrap.Call("setConnectionStatus", "reconnecting")
	}
}
