package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type client struct {
	count    int
	lastSeen time.Time
}

type rateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		clients: make(map[string]*client),
		limit:   limit,
		window:  window,
	}

	go rl.cleanup()
	return rl
}

func (rl *rateLimiter) Middleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Invalid IP", http.StatusBadRequest)
				return
			}

			rl.mu.Lock()

			c, ok := rl.clients[ip]
			if !ok || time.Since(c.lastSeen) > rl.window {
				rl.clients[ip] = &client{count: 1, lastSeen: time.Now()}
				rl.mu.Unlock()
				next.ServeHTTP(w, r)
				return
			}

			if c.count >= rl.limit {
				rl.mu.Unlock()
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			c.count++
			c.lastSeen = time.Now()
			rl.mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

func (rl *rateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, c := range rl.clients {
			if time.Since(c.lastSeen) > rl.window {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}
