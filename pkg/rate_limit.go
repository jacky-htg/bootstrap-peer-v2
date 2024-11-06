package pkg

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu          sync.Mutex
	requests    map[string]int
	lastRequest map[string]time.Time
	limit       int
	window      time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:    make(map[string]int),
		lastRequest: make(map[string]time.Time),
		limit:       limit,
		window:      window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if last, ok := rl.lastRequest[ip]; ok && now.Sub(last) < rl.window {
		if rl.requests[ip] >= rl.limit {
			return false
		}
		rl.requests[ip]++
	} else {
		rl.requests[ip] = 1
	}

	rl.lastRequest[ip] = now
	return true
}
