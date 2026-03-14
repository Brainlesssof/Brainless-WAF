# Anomaly Scoring & Stateful Logic Design (v0.4)

## Anomaly Scoring
Anomaly scoring moves away from simple "first-match-deny" logic to a more nuanced approach where requests are blocked only if their cumulative risk score exceeds a threshold.

### BRF Extensions
- `severity`: Metadata indicating the risk level (CRITICAL, HIGH, MEDIUM, LOW).
- `setvar`: Action to increase a variable (e.g., `tx.anomaly_score=+5`).

### Logic
1. Each rule match increases the `tx.anomaly_score`.
2. A final evaluation phase checks if `tx.anomaly_score >= threshold`.
3. If exceeded, a disruptive action (DENY) is taken.

## Stateful Logic (Phase 1)
Stateful logic allows tracking data across multiple requests from the same source.

### Persistent Collections (In-Memory for v0.4)
- `IP`: A collection keyed by source IP.
- `SESSION`: (Future) Keyed by session ID.

### Example Rule
```brf
# Increase anomaly score on match
SecRule ARGS "@rx union" "id:1001,msg:'SQLi',severity:CRITICAL,setvar:tx.anomaly_score=+5"

# Block if total score > 10
SecAction "id:9000,phase:5,deny,status:403,msg:'Inbound Anomaly Score Exceeded',chain"
SecRule TX:ANOMALY_SCORE "@ge 10" ""
```

## Implementation Plan
1. Update `Rule` struct to support `Severity` and `Variables` (maps).
2. Implement `setvar` handling in `Engine.Evaluate`.
3. Add a dedicated `Threshold` check logic in `WAFProxy`.
