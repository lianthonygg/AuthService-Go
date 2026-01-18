package server

import (
	"database/sql"
	"net/http"
	"time"

	"auth-service/internal/features/user/service"
	"auth-service/internal/features/user/store"
	"auth-service/internal/features/user/transport"
	"auth-service/internal/middleware"
)

func New(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	rl := middleware.NewRateLimiter(5, time.Minute)

	// User Service
	userStore := store.New(db)
	userService := service.New(userStore)
	userHandler := transport.New(*userService)

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	mux.HandleFunc("/users", userHandler.UsersHandler)
	mux.HandleFunc("/users/", userHandler.UserHandlerById)

	handler := middleware.Chain(
		mux,
		middleware.Recover(),
		middleware.RequestId(),
		middleware.Logger(),
		rl.Middleware(),
	)

	return handler
}
