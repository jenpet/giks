package args

import "strings"

type GiksArgs []string

func (ga GiksArgs) Command() string {
	if len(ga) < 2 || isFlag(ga[1]) {
		return "help"
	}
	return ga[1]
}

func (ga GiksArgs) SubCommand() string {
	if len(ga) < 3 || isFlag(ga[2]) {
		return "help"
	}
	return ga[2]
}

func (ga GiksArgs) Args() []string {
	return ga[3:]
}

func (ga GiksArgs) Raw() []string {
	return ga
}

func isFlag(s string) bool {
	return strings.HasPrefix(s, "--")
}


