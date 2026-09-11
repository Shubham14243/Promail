package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterRejectsRequestsOverLimit(t *testing.T) {
	limiter := NewRateLimiter(2, time.Minute)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler := limiter.Middleware(next)

	for requestNumber := 0; requestNumber < 2; requestNumber++ {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.RemoteAddr = "192.0.2.1:1234"
		handler.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusNoContent {
			t.Fatalf("request %d: expected status 204, got %d", requestNumber+1, recorder.Code)
		}
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.RemoteAddr = "192.0.2.1:1234"
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status 429, got %d", recorder.Code)
	}
	if recorder.Header().Get("Retry-After") == "" {
		t.Fatal("expected Retry-After header")
	}
}

func TestRateLimiterSeparatesClients(t *testing.T) {
	limiter := NewRateLimiter(1, time.Minute)
	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	first := httptest.NewRequest(http.MethodGet, "/", nil)
	first.RemoteAddr = "192.0.2.1:1234"
	firstRecorder := httptest.NewRecorder()
	handler.ServeHTTP(firstRecorder, first)

	second := httptest.NewRequest(http.MethodGet, "/", nil)
	second.RemoteAddr = "192.0.2.2:1234"
	secondRecorder := httptest.NewRecorder()
	handler.ServeHTTP(secondRecorder, second)

	if firstRecorder.Code != http.StatusNoContent || secondRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected both clients to be allowed, got %d and %d", firstRecorder.Code, secondRecorder.Code)
	}
}
