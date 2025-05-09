package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestRateLimiterMiddleware(t *testing.T) {

	if os.Getenv("CI") == "true" {
		t.Skip("Skipping rate limiter test in CI environment")
	}

	cfg := config{
		rateLimiter: rateLimiterConfig{
			requestsPerTimeFrame: 20,
			timeFrame:            time.Second * 5,
			enabled:              true,
		},
	}

	app := newTestApplication(t, cfg)
	ts := httptest.NewServer(app.mount())
	defer ts.Close()

	client := &http.Client{}
	mockIP := "192.168.1.1"
	marginOfError := 2

	for i := 0; i < cfg.rateLimiter.requestsPerTimeFrame+marginOfError; i++ {
		req, err := http.NewRequest("GET", ts.URL+"/v1/health", nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		req.Header.Set("X-Forwarded-For", mockIP)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not send request: %v", err)
		}
		io.Copy(io.Discard, resp.Body)
		defer resp.Body.Close()

		if i < cfg.rateLimiter.requestsPerTimeFrame {
			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", resp.Status)
			}
		} else {
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Errorf("expected status Too Many Requests; got %v", resp.Status)
			}
		}
	}
}
