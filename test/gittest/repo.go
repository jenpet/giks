package gittest

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func NewTestRepository(dir string) TestRepository {
	tr := TestRepository{dir: fmt.Sprintf("%s-%d", dir, time.Now().UnixNano())}
	tr.init()
	return tr
}

type TestRepository struct {
	dir string
}

func (tr TestRepository) init() {
	cmd := exec.Command("git", "init", tr.dir)
	if err := cmd.Run(); err != nil {
		panic("could not initiate test git repo. Error: " + err.Error())
	}
}

func (tr TestRepository) Command(arg ...string) (string, error) {
	args := append([]string{"-C", tr.dir}, arg...)
	var cmd = exec.Command("git", args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

func (tr TestRepository) WriteFile(filename, content string) {
	_ = os.WriteFile(filepath.Join(tr.dir, filename), []byte(content), 0777)
}

func (tr TestRepository) AddAll() {
	_, _ = tr.Command("add", "-A")
}

func (tr TestRepository) Commit(msg string) {
	_, _ = tr.Command("commit", "-m", msg)
}

func (tr TestRepository) AbsDir() string {
	abs, _ := filepath.Abs(tr.dir)
	return abs
}
func (tr TestRepository) AbsGitDir() string {
	abs, _ := filepath.Abs(tr.dir)
	return filepath.Join(abs, ".git")
}

func (tr TestRepository) Clean() {
	_ = os.RemoveAll(tr.dir)
}
