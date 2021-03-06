package plugins

import (
	"errors"
	"fmt"
	"github.com/jenpet/giks/log"
	"regexp"
	"strconv"
	"strings"
)

// List holds all built-in giks plugins
var List = []Plugin{
	StringValidator{},
	FileWatcher{},
	ListComparator{},
}

// Get returns a plugin for the given name (identifier). In case none was found an error is returned.
func Get(name string) (Plugin, error) {
	for _, p := range List {
		if p.ID() == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("plugin '%s' not found", name)
}

// Plugin has to be implemented by all built-in plugins in order to be triggered correctly by the giks hook executor
type Plugin interface {
	// Run executes the plugin for a given workingDir and hook. The vars will contain all variables
	// provided by giks, args are the arguments that were passed for the hook execution.
	//
	// The returned boolean indicates whether giks should exit after running the plugin providing relying on additional
	// information given by the error.
	Run(workingDir string, hook string, vars map[string]string, args []string) (bool, error)

	// ID returns the name / identifier of the plugin which can be used within the configuration file
	ID() string
}

func extractVar(key string, vars map[string]string, parseFunc func(val string) error, required bool) error {
	if val, ok := vars[key]; ok {
		if strings.TrimSpace(val) == "" && required {
			return fmt.Errorf("variable '%s' is required but empty", key)
		}
		if err := parseFunc(val); err != nil {
			if err != nil {
				return fmt.Errorf("failed parsing '%s' variable", key)
			}
		}
	} else if required {
		return fmt.Errorf("variable '%s' is required but not set", key)
	}
	return nil
}

func extractStringVar(key string, vars map[string]string, required bool) (string, error) {
	var str string
	err := extractVar(key, vars, func(val string) error {
		str = val
		return nil
	}, required)
	return str, err
}

func extractBoolVar(key string, vars map[string]string, required bool) (bool, error) {
	var b bool
	err := extractVar(key, vars, func(val string) error {
		parsed, err := strconv.ParseBool(val)
		b = parsed
		return err
	}, required)
	return b, err
}

func matchString(input, pattern string) error {
	if strings.TrimSpace(pattern) == "" {
		return errors.New("required pattern missing or empty")
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

func hookUnsupported(hook string, plugin string) (bool, error) {
	log.Warnf("hook '%s' not supported by plugin '%s'", hook, plugin)
	return false, nil
}
