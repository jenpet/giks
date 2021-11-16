package config

import (
	"bytes"
	"errors"
	"fmt"
	"giks/args"
	"giks/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
)

// default filename for giks configs in case none is provided on invocation
const defaultGiksConfigFilename = "giks.yml"

func parseConfig(file string) Config {
	var cfg Config
	var err error
	var once sync.Once
	once.Do(func() {
		var b []byte
		cfgFile := validateConfigFile(file)
		b, err = ioutil.ReadFile(cfgFile)
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
		cfg.ConfigFile = cfgFile
	})
	if err != nil {
		fmt.Printf("Failed parsing giks configuration. Error: %s", err)
		os.Exit(1)
	}
	return cfg
}

func AssembleConfig(ga args.GiksArgs) Config {
	cfg := parseConfig(ga.ConfigFile())
	cfg.GitDir = validateGitDirectory(ga.GitDir())
	cfg.Binary = validateBinaryDirectory(ga.Binary())
	return cfg
}

func validateBinaryDirectory(binary string) string {
	// check if the used binary file exists and might be relative or absolute
	if fi, _ := os.Stat(binary); fi != nil {
		// binary is already an absolute path
		if filepath.IsAbs(binary) {
			return binary
		}
		// get the absolute path to the binary
		path, err := filepath.Abs(binary)
		if err != nil {
			log.Errorf("Could not get absolute path to giks binary. Error: %+v", err)
		}
		return path
	}

	// binary is presumably in the $PATH env var
	path, err := exec.LookPath(binary)
	if err != nil {
		log.Errorf("Could not get absolute path to giks binary. Error: %+v", err)
	}
	return path
}

// TODO: not really validate, rather a check for fallback
func validateConfigFile(file string) string {
	// validate the given input file
	if file != "" {
		file = absoluteFilepath(file)
		if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("The provided config file '%s' does not exist", file)
			os.Exit(1)
		}
		return file
	}

	// use the default by utilizing the cwd
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed retrieving cwd. No config file provided. Fallback with default config not possible.")
		os.Exit(1)
	}
	file = absoluteFilepath(filepath.Join(path, defaultGiksConfigFilename))
	return file
}

// TODO: not really validate, rather a check for fallback
func validateGitDirectory(dir string) string {
	if dir != "" {
		dir = absoluteFilepath(dir)
	} else {
		path, err := os.Getwd()
		if err != nil {
			fmt.Println("Failed retrieving cwd. No config file provided. Fallback with default config not possible.")
			os.Exit(1)
		}
		dir = absoluteFilepath(path)
	}

	// check git availability and get the git directory by executing
	// git -C <git-dir> rev-parse --git-dir
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--git-dir")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed validating git directory '%s'. Error: %+v", dir, err)
		os.Exit(1)
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

func absoluteFilepath(file string) string {
	if filepath.IsAbs(file) {
		return file
	}
	if strings.HasPrefix(file, "~" ) {
		u, err := user.Current()
		if err != nil {
			fmt.Printf("Could not retrieve user home directory due to usage of '%s'. Error: %+v", file, err)
			os.Exit(1)
		}
		file = strings.Replace(file, "~", u.HomeDir, 1)
	}
	file, err := filepath.Abs(file)
	if err != nil {
		fmt.Printf("Could not get absolute filepath for file '%s'. Error: %+v", file, err)
		os.Exit(1)
	}
	return file
}