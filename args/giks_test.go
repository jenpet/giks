package args

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGiksArgs_whenInputsVary_shouldResultInNoError(t *testing.T) {
	argTests := []struct {
		name string
		input []string
		expectedCmd string
		expectedSubCmd string
		expectedArgs []string
	}{
		{
			"full-blown input",
			toArgs("giks hooks exec commit-msg --all=true"),
			"hooks",
			"exec",
			toArgs("commit-msg --all=true"),
		},
		{
			"missing subcommand input",
			toArgs("giks hooks"),
			"hooks",
			"help",
			[]string{},
		},
		{
			"missing command input",
			toArgs("giks"),
			"help",
			"help",
			[]string{},
		},
	}
	for _, tt := range argTests {
		t.Run(tt.name, func(t *testing.T) {
			var ga GiksArgs = tt.input
			assert.Equal(t, tt.expectedCmd, ga.Command(), "expected command and resulting command do not match")
			assert.Equal(t, tt.expectedSubCmd, ga.SubCommand(), "expected sub-command and resulting sub-command do not match")
			assert.Equal(t, tt.expectedArgs, ga.Args(), "expected args and resulting args do not match")
		})
	}
}

func TestGiksArgs_whenInputContainsConfigPath_shouldExtractAccordingly(t *testing.T) {
	var ga GiksArgs = []string{"hooks", "exec", "commit-msg", "FEAT: Hallo", "--config=giks_alternative.yml"}
	assert.Equal(t, "giks_alternative.yml", ga.ConfigFile(), "expected config file and resulting config file does not match")
	assert.NotContains(t, ga.Args(), "--config=giks_alternative.yml", "giks args should not contain global config flags")
}

func toArgs(s string) []string {
	return strings.Split(s, " ")
}
