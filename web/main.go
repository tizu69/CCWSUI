//go:build wasm && js

package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Marlliton/slogpretty"
	"github.com/coder/websocket"
)

func main() {
	slog.SetDefault(slog.New(slogpretty.New(os.Stdout, &slogpretty.Options{
		Level: slog.LevelDebug, AddSource: true, Colorful: true,
		Multiline: true, TimeFormat: "01.02.06 3:04PM",
	})))

	app := NewJSAPI()

	slog.Info("Connecting to gateway!", "url", app.SocketURL())
	ws, _, err := websocket.Dial(context.Background(), app.SocketURL(), &websocket.DialOptions{})
	if err != nil {
		panic(err)
	}
	ws.SetReadLimit(1024 * 1024)
	app.ws = ws

	go app.startReceiveLoop()

	select {}
}
