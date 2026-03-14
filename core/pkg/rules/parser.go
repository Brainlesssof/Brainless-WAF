package rules

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Parser handles reading .rules files and converting them to Rule objects.
type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseFile(path string) ([]Rule, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rules []Rule
	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Regex to parse SecRule line: SecRule VARIABLE "OPERATOR OPERAND" "ACTIONS"
	// Example: SecRule ARGS "@rx union" "id:100,deny,status:403"
	ruleRegex := regexp.MustCompile(`^SecRule\s+([A-Z_]+)\s+"@([a-z]+)\s+(.*)"\s+"(.*)"`)

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		matches := ruleRegex.FindStringSubmatch(line)
		if len(matches) != 5 {
			return nil, fmt.Errorf("invalid rule format at line %d: %s", lineNum, line)
		}

		variable := matches[1]
		operator := matches[2]
		operand := matches[3]
		actionStr := matches[4]

		rule := Rule{
			Variable: variable,
			Operator: operator,
			Operand:  operand,
			Status:   403, // Default
			SetVars:  make(map[string]string),
		}

		// Compile regex if operator is rx
		if operator == "rx" {
			rx, err := regexp.Compile(operand)
			if err != nil {
				return nil, fmt.Errorf("invalid regex at line %d: %v", lineNum, err)
			}
			rule.Rx = rx
		}

		// Parse actions
		actions := strings.Split(actionStr, ",")
		for _, a := range actions {
			a = strings.TrimSpace(a)
			if strings.HasPrefix(a, "id:") {
				id, _ := strconv.Atoi(strings.TrimPrefix(a, "id:"))
				rule.ID = id
			} else if strings.HasPrefix(a, "msg:") {
				rule.Message = strings.Trim(strings.TrimPrefix(a, "msg:"), "'\"")
			} else if strings.HasPrefix(a, "status:") {
				status, _ := strconv.Atoi(strings.TrimPrefix(a, "status:"))
				rule.Status = status
			} else if strings.HasPrefix(a, "severity:") {
				rule.Severity = strings.ToUpper(strings.TrimPrefix(a, "severity:"))
			} else if strings.HasPrefix(a, "setvar:") {
				varParts := strings.Split(strings.TrimPrefix(a, "setvar:"), "=")
				if len(varParts) == 2 {
					rule.SetVars[varParts[0]] = varParts[1]
				}
			} else {
				rule.Actions = append(rule.Actions, a)
			}
		}

		rules = append(rules, rule)
	}

	return rules, scanner.Err()
}
