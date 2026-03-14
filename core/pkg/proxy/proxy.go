package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/brainless-security/brainless-waf/core/pkg/common"
	"github.com/brainless-security/brainless-waf/core/pkg/limiter"
	"github.com/brainless-security/brainless-waf/core/pkg/rules"
)

type WAFProxy struct {
	target  *url.URL
	proxy   *httputil.ReverseProxy
	parser  *Parser
	engine  *rules.Engine
	limiter *limiter.IPVoiceLimiter
}

func NewWAFProxy(targetURL string, engine *rules.Engine, l *limiter.IPVoiceLimiter) (*WAFProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	return &WAFProxy{
		target:  target,
		proxy:   proxy,
		parser:  NewParser(),
		engine:  engine,
		limiter: l,
	}, nil
}

func (p *WAFProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 0. Rate Limiting
	if p.limiter != nil {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		if !p.limiter.Allow(ip) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
	}

	// 1. Initialize security transaction
	tx := common.NewTransaction(r)

	// 2. Parse and normalize request
	if err := p.parser.Parse(tx); err != nil {
		http.Error(w, "Failed to parse request", http.StatusBadRequest)
		return
	}

	// 3. Evaluate rules
	if p.engine != nil {
		result := p.engine.Evaluate(tx)

		// Handle immediate deny actions
		if result.Matched && result.Action == "deny" {
			http.Error(w, result.Message, result.Status)
			return
		}

		// Handle Anomaly Scoring threshold
		// Default threshold = 10 (should be configurable)
		if tx.AnomalyScore >= 10 {
			http.Error(w, "Inbound Anomaly Score Exceeded", http.StatusForbidden)
			return
		}
	}

	p.proxy.ServeHTTP(w, r)
}
