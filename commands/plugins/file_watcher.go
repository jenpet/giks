package plugins

import (
	"bytes"
	"errors"
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
	pattern := ""
	err := extractVar(varFilePattern, vars, func(val string) error {
		pattern = val
		return nil
	}, true)
	if err != nil {
		return false, err
	}

	var cmd *exec.Cmd
	err = extractVar(varCommand, vars, func(val string) error {
		parts := strings.Split(val, " ")
		if len(parts) == 0 {
			return errors.New("no executable command given")
		}
		command := parts[0]
		cargs := ""
		if len(parts) >= 2 {
			cargs = strings.Join(parts[1:], " ")
		}
		cmd = exec.Command(command, cargs)
		cmd.Dir = workingDir
		return nil
	}, true)

	var files []string
	err = extractVar(varFileList, vars, func(val string) error {
		files = strings.Split(val, " ")
		return nil
	}, false)
	if err != nil {
		return false, err
	}
	switch hook {
	case "pre-commit":
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
	default:
		return hookUnsupported(hook)
	}
}

func singleFileMatchesPattern(files []string, pattern string) bool {
	for _, file := range files {
		if err := matchString(file, pattern); err == nil {
			return true
		}
	}
	return false
}
