package config

import (
	"context"
	"fmt"
	"os"
)

// ContextWithConfig adds the configuration to the context to provide easy access later on
func ContextWithConfig(ctx context.Context, file string, gitDir string) context.Context {
	cfg := assembleConfig(file, gitDir)
	return context.WithValue(ctx, ctxConfigKey, cfg)
}

// ContextWithHook adds the targeted hook into the context to provide easy access later on
func ContextWithHook(ctx context.Context, args []string) context.Context {
	if len(args) < 1 {
		fmt.Println("hook name is required but missing")
		os.Exit(1)
	}
	cfg := ConfigFromContext(ctx)
	hook, err := cfg.LookupHook(args[0])
	if err != nil {
		fmt.Printf("failed retrieving '%s' hook. Error: %s\n", args[0], err)
		os.Exit(1)
	}
	return context.WithValue(ctx, ctxHookKey, *hook)
}

func HookFromContext(ctx context.Context) Hook {
	if ctx == nil {
		fmt.Println("could not retrieve hook from context context was nil.")
		os.Exit(1)
	}
	if h, ok := ctx.Value(ctxHookKey).(Hook); ok {
		return h
	}
	fmt.Println("could not retrieve hook from context")
	os.Exit(1)
	return Hook{}
}

// ConfigFromContext returns the configuration from a given context
func ConfigFromContext(ctx context.Context) Config {
	if ctx == nil {
		fmt.Println("could not retrieve config from context context was nil.")
		os.Exit(1)
	}
	if cfg, ok := ctx.Value(ctxConfigKey).(Config); ok {
		return cfg
	}
	fmt.Println("could not retrieve config from context")
	os.Exit(1)
	return Config{}
}

