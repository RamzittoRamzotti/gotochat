package app

import (
	"github.com/RamzittoRamzotti/gotochat.git/internal/domain"
	"github.com/RamzittoRamzotti/gotochat.git/internal/handler"
	"github.com/RamzittoRamzotti/gotochat.git/internal/lib/middleware"

	chi "github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func New() *domain.App {
	return &domain.App{
		Rooms: make(map[string]*domain.Room),
	}
}

func Run(app *domain.App) {
	for i := 0; i < 10; i++ {
		roomID := "room" + string(i)
		room := domain.NewRoom(roomID)
		app.Rooms[roomID] = room
		go room.Run()
	}
	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Use(chiMiddleware.AllowContentType("application/json"))
		r.Use(middleware.AuthMiddleware)

	})
	r.Post("/register", handler.Register())
}
