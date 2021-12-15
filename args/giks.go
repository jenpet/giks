// Package args TODO: Document assumptions about the command, subcommand, hook and flag structure including global flags
package args

import (
	"github.com/jenpet/giks/git"
	"strings"
)

const (
	keyGlobalConfigFlag = "--config"
	keyGlobalGitDirFlag = "--git-dir"
	keyGlobalDebugFlag = "--debug"
)

var globalFlags = []string{keyGlobalGitDirFlag, keyGlobalConfigFlag, keyGlobalDebugFlag}

type GiksArgs []string

func (ga GiksArgs) Binary() string {
	return ga[0]
}

func (ga GiksArgs) Command() string {
	if len(ga.sanitizeArgs()) < 2 || isFlag(ga.sanitizeArgs()[1]) {
		return ""
	}
	return ga.sanitizeArgs()[1]
}

func (ga GiksArgs) SubCommand() string {
	// if there are not enough arguments, the subcommand argument is a flag or a valid git hook treat consider the subcommand absent
	sargs := ga.sanitizeArgs()
	if len(sargs) < 3 || isFlag(sargs[2]) || git.IsValidHook(sargs[2]) {
		return ""
	}
	return ga.sanitizeArgs()[2]
}

func (ga GiksArgs) Hook() string {
	for _, arg := range ga.sanitizeArgs() {
		if git.IsValidHook(arg) {
			return arg
		}
	}
	return ""
}

func (ga GiksArgs) HasHook() bool {
	return ga.Hook() != ""
}

func (ga GiksArgs) ConfigFile() string {
	v, _ := ga.globalFlag(keyGlobalConfigFlag)
	return v
}

func (ga GiksArgs) GitDir() string {
	v, _ := ga.globalFlag(keyGlobalGitDirFlag)
	return v
}

func (ga GiksArgs) Debug() bool {
	_, ok := ga.globalFlag(keyGlobalDebugFlag)
	if !ok {
		return false
	}
	// debug flag does not need a value
	return true
}

func (ga GiksArgs) globalFlag(flag string) (string, bool) {
	for _, arg := range ga {
		// flag is set in general
		if fileArg := strings.Split(arg, flag); len(fileArg) == 2 {
			// flag is set with an actual value
			if val := strings.Split(fileArg[1], "="); len(val) == 2 {
				return val[1], true
			}
			return "", true
		}
	}
	return "", false
}

func (ga GiksArgs) Args() []string {
	args := []string{}
	for _, arg := range ga.sanitizeArgs() {
		if isFlag(arg) {
			args = append(args, arg)
		}
	}
	return args
}

// sanitizeArgs removes all arguments relevant for a global configuration
func (ga GiksArgs) sanitizeArgs() []string {
	var sanatized []string
OUTER:
	for _, arg := range ga {
		for _, flag := range globalFlags {
			if strings.HasPrefix(arg, flag) {
				continue OUTER
			}
		}
		sanatized = append(sanatized, arg)
	}
	return sanatized
}

func (ga GiksArgs) Raw() []string {
	return ga
}

func isFlag(s string) bool {
	return strings.HasPrefix(s, "--")
}
