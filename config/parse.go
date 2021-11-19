package config

import (
	"bytes"
	"fmt"
	"giks/log"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

func parseConfig(r io.Reader) (*Config, error) {
	var err error
	buf := bytes.Buffer{}
	if _, err = buf.ReadFrom(r); err != nil {
		return nil, parseError(err.Error())
	}
	var cfg Config
	if err = yaml.Unmarshal(buf.Bytes(), &cfg); err != nil {
		return nil, parseError(err.Error())
	}
	for name, hook := range cfg.Hooks {
		hook.Name = name
		cfg.Hooks[name] = hook
	}
	if err = cfg.validate(); err != nil {
		return nil, parseError(err.Error())
	}
	return &cfg, err
}

func parseConfigFile(file string) Config {
	absFile := absoluteConfigFile(file)
	fh, err := os.Open(absFile)
	if err != nil {
		log.Errorf("Failed accessing configuration file. Error: %s", err)
	}
	cfg, err := parseConfig(fh)
	cfg.ConfigFile = absFile
	if err != nil {
		log.Errorf("Failed parsing provided configuration. Error: %s", err)
	}
	return *cfg
}

func parseError(reason string) error {
	return fmt.Errorf("provided configuration malformed, reason: %s", reason)
}
