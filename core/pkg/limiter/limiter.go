package limiter

import (
	"sync"

	"golang.org/x/time/rate"
)

// IPVoiceLimiter manages rate limiters for individual IP addresses.
type IPVoiceLimiter struct {
	ips   map[string]*rate.Limiter
	mu    sync.Mutex
	rps   rate.Limit
	burst int
}

func NewIPVoiceLimiter(rps float64, burst int) *IPVoiceLimiter {
	return &IPVoiceLimiter{
		ips:   make(map[string]*rate.Limiter),
		rps:   rate.Limit(rps),
		burst: burst,
	}
}

func (i *IPVoiceLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	l, exists := i.ips[ip]
	if !exists {
		l = rate.NewLimiter(i.rps, i.burst)
		i.ips[ip] = l
	}

	return l
}

func (i *IPVoiceLimiter) Allow(ip string) bool {
	return i.GetLimiter(ip).Allow()
}
