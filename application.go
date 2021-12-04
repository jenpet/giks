package main

import (
	"giks/args"
	"giks/commands"
	"giks/config"
	"os"
)

func main() {
	// parse into specific giks arguments to ease command, subcommand and argument handling
	var ga args.GiksArgs = os.Args
	cfg := config.AssembleConfig(ga)
	commands.Process(cfg, ga)
}
