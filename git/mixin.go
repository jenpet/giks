package git

import (
	"bytes"
	"fmt"
	"github.com/jenpet/giks/log"
	"os/exec"
	"strconv"
	"strings"
)

var mixins = []mixinFunc{
	mixinModifiedFiles,
	mixinStagedFiles,
	mixinHeadFiles,
}

func ApplyMixins(dir string, vars map[string]string) {
	for n, m := range mixins {
		if err := m(dir, vars); err != nil {
			log.Warnf("Failed applying git mixin no. %d. Error: %+v", n, err)
		}
	}
}

type mixinFunc func(dir string, vars map[string]string) error

var mixinModifiedFiles = func(dir string, vars map[string]string) error {
	out, err := execGitCommand(dir, "ls-files", "-m")
	if err != nil {
		return err
	}
	vars["GIKS_MIXIN_MODIFIED_FILES"] = strings.TrimSpace(strings.Replace(out, "\n", " ", -1))
	return nil
}

var mixinStagedFiles = func(dir string, vars map[string]string) error {
	out, err := execGitCommand(dir, "diff", "--cached", "--name-only")
	if err != nil {
		return err
	}
	vars["GIKS_MIXIN_STAGED_FILES"] = strings.Replace(out, "\n", " ", -1)
	return nil
}

var mixinHeadFiles = func(dir string, vars map[string]string) error {
	out, err := execGitCommand(dir, "rev-list", "--count", "HEAD")
	if err != nil {
		return err
	}
	commitAmount, err := strconv.Atoi(out)
	if err != nil {
		return err
	}
	diffCommit := "HEAD~"
	if commitAmount == 1 {
		// use the empty tree when there is only one commit
		diffCommit = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"
	}
	out, err = execGitCommand(dir, "diff", "--name-only", "HEAD", diffCommit)
	if err != nil {
		return err
	}

	vars["GIKS_MIXIN_HEAD_FILES"] = strings.Replace(out, "\n", " ", -1)
	return nil
}

func execGitCommand(dir string, arg ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, arg...)...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed executing git command '%s'. Error: %s", strings.Join(arg, " "), buf.String())
	}
	return strings.TrimSpace(buf.String()), nil
}
