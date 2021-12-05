package commands

import (
	"fmt"
	gargs "github.com/jenpet/giks/args"
	"github.com/jenpet/giks/cli"
	"github.com/jenpet/giks/config"
	"github.com/jenpet/giks/log"
)

func Process(cfg config.Config, gargs gargs.GiksArgs) {
	// actual array of arguments without the binary itself the command and subcommand
	args := gargs.Args()
	switch gargs.Command() {
	case "install":
		if gargs.HasHook() {
			h := cfg.Hook(gargs.Hook())
			installSingleHook(cfg, h, true)
			break
		}
		installHookList(cfg)
	case "uninstall":
		if gargs.HasHook() {
			h := cfg.Hook(gargs.Hook())
			uninstallSingleHook(cfg, h, true)
			break
		}
		uninstallHookList(cfg)
	case "exec":
		if err := executeHook(cfg, gargs); err != nil {
			log.Errorf("failed executing '%s' hook. Error: %s", gargs.Hook(), err)
		}
	case "show":
		cli.PrintTemplate(detailsTemplate, cfg.Hook(gargs.Hook()).ToMap())
	case "list":
		_ = listCommand.Parse(args)
		if listCommand.Parsed() {
			all := *listAllAttr
			cli.PrintTemplate(listTemplate, cfg.HookList(all))
		}
	case "help":
		printHelp(cfg, gargs)
	default:
		fmt.Printf("Unknown command '%s'", gargs.Command())
	}
}
