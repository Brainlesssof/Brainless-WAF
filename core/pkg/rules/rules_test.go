package rules

import (
	"net/http/httptest"
	"testing"

	"github.com/brainless-security/brainless-waf/core/pkg/common"
)

func TestRuleEvaluation(t *testing.T) {
	engine := NewEngine()

	// Add test rules
	engine.Rules = []Rule{
		{
			ID:       1001,
			Variable: "ARGS",
			Operator: "contains",
			Operand:  "select",
			Actions:  []string{"deny"},
			Message:  "SQL Injection Detected",
			Status:   403,
		},
		{
			ID:       1002,
			Variable: "REQUEST_PATH",
			Operator: "streq",
			Operand:  "/admin",
			Actions:  []string{"log"},
			Message:  "Admin Access Attempt",
		},
	}

	// 1. Test Match (Deny)
	req := httptest.NewRequest("GET", "/?query=select+%2A+from+users", nil)
	tx := common.NewTransaction(req)
	tx.Args = req.URL.Query()

	result := engine.Evaluate(tx)
	if !result.Matched {
		t.Errorf("Expected rule 1001 to match")
	}
	if result.Action != "deny" {
		t.Errorf("Expected action 'deny', got %q", result.Action)
	}

	// 2. Test No Match
	req2 := httptest.NewRequest("GET", "/safe-path", nil)
	tx2 := common.NewTransaction(req2)
	result2 := engine.Evaluate(tx2)
	if result2.Matched {
		t.Errorf("Expected no rules to match")
	}
}
