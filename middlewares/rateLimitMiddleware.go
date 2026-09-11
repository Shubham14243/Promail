package middlewares

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	defaultRateLimitRequests = 100
	defaultRateLimitWindow   = time.Minute
)

type rateLimitEntry struct {
	count       int
	windowStart time.Time
}

type RateLimiter struct {
	limit  int
	window time.Duration

	mu      sync.Mutex
	clients map[string]rateLimitEntry
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	if limit <= 0 {
		limit = defaultRateLimitRequests
	}
	if window <= 0 {
		window = defaultRateLimitWindow
	}

	return &RateLimiter{
		limit:   limit,
		window:  window,
		clients: make(map[string]rateLimitEntry),
	}
}

func NewRateLimiterFromEnv() *RateLimiter {
	limit := envInt("APP_RATE_LIMIT_REQUESTS", defaultRateLimitRequests)
	windowSeconds := envInt("APP_RATE_LIMIT_WINDOW_SECONDS", int(defaultRateLimitWindow/time.Second))

	return NewRateLimiter(limit, time.Duration(windowSeconds)*time.Second)
}

func (l *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		client := clientAddress(r)

		l.mu.Lock()
		l.removeExpired(now)

		entry, exists := l.clients[client]
		if !exists || now.Sub(entry.windowStart) >= l.window {
			entry = rateLimitEntry{windowStart: now}
		}

		entry.count++
		remaining := l.limit - entry.count
		if remaining < 0 {
			remaining = 0
		}
		l.clients[client] = entry
		l.mu.Unlock()

		resetAt := entry.windowStart.Add(l.window)
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(l.limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetAt.Unix(), 10))

		if entry.count > l.limit {
			retryAfter := int(time.Until(resetAt).Seconds())
			if retryAfter < 1 {
				retryAfter = 1
			}
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RateLimit(next http.Handler) http.Handler {
	return NewRateLimiterFromEnv().Middleware(next)
}

func (l *RateLimiter) removeExpired(now time.Time) {
	for client, entry := range l.clients {
		if now.Sub(entry.windowStart) >= l.window {
			delete(l.clients, client)
		}
	}
}

func clientAddress(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func envInt(name string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return fallback
	}
	return value
}
