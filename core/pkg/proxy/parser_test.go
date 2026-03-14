package proxy

import (
	"net/http/httptest"
	"testing"

	"github.com/brainless-security/brainless-waf/core/pkg/common"
)

func TestRequestNormalization(t *testing.T) {
	parser := NewParser()

	// 1. Test URL Double Decoding
	doubleEncoded := "http://example.com/search?q=%2527%2520OR%25201%253D1"

	req := httptest.NewRequest("GET", doubleEncoded, nil)
	tx := common.NewTransaction(req)

	if err := parser.Parse(tx); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	expectedArg := "' or 1=1"
	if tx.Args.Get("q") != expectedArg {
		t.Errorf("Expected normalized arg %q, got %q", expectedArg, tx.Args.Get("q"))
	}

	// 2. Test Path Normalization
	reqPath := httptest.NewRequest("GET", "/ADMIN//login", nil)
	txPath := common.NewTransaction(reqPath)
	parser.Parse(txPath)

	expectedPath := "/admin/login"
	if txPath.NormalizedPath != expectedPath {
		t.Errorf("Expected normalized path %q, got %q", expectedPath, txPath.NormalizedPath)
	}

	// 3. Test Unicode Normalization
	unicodeInput := "caf\u0065\u0301"
	reqUni := httptest.NewRequest("GET", "/?name="+unicodeInput, nil)
	txUni := common.NewTransaction(reqUni)
	parser.Parse(txUni)

	expectedUni := "café"
	if txUni.Args.Get("name") != expectedUni {
		t.Errorf("Expected normalized unicode %q, got %q", expectedUni, txUni.Args.Get("name"))
	}
}
