package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/RamzittoRamzotti/gotochat.git/internal/app"
	"github.com/RamzittoRamzotti/gotochat.git/internal/storage/postgres"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:password@localhost:5432/gotodb"
	}

	pool, err := postgres.New(context.Background(), dsn)
	if err != nil {
		logger.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	userStorage := postgres.NewUserStorage(pool)

	a := app.New(userStorage, logger)

	logger.Info("starting server", "addr", ":8080")
	if err := a.Run(":8080"); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
