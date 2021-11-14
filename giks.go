package main

import (
	"context"
	"flag"
	"giks/args"
	"giks/commands"
	"giks/config"
	"os"
)

func main() {
	var ga args.GiksArgs = os.Args

	ctx := config.ContextWithConfig(context.Background())

	switch ga.Command() {
	case "hooks":
		commands.ProcessHooks(ctx, ga)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}