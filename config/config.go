package config

import (
	"context"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"sync"
)

type ctxKey int

const (
	ctxHookKey ctxKey = iota
	ctxConfigKey
)

var validHooks = []string{
	"applypatch-msg",
	"commit-msg",
	"fsmonitor-watchman",
	"post-update",
	"pre-applypatch",
	"pre-commit",
	"pre-merge-commit",
	"pre-push",
	"pre-rebase",
	"pre-receive",
	"prepare-commit-msg",
	"update",
}

func GetConfig() Config {
	var cfg Config
	var err error
	var once sync.Once
	once.Do(func() {
		var b []byte
		b, err = ioutil.ReadFile("config.yml")
		if err != nil {
			return
		}
		err = yaml.Unmarshal(b, &cfg)
		if err != nil {
			err = cfg.validate()
		}
		for name, hook := range cfg.Hooks {
			hook.Name = name
			cfg.Hooks[name]= hook
		}
	})
	if err != nil {
		fmt.Printf("Failed parsing giks configuration. Error: %s", err)
		os.Exit(1)
	}
	return cfg
}

type Config struct {
	Hooks map[string]Hook `yaml:"hooks"`
}

func (c Config) HookList(all bool) map[string]Hook {
	hooks := map[string]Hook{}
	for name, h := range c.Hooks {
		// has to be valid and the all flag also returns disabled ones
		if h.validate() == nil && (all || h.Enabled)  {
			hooks[name] = h
		}
	}
	return hooks
}

func (c Config) Hook(name string) (*Hook, error) {
	for n, h := range c.Hooks {
		if n == name {
			return &h, h.validate()
		}
	}
	return nil, errors.New("unknown hook")
}

func (c Config) validate() error {
	for name,h := range c.Hooks {
		if err := h.validate(); err != nil {
			return fmt.Errorf("hook '%s' is invalid: %s", name, err)
		}
	}
	return nil
}

type Hook struct {
	Enabled bool `yaml:"enabled"`
	Steps []Step `yaml:"steps"`
	Name string `yaml:"-"`
}

func (h Hook) validate() error {
	valid := false
	for _, hook := range validHooks {
		if hook == h.Name {
			valid = true
		}
	}
	if !valid {
		return fmt.Errorf("hook '%s' is not a valid Git hook", h.Name)
	}
	if h.Enabled && len(h.Steps) <= 0 {
		return errors.New("hook enabled but validate steps are missing")
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

// ContextWithHook adds the targeted hook into the context to provide easy access later on
func ContextWithHook(ctx context.Context, args []string) context.Context {
	if len(args) < 1 {
		fmt.Println("hook name is required but missing")
		os.Exit(1)
	}
	cfg := GetConfig()
	hook, err := cfg.Hook(args[0])
	if err != nil {
		fmt.Printf("failed retrieving '%s' hook. Error: %s\n", args[0], err)
		os.Exit(1)
	}
	return context.WithValue(ctx, ctxHookKey, *hook)
}

func HookFromContext(ctx context.Context) Hook {
	if ctx == nil {
		fmt.Println("could not retrieve hook from context context was nil.")
		os.Exit(1)
	}
	if h, ok := ctx.Value(ctxHookKey).(Hook); ok {
		return h
	}
	fmt.Println("could not retrieve hook from context")
	os.Exit(1)
	return Hook{}
}

// ContextWithConfig adds the configuration to the context to provide easy access later on
func ContextWithConfig(ctx context.Context) context.Context {
	cfg := GetConfig()
	return context.WithValue(ctx, ctxConfigKey, cfg)
}

// ConfigFromContext returns the configuration from a given context
func ConfigFromContext(ctx context.Context) Config {
	if ctx == nil {
		fmt.Println("could not retrieve config from context context was nil.")
		os.Exit(1)
	}
	if cfg, ok := ctx.Value(ctxConfigKey).(Config); ok {
		return cfg
	}
	fmt.Println("could not retrieve config from context")
	os.Exit(1)
	return Config{}
}

type Step struct {
	Command string `yaml:"command"`
	Exec string `yaml:"exec"`
	Script string     `yaml:"script"`
	Plugin PluginStep `yaml:"plugin"`
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
		for k,v := range s.Plugin.Vars {
			vars[k]=v
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
	Name string `yaml:"name"`
	Vars map[string]string `yaml:"vars"`
}

func (ps PluginStep) Validate() error {
	if ps.Name == "" {
		return errors.New("provided plugin name is empty")
	}
	return nil
}