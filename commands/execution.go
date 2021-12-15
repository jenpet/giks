package commands

import (
	"fmt"
	gargs "github.com/jenpet/giks/args"
	"github.com/jenpet/giks/commands/plugins"
	"github.com/jenpet/giks/config"
	"github.com/jenpet/giks/errors"
	"github.com/jenpet/giks/git"
	"github.com/jenpet/giks/log"
	"github.com/mattn/go-shellwords"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"text/template"
)

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
	for i, step := range h.Steps {
		// ensure that the variables are up-to-date for every step in case they changed
		// due to previous steps
		vars := giksVars(cfg, gargs)
		if err := executeStep(cfg.WorkingDir, h, step, gargs, vars); err != nil {
			if errors.IsWarningError(err) {
				log.Warnf("failed executing step no. %d. Error: %s", i+1, err)
				continue
			}
			return fmt.Errorf("failed executing step no. %d. Error: %s", i+1, err)
		}
	}
	return nil
}

func executeStep(workingDir string, h config.Hook, s config.Step, args []string, vars map[string]string) error {
	if s.Script != "" {
		return executeScript(workingDir, s.Script, args, vars)
	}

	if s.Command != "" {
		return executeCommand(workingDir, s.Command, args, vars)
	}

	if s.Exec != "" {
		return runExec(workingDir, s.Exec, args, vars)
	}

	if err := s.Plugin.Validate(); err == nil {
		return executePlugin(workingDir, h.Name, s.Plugin, args, vars)
	}

	return errors.New("step seems to be invalid")
}

func executeScript(workingDir string, path string, args []string, vars map[string]string) error {
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
	cmd.Dir = workingDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// TODO: aside from the pre-compiled built-in plugins also support a plugin directory containing bash scripts which can
// TODO: Allow Env variables and commands to be plugin arguments
// be re-used to avoid copy & paste code within the config file
func executePlugin(workingDir string, hook string, pCfg config.PluginStep, args []string, vars map[string]string) error {
	p, err := plugins.Get(pCfg.Name)
	if err != nil {
		return err
	}
	// merge plugin variables with vars given by giks itself and replace values of the plugin with the present value
	// of the giks variable
	for k, v := range pCfg.Vars {
		// in case a value of a plugin variable is the key of a giks variable
		// replace it accordingly
		if val, ok := vars[strings.TrimSpace(v)]; ok {
			vars[k] = val
			continue
		}
		vars[k] = v
	}
	exit, err := p.Run(workingDir, hook, vars, args)
	if err != nil {
		// error message was provided by the plugin configuration use the provided one
		msg := pCfg.ErrorMessage
		// default error message with error
		if msg == "" {
			msg = fmt.Sprintf("failed executing plugin '%s': %+v", pCfg.Name, err)
		}
		// return an error which forces no exit
		if !exit {
			return errors.NewWarningError(msg)
		}
		// error that forces an exit
		return errors.New(msg)
	}
	if pCfg.SuccessMessage != "" {
		log.Info(pCfg.SuccessMessage)
	}
	return err
}

func executeCommand(workingDir string, command string, args []string, vars map[string]string) error {
	args = append([]string{"-c", command}, args...)
	cmd := exec.Command("sh", args...)
	cmd.Dir = workingDir
	cmd.Env = append(cmd.Env, varsToList(vars)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runExec(workingDir string, line string, args []string, vars map[string]string) error {
	if err := os.Chdir(workingDir); err == nil {
		return fmt.Errorf("could not change into working directory '%s'. Error: %+v", workingDir, err)
	}
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
		// TODO: ugly since it assumes that mixins always are arrays
		if strings.Contains(k, "MIXIN") {
			// set the variable as an array
			list = append(list, fmt.Sprintf("%s=(%s)", k, v))
			continue
		}
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
