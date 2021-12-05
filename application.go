package main

import (
	"github.com/jenpet/giks/args"
	"github.com/jenpet/giks/commands"
	"github.com/jenpet/giks/config"
	"os"
)

func main() {
	// parse into specific giks arguments to ease command, subcommand and argument handling
	var ga args.GiksArgs = os.Args
	cfg := config.AssembleConfig(ga)
	commands.Process(cfg, ga)
}
