package commands

import (
	"errors"
	"flag"
	"fmt"
	gargs "giks/args"
	"giks/commands/plugins"
	"giks/config"
	"giks/git"
	"giks/log"
	"github.com/mattn/go-shellwords"
	"os"
	"os/exec"
	"path"
	"syscall"
	"text/template"
)

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
  	{{ if $step.plugin.vars }}vars:
    	{{- range $key, $value := $step.plugin.vars }}
	  - {{ $key }} = {{ $value }}
    	{{- end }}
    {{- end }}
  {{- end }}
{{- end }}
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

func executeHook(cfg config.Config, gargs gargs.GiksArgs) error {
	h := cfg.Hook(gargs.Hook())
	if !h.Enabled {
		return fmt.Errorf("hook '%s' is not enabled", h.Name)
	}
	vars := giksVars(cfg, gargs)
	for i, step := range h.Steps {
		if err := executeStep(h, step, gargs, vars); err != nil {
			return fmt.Errorf("failed executing step no. %d. Error: %s", i+1, err)
		}
	}
	return nil
}

func executeStep(h config.Hook, s config.Step, args []string, vars map[string]string) error {
	if s.Script != "" {
		return executeScript(s.Script, args, vars)
	}

	if s.Command != "" {
		return executeCommand(s.Command, args, vars)
	}

	if s.Exec != "" {
		return runExec(s.Exec, args, vars)
	}

	if err := s.Plugin.Validate(); err == nil {
		return executePlugin(h.Name, s.Plugin, args, vars)
	}

	return errors.New("step seems to be invalid")
}

func executeScript(path string, args []string, vars map[string]string) error {
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
	cmd.Env = append(os.Environ(), varsToList(vars)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// TODO: aside from the pre-compiled built-in plugins also support a plugin directory containing bash scripts which can
// be re-used to avoid copy & paste code within the config file
func executePlugin(hook string, pCfg config.PluginStep, args []string, vars map[string]string) error {
	p, err := plugins.Get(pCfg.Name)
	if err != nil {
		return err
	}
	// merge plugin variables with vars given by giks itself
	for k, v := range pCfg.Vars {
		vars[k] = v
	}
	exit, err := p.Run(hook, vars, args)
	if err != nil && exit {
		log.Errorf("failed executing plugin '%s'. Error: %+v", hook, err)
	}
	return err
}

func executeCommand(command string, args []string, vars map[string]string) error {
	args = append([]string{"-c", command}, args...)
	cmd := exec.Command("sh", args...)
	cmd.Env = append(cmd.Env, varsToList(vars)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runExec(line string, args []string, vars map[string]string) error {
	parts, err := shellwords.Parse(line)
	if err != nil {
		return fmt.Errorf("could not parse exec '%s'", line)
	}

	bin := parts[0]
	path, err := exec.LookPath(bin)
	if err != nil {
		return fmt.Errorf("binary not found for exec '%s'", line)
	}
	env := append(os.Environ(), varsToList(vars)...)
	args = append(parts, args...)
	return syscall.Exec(path, args, env)
}

func varsToList(envs map[string]string) []string {
	var list []string
	for k, v := range envs {
		list = append(list, fmt.Sprintf("%s=%s", k, v))
	}
	return list
}

func giksVars(cfg config.Config, gargs gargs.GiksArgs) map[string]string {
	vars := map[string]string{}
	git.ApplyMixins(path.Dir(cfg.GitDir), vars)
	vars["GIKS_HOOK_TYPE"] = gargs.Hook()
	return vars
}
