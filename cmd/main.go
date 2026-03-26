package main

import (
	"log/slog"
	"os"

	"github.com/RamzittoRamzotti/gotochat.git/internal/app"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	app := app.New()

}
