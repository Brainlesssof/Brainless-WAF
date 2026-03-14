package main

import (
	"log"
	"net/http"

	"github.com/brainless-security/brainless-waf/core/pkg/common"
	"github.com/brainless-security/brainless-waf/core/pkg/limiter"
	"github.com/brainless-security/brainless-waf/core/pkg/proxy"
	"github.com/brainless-security/brainless-waf/core/pkg/rules"
)

func main() {
	// Default values
	listenAddr := ":80"
	upstreamURL := "http://localhost:8080"
	var cfg *common.Config

	// Try to load config from file
	c, err := common.LoadConfig("config/config.yaml")
	if err == nil {
		cfg = c
		listenAddr = cfg.Server.Listen
		upstreamURL = cfg.Server.Upstream
		log.Printf("Loaded configuration from config/config.yaml")
	} else {
		log.Printf("Using default bootstrap configuration (no config/config.yaml found)")
	}

	// Initialize Rate Limiter
	var l *limiter.IPVoiceLimiter
	if cfg != nil && cfg.RateLimiting.Enabled {
		l = limiter.NewIPVoiceLimiter(cfg.RateLimiting.RPS, cfg.RateLimiting.Burst)
		log.Printf("Rate limiting enabled: RPS=%.2f, Burst=%d", cfg.RateLimiting.RPS, cfg.RateLimiting.Burst)
	}

	// Initialize Rule Engine
	engine := rules.NewEngine()
	ruleParser := rules.NewParser()
	loadedRules, err := ruleParser.ParseFile("rules/default.rules")
	if err == nil {
		engine.Rules = loadedRules
		log.Printf("Successfully loaded %d rules from rules/default.rules", len(loadedRules))
	} else {
		log.Printf("No rules loaded: %v", err)
	}

	wafProxy, err := proxy.NewWAFProxy(upstreamURL, engine, l)
	if err != nil {
		log.Fatalf("Failed to initialize WAF proxy: %v", err)
	}

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	// All other traffic goes through the proxy
	http.Handle("/", wafProxy)

	// Start server(s)
	if cfg != nil && cfg.Server.TLS.Enabled {
		go func() {
			log.Printf("WAF Core Engine (HTTPS) listening on %s", cfg.Server.TLS.ListenTLS)
			if err := http.ListenAndServeTLS(cfg.Server.TLS.ListenTLS, cfg.Server.TLS.CertFile, cfg.Server.TLS.KeyFile, nil); err != nil {
				log.Fatalf("HTTPS server failed: %v", err)
			}
		}()
	}

	log.Printf("WAF Core Engine (HTTP) listening on %s", listenAddr)
	log.Printf("Proxying traffic to %s", upstreamURL)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
