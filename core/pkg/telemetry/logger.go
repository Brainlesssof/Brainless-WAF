package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/brainless-security/brainless-waf/core/pkg/common"
)

// SecurityEvent represents a structured log entry for a security match.
type SecurityEvent struct {
	Timestamp    string `json:"timestamp"`
	RequestID    string `json:"request_id"`
	ClientIP     string `json:"client_ip"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	MatchedRules []int  `json:"matched_rules"`
	AnomalyScore int    `json:"anomaly_score"`
	Action       string `json:"action"`
	Message      string `json:"message"`
}

// Logger handles structured logging of security events.
type Logger struct {
	output *os.File
}

func NewLogger() *Logger {
	return &Logger{
		output: os.Stderr, // Default to stderr for container logs
	}
}

func (l *Logger) LogMatch(tx *common.Transaction, msg string) {
	event := SecurityEvent{
		Timestamp:    time.Now().Format(time.RFC3339),
		RequestID:    tx.ID,
		ClientIP:     tx.Request.RemoteAddr,
		Method:       tx.Request.Method,
		Path:         tx.NormalizedPath,
		MatchedRules: tx.MatchedRules,
		AnomalyScore: tx.AnomalyScore,
		Action:       tx.Action,
		Message:      msg,
	}

	data, err := json.Marshal(event)
	if err == nil {
		fmt.Fprintln(l.output, string(data))
	}
}
