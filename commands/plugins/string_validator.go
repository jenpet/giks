package plugins

import (
	"fmt"
	"os"
	"strconv"
)

const (
	varFailOnMismatch    = "FAIL_ON_MISMATCH"
	varValidationPattern = "VALIDATION_PATTERN"
)

type StringValidator struct{}

func (sv StringValidator) ID() string {
	return "string-validator"
}

func (sv StringValidator) Run(hook string, vars map[string]string, args []string) (bool, error) {
	failOnMismatch := false
	err := extractVar(varFailOnMismatch, vars, func(val string) error {
		var err error
		failOnMismatch, err = strconv.ParseBool(val)
		return err
	}, false)
	if err != nil {
		return false, err
	}
	switch hook {
	case "commit-msg":
		b, err := os.ReadFile(args[0])
		if err != nil {
			return false, fmt.Errorf("could not read file '%s'", err)
		}
		return failOnMismatch, matchString(string(b), vars[varValidationPattern])
	default:
		return hookUnsupported(hook)
	}
}
