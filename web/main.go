//go:build wasm && js

package main

import (
	"log/slog"
	"os"

	"github.com/Marlliton/slogpretty"
)

func main() {
	slog.SetDefault(slog.New(slogpretty.New(os.Stdout, &slogpretty.Options{
		Level: slog.LevelDebug, AddSource: true, Colorful: true,
		Multiline: true, TimeFormat: "02.01.06 3:04PM",
	})))

	app := NewJSAPI()
	go app.run()
	select {}
}
