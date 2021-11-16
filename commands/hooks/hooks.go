package hooks

import (
	"context"
	"errors"
	"flag"
	"fmt"
	gargs "giks/args"
	"giks/config"
	"giks/util"
	"github.com/mattn/go-shellwords"
	"os"
	"os/exec"
	"syscall"
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
  	{{ if $step.plugin.args }}vars:
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

func ProcessHooks(ctx context.Context, gargs gargs.GiksArgs) {
	// actual array of arguments without the binary itself the command and subcommand
	args := gargs.Args()
	// retrieve provided config from context
	cfg := config.ConfigFromContext(ctx)
	switch gargs.SubCommand() {
	case "list":
		_ = listCommand.Parse(args)
		if listCommand.Parsed() {
			all := *listAllAttr
			util.PrintTemplate(listTemplate, cfg.HookList(all))
		}
	case "show":
		util.PrintTemplate(detailsTemplate, cfg.Hook(gargs.Hook()).ToMap())
	case "exec":
		if err := executeHook(cfg.Hook(gargs.Hook()), args); err != nil {
			fmt.Printf("failed executing '%s' hook. Error: %s\n", args[0], err)
			os.Exit(1)
		}
	case "install":
		if gargs.HasHook() {
			h := cfg.Hook(gargs.Hook())
			installSingleHook(cfg, h, true)
			break
		}
		installHookList(cfg)
	case "uninstall":
		if gargs.HasHook() {
			h := cfg.Hook(gargs.Hook())
			uninstallSingleHook(cfg, h, true)
			break
		}
		uninstallHookList(cfg)
	case "help":
		fmt.Println("help text")
	default:
		fmt.Printf("Unknown subcommand '%s'", gargs.SubCommand())
	}
}

func executeHook(h config.Hook, args []string) error {
	if !h.Enabled {
		return fmt.Errorf("hook '%s' is not enabled", h.Name)
	}
	for i, step := range h.Steps {
		if err := executeStep(h, step, args); err != nil {
			return fmt.Errorf("failed executing step no. %d. Error: %s", i+1, err)
		}
	}
	return nil
}

func executeStep(h config.Hook, s config.Step, args []string) error {
	if s.Script != "" {
		return executeScript(h.Name, s.Script, args, nil)
	}

	if s.Command != "" {
		return executeCommand(h.Name, s.Command, args)
	}

	if s.Exec != "" {
		return runExec(h.Name, s.Exec, args)
	}

	if err := s.Plugin.Validate(); err == nil {
		return executePlugin(h.Name, s.Plugin, args)
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

func executePlugin(hook string, plugin config.PluginStep, args []string) error {
	pluginPath := fmt.Sprintf("./%s/%s.sh", pluginDirectory, plugin.Name)
	err := executeScript(hook, pluginPath, args, plugin.Vars)
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("the given plugin '%s' does not exist", plugin.Name)
	}
	return err
}

func executeCommand(hook string, command string, args []string) error {
	args = append([]string{"-c", command}, args...)
	cmd := exec.Command("sh", args...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("GIKS_HOOK_TYPE=%s", hook))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runExec(hook string, line string, args []string) error {
	parts, err := shellwords.Parse(line)
	if err != nil {
		return fmt.Errorf("could not parse exec '%s'", line)
	}

	bin := parts[0]
	path, err := exec.LookPath(bin)
	if err != nil {
		return fmt.Errorf("binary not found for exec '%s'", line)
	}
	env := os.Environ()
	env = append(env, fmt.Sprintf("GIKS_HOOK_TYPE=%s", hook))
	args = append(parts, args...)
	return syscall.Exec(path, args, env)
}

func envsToList(envs map[string]string) []string {
	var list []string
	for k,v := range envs {
		list = append(list, fmt.Sprintf("%s=%s", k,v))
	}
	return list
}