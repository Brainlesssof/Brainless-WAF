package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/brainless-security/brainless-waf/core/pkg/common"
	"github.com/brainless-security/brainless-waf/core/pkg/rules"
)

type WAFProxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
	parser *Parser
	engine *rules.Engine
}

func NewWAFProxy(targetURL string, engine *rules.Engine) (*WAFProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	return &WAFProxy{
		target: target,
		proxy:  proxy,
		parser: NewParser(),
		engine: engine,
	}, nil
}

func (p *WAFProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		if result.Matched && result.Action == "deny" {
			http.Error(w, result.Message, result.Status)
			return
		}
	}

	p.proxy.ServeHTTP(w, r)
}
