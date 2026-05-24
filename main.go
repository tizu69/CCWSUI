package main

import (
	"log/slog"
	"os"
	"sync"

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
		Multiline: true, TimeFormat: "01.02.06 3:04PM",
	})))

	godotenv.Load()

	addr := os.Getenv("CCWSUI_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	app := CCWSUI{
		addr:  addr,
		rooms: getCoreRooms(),
	}
	app.Run()
}

func (app *CCWSUI) room(id string) *Room {
	app.roomsmu.RLock()
	defer app.roomsmu.RUnlock()
	return app.rooms[id]
}
