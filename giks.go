package main

import (
	"context"
	"fmt"
	"giks/args"
	"giks/commands/hooks"
	"giks/config"
	"os"
)

func main() {
	// parse into specific giks arguments to ease command, subcommand and argument handling
	var ga args.GiksArgs = os.Args

	// enrich the context with the given config
	ctx := config.ContextWithConfig(context.Background(), ga.ConfigFile(), ga.GitDir())

	switch ga.Command() {
	case "hooks":
		hooks.ProcessHooks(ctx, ga)
	case "help":
		cfg := config.ConfigFromContext(ctx)
		fmt.Println("Help text")
		fmt.Printf("giks config: '%s' git directory: '%s'", cfg.ConfigFile, cfg.GitDir)
	default:
		fmt.Println("unknown command")
		os.Exit(1)
	}
}