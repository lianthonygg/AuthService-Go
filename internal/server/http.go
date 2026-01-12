package server

import (
	"auth-service/internal/middleware"
	"net/http"
	"time"
)

func New() http.Handler {
	mux := http.NewServeMux()
	rl := middleware.NewRateLimiter(100, time.Minute)

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	handler := middleware.Chain(
		mux,
		middleware.Recover(),
		middleware.RequestId(),
		middleware.Logger(),
		rl.Middleware(),
	)

	return handler
}
