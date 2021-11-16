package hooks

import (
	"bytes"
	"errors"
	"fmt"
	"giks/config"
	"giks/git"
	"giks/log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// hookMask holds the access mask for installed hooks
const hookMask = 0755

var (
	errHookExternallyManaged = errors.New("hook is externally managed")
	errHookAlreadyInstalled = errors.New("hook is already installed")
	errHookNotInstalled = errors.New("hook is not installed")
)

var hookTemplateString = `
# GIKS-ZONE!
# This {{ .name }} hook is managed via giks (https://github.com/jenpet/giks).
# You should not alter this file manually except you do it tenderly and know what you are actually doing.
# To remove this hook run 'giks hooks uninstall {{ .name }}'.
{{ .command }}
`

func installSingleHook(cfg config.Config, h config.Hook, confirmation bool) {
	if confirmation {
		verifyUserConfirmation(fmt.Sprintf("Do you want to install hook '%s' for git directory '%s'", h.Name, cfg.GitDir))
	}
	if err := installHook(cfg, h.Name, false); err != nil {
		if errors.Is(err, errHookAlreadyInstalled) || errors.Is(err, errHookExternallyManaged){
			log.Warnf("Hook '%s' was not installed. Reason: %+v", h.Name, err)
			return
		}
		log.Errorf("Hook '%s' could not be installed. Error: %s", h.Name, err)
	}
}

func uninstallSingleHook(cfg config.Config, h config.Hook, confirmation bool) {
	if confirmation {
		verifyUserConfirmation(fmt.Sprintf("Do you want to uninstall hook '%s' for git directory '%s'", h.Name, cfg.GitDir))
	}
	if err := uninstallHook(cfg, h.Name); err != nil {
		if errors.Is(err, errHookNotInstalled) || errors.Is(err, errHookExternallyManaged){
			log.Warnf("Hook '%s' was not uninstalled. Reason: %+v", h.Name, err)
			return
		}
		log.Errorf("Hook '%s' could not be uninstalled. Error: %+v", h.Name, err)
	}
}

func installHookList(cfg config.Config) {
	msg := fmt.Sprintf("Do you want to install the '%s' hook(s) for git directory '%s'",
		strings.Join(cfg.HookListNames(false), ", "),
		cfg.GitDir)
	verifyUserConfirmation(msg)
	for _, h := range cfg.HookList(false) {
		installSingleHook(cfg, h, false)
	}
}

func uninstallHookList(cfg config.Config) {
	msg := fmt.Sprintf("Do you want to uninstall the '%s' hook(s) for git directory '%s'",
		strings.Join(cfg.HookListNames(false), ", "),
		cfg.GitDir)
	verifyUserConfirmation(msg)
	for _, h := range cfg.HookList(false) {
		uninstallSingleHook(cfg, h, false)
	}
}

func installHook(cfg config.Config, hookName string, force bool) error {
	ok, err := hookIsInstalled(cfg, hookName)
	if err != nil {
		return err
	}
	if ok && !force {
		return errHookAlreadyInstalled
	}
	fileName := hookFileName(cfg.GitDir, hookName)
	content := hookFileContent(cfg, hookName)
	err = os.WriteFile(fileName, []byte(content), hookMask)
	if err != nil {
		log.Errorf("failed writing hook file '%s'. Error: %+v", fileName, err)
	}
	log.Infof("Installed hook '%s' in '%s'", hookName, fileName)
	return nil
}

func uninstallHook(cfg config.Config, hookName string) error {
	ok, err := hookIsInstalled(cfg, hookName)
	if err != nil  {
		return err
	}
	if !ok {
		return errHookNotInstalled
	}
	fileName := hookFileName(cfg.GitDir, hookName)
	if err = os.Remove(fileName); err != nil {
		log.Errorf("failed removing hook file '%s'. Error: %+v", fileName, err)
	}
	log.Infof("Uninstalled hook '%s' by removing '%s'", hookName, fileName)
	return nil
}

func hookIsInstalled(cfg config.Config, hookName string) (bool, error) {
	file := hookFileName(cfg.GitDir, hookName)
	if _, err := os.Stat(file); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		log.Errorf("failed checking hook installation for hook '%s'. Error: %+v", hookName, err)
	}
	content := hookFileContent(cfg, hookName)
	b, err := os.ReadFile(file)
	if err != nil {
		log.Errorf("could not read hook file '%s'. Error: %+v", file, err)
	}
	if strings.TrimSpace(string(b)) != content {
		return true, errHookExternallyManaged
	}
	return true, nil
}

func hookFileName(gitDir string, hookName string) string {
	return filepath.Join(gitDir, "hooks", hookName)
}

func hookFileContent(cfg config.Config, hookName string) string {
	cmd, err := commandString(cfg, hookName)
	if err != nil {
		log.Errorf("could not retrieve command string for hook '%s'. Error: %+v", hookName, err)
	}
	var content bytes.Buffer
	tpl, _ := template.New("hook").Parse(hookTemplateString)
	data := map[string]string {
		"name": hookName,
		"command": cmd,
	}
	_ = tpl.Execute(&content, data)
	return strings.TrimSpace(content.String())
}

func commandString(cfg config.Config, hookName string) (string, error) {
	cmd := fmt.Sprintf("%s hooks exec %s --config=%s", cfg.Binary, hookName, cfg.ConfigFile)
	switch hookName {
	case git.HookCommitMsg:
		return fmt.Sprintf("%s ${1}", cmd), nil
	case git.HookPreCommit:
		return cmd, nil
	}
	return "", errors.New(fmt.Sprintf("installation with hook '%s' is not supported", hookName))
}

func verifyUserConfirmation(msg string) {
	log.Infof("%s (y/n)?\n", msg)
	if !readUserConfirmation() {
		log.Info("Operation cancelled due to user selection.")
		os.Exit(1)
	}
}

func readUserConfirmation() bool {
	var confirmation string
	_, err := fmt.Scanf("%s", &confirmation)
	if err != nil {
		log.Errorf("failed reading user input. Error: %+v", err)
	}
	return "y" == strings.ToLower(confirmation)
}
