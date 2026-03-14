package main

import (
	"log"
	"net/http"

	"github.com/brainless-security/brainless-waf/core/pkg/common"
	"github.com/brainless-security/brainless-waf/core/pkg/proxy"
	"github.com/brainless-security/brainless-waf/core/pkg/rules"
)

func main() {
	// Default values
	listenAddr := ":80"
	upstreamURL := "http://localhost:8080"

	// Try to load config from file
	cfg, err := common.LoadConfig("config/config.yaml")
	if err == nil {
		listenAddr = cfg.Server.Listen
		upstreamURL = cfg.Server.Upstream
		log.Printf("Loaded configuration from config/config.yaml")
	} else {
		log.Printf("Using default bootstrap configuration (no config/config.yaml found)")
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

	wafProxy, err := proxy.NewWAFProxy(upstreamURL, engine)
	if err != nil {
		log.Fatalf("Failed to initialize WAF proxy: %v", err)
	}

	http.HandleFunc("/", wafProxy.ServeHTTP)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Printf("Brainless WAF Core starting on %s...", listenAddr)
	log.Printf("Proxying traffic to %s", upstreamURL)

	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
