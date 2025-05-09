package ratelimiter

import (
	"sync"
	"time"
)

type Limiter interface {
	Allow(ip string) (bool, time.Duration)
}

type FixedWindowRateLimiter struct {
	sync.RWMutex
	client               map[string]int
	requestsPerTimeFrame int
	timeFrame            time.Duration
}

func NewFixedWindowRateLimiter(requestsPerTimeFrame int, timeFrame time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		client:               make(map[string]int),
		requestsPerTimeFrame: requestsPerTimeFrame,
		timeFrame:            timeFrame,
	}
}

func (rl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	rl.RLock()
	count, exists := rl.client[ip]
	rl.RUnlock()

	if !exists || count < rl.requestsPerTimeFrame {
		rl.Lock()
		if !exists {
			go rl.resetCount(ip)
		}

		rl.client[ip]++
		rl.Unlock()
		return true, 0
	}

	return false, rl.timeFrame
}

func (rl *FixedWindowRateLimiter) resetCount(ip string) {
	time.Sleep(rl.timeFrame)
	rl.Lock()
	delete(rl.client, ip)
	rl.Unlock()
}
