package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/puremike/social-go/internal/auth"
	"github.com/puremike/social-go/internal/ratelimiter"
	"github.com/puremike/social-go/internal/store"
	"github.com/puremike/social-go/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	//  rate limiter
	rateLimiter := ratelimiter.NewFixedWindowRateLimiter(cfg.rateLimiter.requestsPerTimeFrame, cfg.rateLimiter.timeFrame)

	return &application{
		logger:        zap.NewNop().Sugar(),
		store:         store.NewMockStore(),
		cacheStorage:  cache.NewMockCacheStore(),
		authenticator: &auth.TestAuthenticator{},
		rateLimiter:   rateLimiter,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponse(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("want %d; got %d", expected, actual)
	}
}
