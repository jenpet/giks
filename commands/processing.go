package commands

import (
	"flag"
	"fmt"
	gargs "github.com/jenpet/giks/args"
	"github.com/jenpet/giks/cli"
	"github.com/jenpet/giks/config"
	"github.com/jenpet/giks/log"
	"github.com/jenpet/giks/meta"
)

var showCommand = flag.NewFlagSet("show", flag.ExitOnError)
var showAllAttr = showCommand.Bool("all", false, "include disabled hooks")

func Process(cfg config.Config, gargs gargs.GiksArgs) {
	// actual array of arguments without the binary itself the command and subcommand
	args := gargs.Args(true)
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
		if gargs.HasHook() {
			cli.PrintTemplate(detailsTemplate, cfg.Hook(gargs.Hook()).ToMap())
			break
		}
		_ = showCommand.Parse(args)
		if showCommand.Parsed() {
			all := *showAllAttr
			cli.PrintTemplate(listTemplate, cfg.HookList(all))
		}
	case "help":
		printHelp(cfg, gargs)
	case "version":
		fmt.Printf("Version: %s (#%s)", meta.Version(), meta.CommitHash())
	default:
		fmt.Printf("Unknown command. Use `giks help` for more information.")
	}
}
