package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Mixin interface {
	ID() string
	Enrich(dir string, vars map[string]string) (bool, error)
}

var Mixins = []Mixin{
	StagedFilesLister{},
}

// GetMixin returns a plugin for the given name (identifier). In case none was found an error is returned.
func GetMixin(name string) (Mixin, error) {
	for _, p := range Mixins {
		if p.ID() == name || p.ID() == strings.ToUpper(name) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("mixin '%s' not found", name)
}

type StagedFilesLister struct{}

func (sfl StagedFilesLister) ID() string {
	return "STAGED_FILES"
}

func (sfl StagedFilesLister) Enrich(dir string, vars map[string]string) (bool, error) {
	cmd := exec.Command("git", "-C", dir, "diff", "--cached", "--name-only")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return false, err
	}
	s := buf.String()
	files := strings.Replace(s, "\n", " ", -1)
	vars["GIKS_MIXIN_STAGED_FILES"] = files
	return true, nil
}
