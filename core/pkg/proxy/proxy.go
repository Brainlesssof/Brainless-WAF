package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type WAFProxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
	parser *Parser
}

func NewWAFProxy(targetURL string) (*WAFProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	return &WAFProxy{
		target: target,
		proxy:  proxy,
		parser: NewParser(),
	}, nil
}

func (p *WAFProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Initialize security transaction
	tx := NewTransaction(r)

	// 2. Parse and normalize request
	if err := p.parser.Parse(tx); err != nil {
		http.Error(w, "Failed to parse request", http.StatusBadRequest)
		return
	}

	// TODO: Phase 1.3 -> Evaluate rules using tx

	p.proxy.ServeHTTP(w, r)
}
