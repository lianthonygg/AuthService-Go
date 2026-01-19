package server

import (
	"database/sql"
	"net/http"
	"time"

	"auth-service/internal/config"
	authService "auth-service/internal/features/auth/service"
	authStore "auth-service/internal/features/auth/store"
	authHandler "auth-service/internal/features/auth/transport"
	"auth-service/internal/features/user/service"
	"auth-service/internal/features/user/store"
	"auth-service/internal/features/user/transport"
	"auth-service/internal/middleware"
	"auth-service/internal/shared/security"
)

func New(db *sql.DB, config *config.Config) http.Handler {
	mux := http.NewServeMux()
	rl := middleware.NewRateLimiter(5, time.Minute)

	// Security Services
	hasher := security.NewHasher()
	generator := security.NewGenerator(config)

	// User Service
	userStore := store.New(db, hasher)
	userService := service.New(userStore)
	userHandler := transport.New(*userService)

	// Auth Service
	autheStore := authStore.NewAuthStore(db)
	autheService := authService.New(*userService, autheStore, generator, hasher)
	autheHandler := authHandler.New(*autheService)

	// Auth Endpoints
	mux.Handle("GET /auth/me", middleware.JWTAuth(config.SecretKey)(http.HandlerFunc(autheHandler.AuthMeHandler)))
	mux.HandleFunc("POST /auth/login", autheHandler.LoginHandler)
	mux.HandleFunc("POST /auth/register", autheHandler.RegisterHandler)
	mux.HandleFunc("POST /auth/refresh", autheHandler.RefreshTokenHandler)
	mux.HandleFunc("POST /auth/logout", autheHandler.RefreshTokenRevokedHandler)

	// Users Endpoints
	mux.HandleFunc("GET /users", userHandler.GetAllUsersHandler)
	mux.HandleFunc("GET /users/", userHandler.GetByIdUserHandler)
	mux.HandleFunc("POST /users", userHandler.CreateUserHandler)
	mux.HandleFunc("PUT /users/", userHandler.UpdateUserHandler)
	mux.HandleFunc("DELETE /users/", userHandler.RemoveUserHandler)

	handler := middleware.Chain(
		mux,
		middleware.Recover(),
		middleware.RequestId(),
		middleware.Logger(),
		rl.Middleware(),
	)

	return handler
}
