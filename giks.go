package main

import (
	"flag"
	"fmt"
	"giks/commands"
	"giks/config"
	"os"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("Failed parsing giks configuration. Error: %s", err)
		return
	}

	if len(os.Args) < 2 {
		fmt.Printf("help text")
		os.Exit(1)
	}

	// subcommands
	switch os.Args[1] {
	case "hooks":
		commands.ProcessHooks(os.Args[2:], cfg)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}