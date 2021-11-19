package config

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestParseConfigFile_whenInputIsValid_shouldParseCorrectly(t *testing.T) {
	cfg := parseConfigFile("../test/files/giks-testconfig.yml")
	// test HookList() and HookListNames()
	assert.NotContains(t, cfg.HookListNames(false), "pre-push", "hook list should not return disabled hook")
	assert.Len(t, cfg.HookList(false), 2, "hook list should be filtered for active hooks")
	assert.Contains(t, cfg.HookListNames(true), "pre-push", "hook list should not return disabled hook")
	assert.Len(t, cfg.HookList(true), 3, "hook list should be filtered for active hooks")

	// test Hook() and LookupHook()
	assert.Equal(t, Hook{false, nil, "pre-push"}, cfg.Hook("pre-push"), "hook from the config should be returned")
	lookup, err := cfg.LookupHook("absent")
	assert.Nil(t, lookup, "no hook result expected when looking up an absent hook")
	assert.Error(t, err, "error expected when looking up an absent hook")

	// test ToMap()
	h := cfg.Hook("pre-commit").ToMap()
	assert.Len(t, h["steps"], 1, "hook map step amount should match config")
	assert.Equal(t, "pre-commit", h["name"].(string), "hook map's name should match config")
	assert.True(t, h["enabled"].(bool), "hook map's status should match config")
}

func TestParseConfig_whenInputIsInvalid_shouldReturnError(t *testing.T) {
	configTests := []struct {
		name  string
		input io.Reader
	}{
		{
			"malformed YAML",
			strings.NewReader(`!foobar`),
		},
		{
			"erroneous reader",
			new(errReader),
		},
		{
			"unsupported hook",
			strings.NewReader(`
hooks:
  foo:
    enabled: true`),
		},
	}

	for _, tt := range configTests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := parseConfig(tt.input)
			assert.Nil(t, cfg, "no config expected when providing an invalid configuration input")
			assert.Error(t, err, "error expected when providing an invalid configuration input")
		})
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("artificial error")
}
