package app

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/RamzittoRamzotti/gotochat.git/internal/domain"
	"github.com/RamzittoRamzotti/gotochat.git/internal/handler"
	"github.com/RamzittoRamzotti/gotochat.git/internal/lib/middleware"
	"github.com/RamzittoRamzotti/gotochat.git/internal/metrics"
	auth "github.com/RamzittoRamzotti/gotochat.git/internal/services/auth"
	chat "github.com/RamzittoRamzotti/gotochat.git/internal/services/chat"
	chi "github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

type App struct {
	authService *auth.AuthService
	chatService *chat.ChatService
}

func New(userProvider auth.UserProvider, logger *slog.Logger) *App {
	return &App{
		authService: auth.NewAuthService(userProvider, logger),
		chatService: chat.NewChatService(logger),
	}
}

func (a *App) Run(addr string) error {
	for i := 0; i < 10; i++ {
		roomID := "room" + string(rune('0'+i))
		room := domain.NewRoom(roomID)
		a.chatService.AddRoom(room)
		go room.Poll()
	}

	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	metrics.Register()

	registerHandler := handler.NewRegisterHandler(a.authService)
	loginHandler := handler.NewLoginHandler(a.authService)
	chatHandler := handler.NewChatHandler(a.chatService)

	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.Handle("/metrics", promhttp.Handler())

	r.Post("/register", registerHandler.Register())
	r.Post("/login", loginHandler.Login())

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/chat/{roomID}", chatHandler.Chat())
	})

	return http.ListenAndServe(addr, r)
}
