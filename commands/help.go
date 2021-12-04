package commands

import (
	"flag"
	gargs "giks/args"
	"giks/cli"
	"giks/config"
	"strings"
	"text/template"
)

var helpCommand = flag.NewFlagSet("help", flag.ExitOnError)
var debugAttr = helpCommand.Bool("debug", false, "show additional debug information")

var helpTemplateString = `
HALP TEXT
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
