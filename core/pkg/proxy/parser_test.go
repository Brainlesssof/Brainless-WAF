package proxy

import (
	"net/http/httptest"
	"testing"
)

func TestRequestNormalization(t *testing.T) {
	parser := NewParser()

	// 1. Test URL Double Decoding
	doubleEncoded := "http://example.com/search?q=%2527%2520OR%25201%253D1"
	// %2527 -> %27 -> '
	// %2520 -> %20 -> (space)

	req := httptest.NewRequest("GET", doubleEncoded, nil)
	tx := NewTransaction(req)

	if err := parser.Parse(tx); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	expectedArg := "' or 1=1" // Lowercased because of NormalizeString
	if tx.Args.Get("q") != expectedArg {
		t.Errorf("Expected normalized arg %q, got %q", expectedArg, tx.Args.Get("q"))
	}

	// 2. Test Path Normalization
	reqPath := httptest.NewRequest("GET", "/ADMIN//login", nil)
	txPath := NewTransaction(reqPath)
	parser.Parse(txPath)

	expectedPath := "/admin/login"
	if txPath.NormalizedPath != expectedPath {
		t.Errorf("Expected normalized path %q, got %q", expectedPath, txPath.NormalizedPath)
	}

	// 3. Test Unicode Normalization (NFC)
	// 'e' (U+0065) + acute accent (U+0301) should become 'é' (U+00E9) in NFC
	unicodeInput := "caf\u0065\u0301" // café (decomposed)
	reqUni := httptest.NewRequest("GET", "/?name="+unicodeInput, nil)
	txUni := NewTransaction(reqUni)
	parser.Parse(txUni)

	expectedUni := "café"
	if txUni.Args.Get("name") != expectedUni {
		// Note: The comparison might fail if the source file itself isn't NFC,
		// but Go strings are typically UTF-8 and the literal "café" is NFC.
		t.Errorf("Expected normalized unicode %q, got %q", expectedUni, txUni.Args.Get("name"))
	}
}
