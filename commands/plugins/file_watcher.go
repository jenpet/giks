package plugins

import (
	"errors"
	"fmt"
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

func (fw FileWatcher) Run(hook string, vars map[string]string, args []string) (bool, error) {
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
			err = cmd.Run()
			if err != nil {
				return false, fmt.Errorf("files matched pattern but command failed: %+v", err)
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