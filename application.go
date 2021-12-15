package main

import (
	"github.com/jenpet/giks/args"
	"github.com/jenpet/giks/commands"
	"github.com/jenpet/giks/config"
	"github.com/jenpet/giks/log"
	"os"
)

func main() {
	// parse into specific giks arguments to ease command, subcommand and argument handling
	var ga args.GiksArgs = os.Args
	cfg := config.AssembleConfig(ga)
	// initialize the logger in case debug logging is required
	log.Init(ga.Debug())
	commands.Process(cfg, ga)
}
