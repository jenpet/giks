package plugins

import (
	"fmt"
	"os"
)

const (
	varFailOnMismatch    = "FAIL_ON_MISMATCH"
	varValidationPattern = "VALIDATION_PATTERN"
)

type StringValidator struct{}

func (sv StringValidator) ID() string {
	return "string-validator"
}

func (sv StringValidator) Run(workingDir string, hook string, vars map[string]string, args []string) (bool, error) {
	failOnMismatch, err := extractBoolVar(varFailOnMismatch, vars, false)
	if err != nil {
		return true, err
	}

	pat, err := extractStringVar(varValidationPattern, vars, true)
	if err != nil {
		return true, err
	}
	switch hook {
	// TODO: need re-work. In the current setup it is only usable for a commit-msg hook
	case "commit-msg":
		b, err := os.ReadFile(args[0])
		if err != nil {
			return false, fmt.Errorf("could not read file '%s'", err)
		}
		return failOnMismatch, matchString(string(b), pat)
	default:
		return hookUnsupported(hook, sv.ID())
	}
}
