package commands

import (
	"errors"
	"flag"
	"fmt"
	"giks/config"
	"os"
	"os/exec"
	"text/template"
)

const pluginDirectory = "plugins"

var listCommand = flag.NewFlagSet("list", flag.ExitOnError)
var listAllAttr = listCommand.Bool("all", false, "include disabled hooks")

var listTemplateString = `
HOOK				| ENABLED				| STEPS
{{- range $hook, $value := . }}
{{ $hook }}				| {{ $value.Enabled }}				| {{ len $value.Steps -}}
{{ end }}
`

var detailsTemplateString = `
HOOK: {{ .name }}
ENABLED: {{ .enabled }}
STEPS: {{ len .steps }}
{{- range $idx, $step := .steps }}
  {{ if $step.command }}{{ $idx }}.)	command: '{{ $step.command }}'
  {{- else if $step.exec }}{{ $idx }}.)	exec: '{{ $step.exec }}'
  {{- else if $step.script }}{{ $idx }}.)	script: '{{ $step.script }}'
  {{- else if $step.plugin }}{{ $idx }}.)	plugin: '{{ $step.plugin.name }}'
  	{{ if $step.plugin.args }}args:
    	{{- range $key, $value := $step.plugin.args }}
	  - {{ $key }} = {{ $value }}
    	{{- end }}
    {{ end }}
  {{ end -}}
{{ end -}}
`


var listTemplate *template.Template
var detailsTemplate *template.Template

func init() {
	var err error
	listTemplate, err = template.New("list").Parse(listTemplateString)
	detailsTemplate, err = template.New("details").Parse(detailsTemplateString)

	if err != nil {
		panic(err)
	}
}

func ProcessHooks(args []string, cfg config.Config) {
	switch args[0] {
		case "list":
			_ = listCommand.Parse(args[1:])
			if listCommand.Parsed() {
				all := *listAllAttr
				printTemplate(listTemplate, cfg.HookList(all))
			}
		case "show":
			if len(args[1:]) < 1 {
				fmt.Println("missing hook name")
				os.Exit(1)
			}
			hook, err := cfg.Hook(args[1])
			if err != nil {
				fmt.Printf("failed retrieving '%s' hook. Error: %s\n", args[1], err)
				os.Exit(1)
			}

			printTemplate(detailsTemplate, hook.ToMap())
	case "exec":
		if len(args[1:]) < 1 {
			fmt.Println("missing hook name")
			os.Exit(1)
		}
		hook, err := cfg.Hook(args[1])
		if err != nil {
			fmt.Printf("failed retrieving '%s' hook. Error: %s\n", args[1], err)
			os.Exit(1)
		}
		if err := executeHook(*hook, args[2:]); err != nil {
			fmt.Printf("failed executing '%s' hook. Error: %s\n", args[1], err)
			os.Exit(1)
		}
	}
}

func executeHook(h config.Hook, args []string) error {
	if !h.Enabled {
		return fmt.Errorf("hook '%s' is not enabled", h.Name)
	}
	for i, step := range h.Steps {
		if err := executeStep(h.Name, step, args); err != nil {
			return fmt.Errorf("failed executing step no. %d. Error: %s", i+1, err)
		}
	}
	return nil
}

func executeStep(hook string, s config.Step, args []string) error {
	if s.Script != "" {
		return executeScript(hook, s.Script, args, nil)
	}
	if err := s.Plugin.Validate(); err == nil {
		return executePlugin(hook, s.Plugin)
	}
	return errors.New("step seems to be invalid")
}

func executeScript(hook string, path string, args []string, envs map[string]string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	var bin = path

	// check whether the script is executable. If not use the shell to execute the script
	if !stat.IsDir() && stat.Mode()&0100 == 0 {
		bin = "sh"
		args = append([]string{path}, args...)
	}

	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), envsToList(envs)...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("GIKS_HOOK_TYPE=%s", hook))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func executePlugin(hook string, plugin config.PluginStep) error {
	pluginPath := fmt.Sprintf("./%s/%s.sh", pluginDirectory, plugin.Name)
	err := executeScript(hook, pluginPath, nil, plugin.Args)
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("the given plugin '%s' does not exist", plugin.Name)
	}
	return err
}

func envsToList(envs map[string]string) []string {
	var list []string
	for k,v := range envs {
		list = append(list, fmt.Sprintf("%s=%s", k,v))
	}
	return list
}