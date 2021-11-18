package main

import (
	"fmt"
	"giks/args"
	"giks/commands"
	"giks/config"
	"giks/log"
	"os"
)

func main() {
	// parse into specific giks arguments to ease command, subcommand and argument handling
	var ga args.GiksArgs = os.Args
	cfg := config.AssembleConfig(ga)
	switch ga.Command() {
	case "hooks":
		commands.ProcessHooks(cfg, ga)
	case "help":
		fmt.Println("Help text")
		fmt.Printf("giks binary: '%s', giks config: '%s', git directory: '%s'", cfg.Binary, cfg.ConfigFile, cfg.GitDir)
	default:
		log.Errorf("Command '%s' is unknown. Run `giks help` in order to get some ideas on the hook game.")
	}
}