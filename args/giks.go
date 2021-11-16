// Package args TODO: Document assumptions about the command, subcommand, hock and flag structure including global flags
package args

import (
	"fmt"
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
		return "help"
	}
	return ga.sanitizeArgs()[1]
}

func (ga GiksArgs) SubCommand() string {
	if len(ga.sanitizeArgs()) < 3 || isFlag(ga.sanitizeArgs()[2]) {
		return "help"
	}
	return ga.sanitizeArgs()[2]
}

func (ga GiksArgs) Hook() string {
	if len(ga.sanitizeArgs()) < 4 || isFlag(ga.sanitizeArgs()[3]) {
		return ""
	}
	return ga.sanitizeArgs()[3]
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
	args := ga.sanitizeArgs()
	if len(args) < 5 {
		return []string{}
	}
	return args[4:]
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


