package plugins

import (
	"bytes"
	"fmt"
	"github.com/jenpet/giks/log"
	"os/exec"
	"strings"
)

const (
	varFilePattern = "FILE_WATCHER_PATTERN"
	varCommand     = "FILE_WATCHER_COMMAND"
	varFileList    = "FILE_WATCHER_FILES_LIST"
)

type FileWatcher struct{}

func (fw FileWatcher) ID() string {
	return "file-watcher"
}

func (fw FileWatcher) Run(workingDir string, hook string, vars map[string]string, args []string) (bool, error) {
	pattern, err := extractStringVar(varFilePattern, vars, true)
	if err != nil {
		return false, err
	}

	var cmd *exec.Cmd
	err = extractVar(varCommand, vars, func(val string) error {
		cmd = exec.Command("sh", []string{"-c", val}...)
		cmd.Dir = workingDir
		return nil
	}, true)

	filesString, err := extractStringVar(varFileList, vars, false)
	if err != nil {
		return false, err
	}
	files := strings.Split(filesString, " ")
	if singleFileMatchesPattern(files, pattern) {
		log.Infof("[%s]: changes detected. Running '%s'...", fw.ID(), strings.Join(cmd.Args, " "))
		var buf bytes.Buffer
		cmd.Stdout = &buf
		cmd.Stderr = &buf
		err = cmd.Run()
		if err != nil {
			return false, fmt.Errorf("files matched pattern but command failed: %+v: %s", err, buf.String())
		}
		return true, nil
	}
	return true, nil
}

func singleFileMatchesPattern(files []string, pattern string) bool {
	for _, file := range files {
		if err := matchString(file, pattern); err == nil {
			return true
		}
	}
	return false
}
