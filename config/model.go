package config

import (
	"errors"
	"fmt"
	"giks/git"
	"giks/log"
)

// Config holds the config information provided by the used configuration file and additional
// meta information which is available at runtime.
type Config struct {
	// absolute path to the used configuration file
	ConfigFile string `yaml:"-"`
	// absolute path to the affected git repository
	GitDir string `yaml:"-"`
	// absolute path to the giks binary file
	Binary string `yaml:"-"`
	// parsed hook configurations based on the configuration file
	Hooks map[string]Hook `yaml:"hooks"`
}

func (c Config) HookList(all bool) map[string]Hook {
	hooks := map[string]Hook{}
	for name, h := range c.Hooks {
		// has to be valid and the all flag also returns disabled ones
		if h.validate() == nil && (all || h.Enabled) {
			hooks[name] = h
		}
	}
	return hooks
}

func (c Config) HookListNames(all bool) []string {
	list := c.HookList(all)
	names := make([]string, len(list))
	i := 0
	for name := range list {
		names[i] = name
		i++
	}
	return names
}

func (c Config) LookupHook(name string) (*Hook, error) {
	if name == "" {
		return nil, errors.New("provided hook is empty")
	}
	for n, h := range c.Hooks {
		if n == name {
			return &h, h.validate()
		}
	}
	return nil, errors.New("unknown hook")
}

func (c Config) Hook(name string) Hook {
	h, err := c.LookupHook(name)
	if err != nil {
		log.Errorf("error: could not find hook '%s'. Error: %+v", name, err)
	}
	return *h
}

func (c Config) validate() error {
	for name, h := range c.Hooks {
		if err := h.validate(); err != nil {
			return fmt.Errorf("hook '%s' is invalid: %s", name, err)
		}
	}
	return nil
}

type Hook struct {
	Enabled bool   `yaml:"enabled"`
	Steps   []Step `yaml:"steps"`
	Name    string `yaml:"-"`
}

func (h Hook) validate() error {
	valid := false
	for _, hook := range git.Hooks {
		if hook == h.Name {
			valid = true
		}
	}
	if !valid {
		return fmt.Errorf("hook '%s' is not a valid Git hook", h.Name)
	}
	return nil
}

func (h Hook) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	m["name"] = h.Name
	m["enabled"] = h.Enabled
	steps := make([]map[string]interface{}, len(h.Steps))
	for idx, step := range h.Steps {
		steps[idx] = step.ToMap()
	}
	m["steps"] = steps
	return m
}

type Step struct {
	Command string     `yaml:"command"`
	Exec    string     `yaml:"exec"`
	Script  string     `yaml:"script"`
	Plugin  PluginStep `yaml:"plugin"`
}

func (s Step) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	if s.Command != "" {
		m["command"] = s.Command
	}

	if s.Exec != "" {
		m["exec"] = s.Exec
	}

	if s.Script != "" {
		m["script"] = s.Script
	}

	if s.Plugin.Validate() == nil {
		vars := map[string]string{}
		for k, v := range s.Plugin.Vars {
			vars[k] = v
		}
		info := map[string]interface{}{}
		info["name"] = s.Plugin.Name
		info["vars"] = vars
		m["plugin"] = info
	}
	return m
}

func (s Step) validate() error {
	i := 0
	if s.Command != "" {
		i++
	}

	if s.Exec != "" {
		i++
	}

	if s.Script != "" {
		i++
	}

	if s.Plugin.Validate() != nil {
		i++
	}

	if i != 1 {
		return errors.New("too many or too few step entry-points provided. Only one of 'command', 'exec', 'script' or 'plugin' is possible")
	}
	return nil
}

type PluginStep struct {
	Name string            `yaml:"name"`
	Vars map[string]string `yaml:"vars"`
}

func (ps PluginStep) Validate() error {
	if ps.Name == "" {
		return errors.New("provided plugin name is empty")
	}
	return nil
}
