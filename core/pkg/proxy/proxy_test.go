package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	// Simple test for the health check concept
	// In the real app, this is handled in main()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
}

func TestProxyForwarding(t *testing.T) {
	// 1. Create a mock upstream server
	backendResponse := "Hello from upstream"
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()

	// 2. Create the WAF Proxy pointing to the backend
	p, err := NewWAFProxy(backend.URL)
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	// 3. Create a request to the proxy
	req := httptest.NewRequest("GET", "/some-path", nil)
	w := httptest.NewRecorder()

	// 4. Serve the request
	p.ServeHTTP(w, req)

	// 5. Verify the response
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if string(body) != backendResponse {
		t.Errorf("Expected body %q, got %q", backendResponse, string(body))
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}
