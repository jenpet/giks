package plugins

import (
	"errors"
	"fmt"
	"giks/log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type StringValidator struct {}

func (sv StringValidator) ID() string {
	return "string-validator"
}

func (sv StringValidator) Run(hook string, vars map[string]string, args []string) (bool, error) {
	failOnMismatch := false
	if val, ok := vars["FAIL_ON_MISMATCH"]; ok {
		var err error
		failOnMismatch, err = strconv.ParseBool(val)
		if err != nil {
			return false, errors.New("failed parsing 'FAIL_ON_MISMATCH' flag")
		}
	}
	switch hook {
	case "commit-msg":
		b, err := os.ReadFile(args[0])
		if err != nil {
			return false, fmt.Errorf("could not read file '%s'", err)
		}
		return failOnMismatch, validateString(string(b), vars["VALIDATION_PATTERN"])
	default:
		log.Warnf("hook '%s' not supported by string validator")
		return false, nil
	}
}

func validateString(input, pattern string) error {
	if strings.TrimSpace(pattern) == "" {
		return errors.New("required variable 'VALIDATION_PATTERN' missing or empty")
	}
	r, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("provided pattern '%s' can not be used as a regexp", pattern)
	}
	if !r.MatchString(input) {
		return fmt.Errorf("input '%s' does not match required pattern '%s'", input, pattern)
	}
	return nil
}