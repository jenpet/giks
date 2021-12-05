package config

import (
	"bytes"
	"errors"
	"github.com/jenpet/giks/args"
	"github.com/jenpet/giks/log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

// default filename for giks configs in case none is provided on invocation
const defaultGiksConfigFilename = "giks.yml"

// AssembleConfig takes giks specific arguments and parses the configuration file for giks in order to return a config.
// Additionally, it sanitizes the given inputs targeting files and returns a configuration which can be
// used without bothering about paths.
func AssembleConfig(ga args.GiksArgs) Config {
	cfg := parseConfigFile(ga.ConfigFile())
	cfg.GitDir = absoluteGitDirectory(ga.GitDir())
	cfg.WorkingDir = path.Dir(cfg.GitDir)
	cfg.Binary = absoluteBinaryPath(ga.Binary())
	return cfg
}

// absoluteBinaryPath determines the absolute path of the binary provided as a string.
// Three different scenarios are addressed:
// - binary string is already provided in an absolute way
// - binary string is provided relatively to the cwd
// - binary string is no file at all and presumably in the $PATH of the machine
func absoluteBinaryPath(binary string) string {
	if binary == "" {
		log.Error("Could not determine absolute path of binary. Provided binary name or filepath was empty.")
	}
	// check if the used binary file exists and might be relative or absolute
	if fi, _ := os.Stat(binary); fi != nil {
		// binary is already an absolute path
		if filepath.IsAbs(binary) {
			return binary
		}
		// get the absolute path to the binary
		path, err := filepath.Abs(binary)
		if err != nil {
			log.Errorf("Could not get absolute path to '%s' binary. Error: %+v", binary, err)
		}
		return path
	}

	// binary is presumably in the $PATH env var
	path, err := exec.LookPath(binary)
	if err != nil {
		log.Errorf("Could not get absolute path to '%s' binary. Error: %+v", binary, err)
	}
	return path
}

// absoluteConfigFile determines the absolute path of the provided configuration file.
// In case no file was provided it will assume that the default configuration file is used in the cwd.
// Absence of the configuration file will cause giks to exit.
func absoluteConfigFile(file string) string {
	// validate the given input file
	if file != "" {
		file = absoluteFilepath(file)
		if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
			log.Errorf("The provided config file '%s' does not exist", file)
		}
		return file
	}

	// use the default by utilizing the cwd
	path, err := os.Getwd()
	if err != nil {
		log.Error("Failed retrieving cwd. No config file provided. Fallback with default config not possible.")
	}
	file = absoluteFilepath(filepath.Join(path, defaultGiksConfigFilename))
	if _, err = os.Stat(file); errors.Is(err, os.ErrNotExist) {
		log.Errorf("No valid config file provided. Tried to read the default configuration '%s' but it does not exist.", file)
	}
	return file
}

// absoluteGitDirectory looks up the responsible git directory originating from a given directory. The absence of a
// directory results in a fallback to the absolute path to the cwd.
// Attention: the git command has to be present in the $PATH variable in order to identify the git directory.
func absoluteGitDirectory(dir string) string {
	if dir != "" {
		dir = absoluteFilepath(dir)
	} else {
		path, err := os.Getwd()
		if err != nil {
			log.Error("Failed retrieving cwd. No config file provided. Fallback with default config not possible.")
		}
		dir = absoluteFilepath(path)
	}

	// check git availability and get the git directory by executing
	// git -C <git-dir> rev-parse --git-dir
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--git-dir")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		log.Errorf("Failed validating git directory '%s'. Error: %+v", dir, err)
	}

	// if the output of the git command is not an absolute directory it is a child of the given dir
	gitDir := strings.TrimSpace(buf.String())
	if !filepath.IsAbs(gitDir) {
		dir = filepath.Join(dir, gitDir)
	} else {
		dir = gitDir
	}
	return dir
}

// absoluteFilepath returns the absolute path to a given file. Since '~' does not get resolved by the golang standard
// library it will is manually replaced within this function.
func absoluteFilepath(file string) string {
	if filepath.IsAbs(file) {
		return file
	}
	if strings.HasPrefix(file, "~") {
		u, err := user.Current()
		if err != nil {
			log.Errorf("Could not retrieve user home directory due to usage of '%s'. Error: %+v", file, err)
		}
		file = strings.Replace(file, "~", u.HomeDir, 1)
	}
	file, err := filepath.Abs(file)
	if err != nil {
		log.Errorf("Could not get absolute filepath for file '%s'. Error: %+v", file, err)
	}
	return file
}
