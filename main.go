package main

import (
	"log/slog"
	"os"

	"github.com/Marlliton/slogpretty"
	"github.com/joho/godotenv"
)

type CCWSUI struct {
	addr string
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
		addr: addr,
	}
	app.Run()
}
