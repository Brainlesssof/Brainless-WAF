package limiter

import (
	"testing"
	"time"
)

func TestIPVoiceLimiter(t *testing.T) {
	// 1 request per second, burst of 1
	l := NewIPVoiceLimiter(1.0, 1)
	ip := "192.168.1.1"

	// First request should be allowed
	if !l.Allow(ip) {
		t.Errorf("Expected first request to be allowed")
	}

	// Second request immediately after should be blocked
	if l.Allow(ip) {
		t.Errorf("Expected second request to be blocked")
	}

	// Different IP should be allowed
	if !l.Allow("10.0.0.1") {
		t.Errorf("Expected request from different IP to be allowed")
	}

	// Wait for 1 second
	time.Sleep(1100 * time.Millisecond)

	// Should be allowed again
	if !l.Allow(ip) {
		t.Errorf("Expected request to be allowed after 1 second")
	}
}
