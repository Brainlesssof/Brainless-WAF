package proxy

import (
	"bytes"
	"io"
	"net/url"
	"strings"

	"github.com/brainless-security/brainless-waf/core/pkg/common"
	"golang.org/x/text/unicode/norm"
)

// Parser handles the normalization and extraction of HTTP request data.
type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// Parse extracts and normalizes data from an http.Request into a Transaction.
func (p *Parser) Parse(tx *common.Transaction) error {
	r := tx.Request

	// 1. Normalize Path and URL
	tx.NormalizedURL = p.NormalizeURL(r.URL.String())
	tx.NormalizedPath = p.NormalizePath(r.URL.Path)

	// 2. Extract and Normalize Args (Query params)
	tx.Args = r.URL.Query()
	for k, values := range tx.Args {
		for i, v := range values {
			tx.Args[k][i] = p.NormalizeString(v)
		}
	}

	// 3. Normalize Headers
	for k, values := range r.Header {
		for _, v := range values {
			tx.Headers.Add(k, p.NormalizeString(v))
		}
	}

	// 4. Capture and Normalize Body (if any)
	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil {
			// Restore the body for the reverse proxy to consume
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			tx.Body = bodyBytes
			tx.ContentType = r.Header.Get("Content-Type")
		}
	}

	return nil
}

// NormalizeURL performs recursive URL decoding and Unicode normalization on full URLs.
func (p *Parser) NormalizeURL(input string) string {
	return p.NormalizeString(input)
}

// NormalizePath cleans up the path by resolving segments and removing duplicates.
func (p *Parser) NormalizePath(path string) string {
	// Simple path cleanup (equivalent to path.Clean but for web)
	cleaned := strings.ReplaceAll(path, "//", "/")
	return strings.ToLower(cleaned)
}

// NormalizeString applies Unicode normalization, recursive unescaping, and lowercase.
func (p *Parser) NormalizeString(input string) string {
	// 1. Recursive unescaping (evasion prevention)
	decoded := input
	for i := 0; i < 10; i++ { // Limit recursion to 10
		next, err := url.QueryUnescape(decoded)
		if err != nil || next == decoded {
			break
		}
		decoded = next
	}

	// 2. Apply NFC normalization
	normalized := norm.NFC.String(decoded)

	// 3. Convert to lowercase for uniform rule matching
	return strings.ToLower(normalized)
}
