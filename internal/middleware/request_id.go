package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type requestIDKeyType struct{}

var requestIDKey = requestIDKeyType{}

func RequestId() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqId := r.Header.Get("X-Request-ID")
			if reqId == "" {
				reqId = uuid.NewString()
			}

			ctx := context.WithValue(r.Context(), requestIDKey, reqId)
			w.Header().Set("X-Request-ID", reqId)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}