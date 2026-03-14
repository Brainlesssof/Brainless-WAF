package rules

import (
	"regexp"
	"strings"

	"github.com/brainless-security/brainless-waf/core/pkg/common"
)

// Rule represents a single security rule in BRF.
type Rule struct {
	ID       int
	Variable string // ARGS, REQUEST_URI, HEADERS, REQUEST_BODY
	Operator string // rx, contains, streq
	Operand  string
	Actions  []string // deny, allow, log, pass
	Message  string
	Status   int // HTTP status code

	// Pre-compiled components
	Rx *regexp.Regexp
}

// Result represents the outcome of a rule evaluation.
type Result struct {
	Matched bool
	Action  string
	Message string
	Status  int
}

// Engine maintains the set of rules and performs evaluations.
type Engine struct {
	Rules []Rule
}

func NewEngine() *Engine {
	return &Engine{
		Rules: []Rule{},
	}
}

func (e *Engine) Evaluate(tx *common.Transaction) Result {
	for _, rule := range e.Rules {
		// 1. Extract variable value
		value := e.extractVariable(tx, rule.Variable)

		// 2. Apply operator
		matched := false
		switch rule.Operator {
		case "rx":
			if rule.Rx != nil {
				matched = rule.Rx.MatchString(value)
			}
		case "contains":
			matched = strings.Contains(value, rule.Operand)
		case "streq":
			matched = (value == rule.Operand)
		}

		if matched {
			tx.MatchedRules = append(tx.MatchedRules, rule.ID)

			// Handle actions
			for _, action := range rule.Actions {
				if action == "deny" {
					tx.Action = "deny"
					return Result{
						Matched: true,
						Action:  "deny",
						Message: rule.Message,
						Status:  rule.Status,
					}
				}
			}
		}
	}

	return Result{Matched: false, Action: "allow"}
}

func (e *Engine) extractVariable(tx *common.Transaction, variable string) string {
	switch variable {
	case "REQUEST_URI":
		return tx.NormalizedURL
	case "REQUEST_PATH":
		return tx.NormalizedPath
	case "ARGS":
		var builder strings.Builder
		for k, values := range tx.Args {
			builder.WriteString(k)
			builder.WriteString("=")
			builder.WriteString(strings.Join(values, " "))
			builder.WriteString(" ")
		}
		return builder.String()
	case "REQUEST_HEADERS":
		var builder strings.Builder
		for k, values := range tx.Headers {
			builder.WriteString(k)
			builder.WriteString(": ")
			builder.WriteString(strings.Join(values, " "))
			builder.WriteString("\n")
		}
		return builder.String()
	case "REQUEST_BODY":
		return string(tx.Body)
	}
	return ""
}
