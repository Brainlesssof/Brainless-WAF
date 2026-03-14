# Brainless Rule Format (BRF) Specification

BRF is a lightweight, line-based format inspired by ModSecurity but optimized for high-performance processing in Go.

## Example Rule

```brf
# Block SQL Injection in ARGS
SecRule ARGS "@rx (?i)(union|select|insert|update|delete|drop)" "id:1001,phase:2,deny,status:403,msg:'SQL Injection Detected'"

# Block XSS in Request body
SecRule REQUEST_BODY "@contains <script>" "id:1002,phase:2,deny,status:403,msg:'XSS Attempt'"
```

## Structure
- **SecRule**: Keyword to define a rule.
- **Variable**: The part of the request to inspect (e.g., `ARGS`, `REQUEST_URI`, `HEADERS`).
- **Operator**: The matching logic, prefixed with `@` (e.g., `@rx`, `@contains`, `@streq`).
- **Actions**: comma-separated metadata and actions (e.g., `id`, `msg`, `deny`, `allow`).

## Supported Components (v0.3)
### Variables
- `ARGS`: All query string parameters.
- `REQUEST_URI`: The full normalized URI.
- `REQUEST_HEADERS`: All request headers.
- `REQUEST_BODY`: The raw request body (normalized).

### Operators
- `@rx`: Regular expression match.
- `@contains`: Substring match.
- `@streq`: Exact string equality.

### Actions
- `id`: unique rule identifier.
- `msg`: descriptive message for logging.
- `deny`: immediately stop processing and block request.
- `pass`: continue to next rule if matched.
- `status`: HTTP status code to return on deny.
