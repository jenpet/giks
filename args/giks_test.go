package args

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGiksArgs_whenInputsVary_shouldResultInNoError(t *testing.T) {
	argTests := []struct {
		name           string
		input          []string
		expectedBinary string
		expectedCmd    string
		expectedSubCmd string
		expectedHook   string
		expectedArgs   []string
	}{
		{
			"full-blown input",
			toArgs("giks hooks exec commit-msg --all=true"),
			"giks",
			"hooks",
			"exec",
			"commit-msg",
			toArgs("--all=true"),
		},
		{
			"missing subcommand input",
			toArgs("giks hooks"),
			"giks",
			"hooks",
			"",
			"",
			[]string{},
		},
		{
			"missing command input",
			toArgs("giks"),
			"giks",
			"",
			"",
			"",
			[]string{},
		},
		{
			"subcommand is hook",
			toArgs("giks exec commit-msg --all=true"),
			"giks",
			"exec",
			"",
			"commit-msg",
			toArgs("--all=true"),
		},
	}
	for _, tt := range argTests {
		t.Run(tt.name, func(t *testing.T) {
			var ga GiksArgs = tt.input
			assert.Equal(t, tt.expectedBinary, ga.Binary(), "expected binary and resulting binary do not match")
			assert.Equal(t, tt.expectedCmd, ga.Command(), "expected command and resulting command do not match")
			assert.Equal(t, tt.expectedSubCmd, ga.SubCommand(), "expected sub-command and resulting sub-command do not match")
			assert.Equal(t, tt.expectedHook, ga.Hook(), "expected hook and resulting hook do not match")
			if tt.expectedHook != "" {
				assert.True(t, ga.HasHook(), "expected hook to be present")
			}
			assert.Equal(t, tt.expectedArgs, ga.Args(), "expected args and resulting args do not match")
		})
	}
}

func TestGiksArgs_whenInputHasGlobalFlags_shouldSanitizeAccordingly(t *testing.T) {
	input := []string{"hooks", "exec", "--config=giks_alternative.yml", "--git-dir=/foo/bar/.git/", "commit-msg", "FEAT: Hallo"}
	var ga GiksArgs = input
	assert.Equal(t, "giks_alternative.yml", ga.ConfigFile(), "expected config file and resulting config file does not match")
	assert.Equal(t, "/foo/bar/.git/", ga.GitDir(), "expected git dir and resulting git dir does not match")
	assert.NotContains(t, ga.Args(), "--config=giks_alternative.yml", "giks args should not contain global config flags")
	assert.NotContains(t, ga.Args(), "--git-dir=/foo/bar/.git/", "giks args should not contain global config flags")
	assert.Equal(t, input, ga.Raw(), "giks args should still contain raw arguments")
}

func toArgs(s string) []string {
	return strings.Split(s, " ")
}
