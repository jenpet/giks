package config

import (
	"giks/args"
	"giks/test/gittest"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"testing"
)

func TestAssembleConfig_shouldParseConfigAndValidatePaths(t *testing.T) {
	r := gittest.NewTestRepository("../test/output/git-dir")
	defer r.Clean()
	// using the 'true' command as a giks command test replacement
	var ga args.GiksArgs = []string{"true", "hooks", "exec", "commit-msg", "--config=../test/files/giks-testconfig.yml", "--git-dir=" + r.AbsDir()}
	absCfg, _ := filepath.Abs("../test/files/giks-testconfig.yml")
	bin, _ := exec.LookPath("true")

	cfg := AssembleConfig(ga)
	assert.Equal(t, absCfg, cfg.ConfigFile, "config file has to be available via an absolute path")
	assert.Equal(t, bin, cfg.Binary, "binary has to be available via an absolute path")
	assert.Equal(t, r.AbsGitDir(), cfg.GitDir, "git directory has to be available via an absolute path")
}

func TestAbsoluteFilePath_shouldResolveCorrectly(t *testing.T) {
	pathTests := []struct {
		name        string
		input       string
		expectedAbs string
	}{
		{
			"home dir",
			"~/foo/bar",
			func() string { u, _ := user.Current(); return u.HomeDir + "/foo/bar" }(),
		},
		{
			"relative",
			"./test",
			func() string { cwd, _ := os.Getwd(); return cwd + "/test" }(),
		},
		{
			"absolute",
			"/foo/bar",
			"/foo/bar",
		},
	}
	for _, tt := range pathTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedAbs, absoluteFilepath(tt.input), "input does not match the expected absolute output")
		})
	}
}
