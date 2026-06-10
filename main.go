package main

import (
	"log/slog"
	"os"
	"sync"

	"g.tizu.dev/CCWSUI/predefined"
	"github.com/Marlliton/slogpretty"
	"github.com/joho/godotenv"
)

type CCWSUI struct {
	addr string

	rooms   map[string]*Room
	roomsmu sync.RWMutex
}

func main() {
	slog.SetDefault(slog.New(slogpretty.New(os.Stdout, &slogpretty.Options{
		Level: slog.LevelDebug, AddSource: true, Colorful: true,
		Multiline: true, TimeFormat: "02.01.06 3:04PM",
	})))

	if files, err := staticFS.ReadDir("static/item"); err != nil || len(files) < 10 {
		slog.Warn("Failed to list item textures, items won't render!", "err", err,
			"1", "Grab IconExporter into a modpack you wish to export from",
			"2", "Load a world",
			"3", "Run '/iconexporter config general.fileNameHashComponents false'",
			"4", "Run '/iconexporter export 256'",
			"5", "Run 'rm *\\{*' to remove NBT files",
			"6", "Place the PNGs into static/item")
	}

	godotenv.Load()
	addr := os.Getenv("CCWSUI_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	app := CCWSUI{
		addr:  addr,
		rooms: make(map[string]*Room),
	}

	app.rooms["home"] = NewPredefinedRoom(predefined.NewHome(), "home")
	app.rooms["docs"] = NewPredefinedRoom(predefined.NewDocs(), "docs")

	app.Run()
}

func (app *CCWSUI) room(id string) *Room {
	app.roomsmu.RLock()
	defer app.roomsmu.RUnlock()
	return app.rooms[id]
}

func (app *CCWSUI) slugTaken(slug string) bool {
	app.roomsmu.RLock()
	defer app.roomsmu.RUnlock()
	_, ok := app.rooms[slug]
	return ok
}
