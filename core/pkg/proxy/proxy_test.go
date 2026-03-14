package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
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
	backendResponse := "Hello from upstream"
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()

	p, err := NewWAFProxy(backend.URL, nil, nil, nil, 10)
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	req := httptest.NewRequest("GET", "/some-path", nil)
	w := httptest.NewRecorder()

	p.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if string(body) != backendResponse {
		t.Errorf("Expected body %q, got %q", backendResponse, string(body))
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}
