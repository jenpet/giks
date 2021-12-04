// Package args TODO: Document assumptions about the command, subcommand, hook and flag structure including global flags
package args

import (
	"fmt"
	"giks/git"
	"strings"
)

const (
	keyGlobalConfigFlag = "--config"
	keyGlobalGitDirFlag = "--git-dir"
)

var globalFlags = []string{keyGlobalGitDirFlag, keyGlobalConfigFlag}

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
	return ga.getGlobalFlag(keyGlobalConfigFlag)
}

func (ga GiksArgs) GitDir() string {
	return ga.getGlobalFlag(keyGlobalGitDirFlag)
}

func (ga GiksArgs) getGlobalFlag(flag string) string {
	for _, arg := range ga {
		if fileArg := strings.Split(arg, fmt.Sprintf("%s=", flag)); len(fileArg) == 2 {
			return fileArg[1]
		}
	}
	return ""
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
