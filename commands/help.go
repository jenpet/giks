package commands

import (
	"flag"
	gargs "github.com/jenpet/giks/args"
	"github.com/jenpet/giks/cli"
	"github.com/jenpet/giks/config"
	"strings"
	"text/template"
)

var helpCommand = flag.NewFlagSet("help", flag.ExitOnError)
var debugAttr = helpCommand.Bool("debug", false, "show additional debug information")

var helpTemplateString = `

Usage: giks COMMAND [SUBCOMMAND] [OPTIONS]

Global Options:

--config		Path to the giks configuration (default: ${PWD}/giks.yml)
--git-dir		Path to the Git directory which should be managed by giks (default: ${PWD}/.git)


Commands:

install [HOOK] Installs a given hook based on the configuration into the target directory. 
	If no hook is provided it will install all enabled hooks of the configuration.

uninstall [HOOK] Removes a given hook based on the configuration from the target directory. 
	If no hook is provided all hooks will be removed.

exec HOOK Executes a given hook according to the configuration provided.

show [HOOK] [--all] Displays detailed information about the used configuration (i.e. list of hooks). 
	If a hook is provided it will show the details for the specific hook. Adding the --all flag also lists disabled hooks.

{{ if .debug }}
Binary:		{{ .debug.binary }}
Config:		{{ .debug.config }}
Git directory:		{{ .debug.gitdir }}
Arguments:		{{ .debug.args }}
{{- end }}
`

var helpTemplate *template.Template

func init() {
	var err error
	helpTemplate, err = template.New("help").Parse(helpTemplateString)
	if err != nil {
		panic(err)
	}
}

func printHelp(cfg config.Config, gargs gargs.GiksArgs) {
	debug := map[string]string{
		"binary": cfg.Binary,
		"config": cfg.ConfigFile,
		"gitdir": cfg.GitDir,
		"args":   strings.Join(gargs.Args(), ""),
	}
	var data map[string]interface{}
	_ = helpCommand.Parse(gargs.Args())
	if helpCommand.Parsed() && *debugAttr {
		data = map[string]interface{}{"debug": debug}
	}
	cli.PrintTemplate(helpTemplate, data)
}
