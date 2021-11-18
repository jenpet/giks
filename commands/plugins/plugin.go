package plugins

import (
	"fmt"
)

// List holds all built-in giks plugins
var List = []Plugin{
	StringValidator{},
}

// Get returns a plugin for the given name (identifier). In case none was found an error is returned.
func Get(name string) (Plugin, error) {
	for _, p := range List {
		if p.ID() == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("plugin '%s' not found", name)
}

// Plugin has to be implemented by all built-in plugins in order to be triggered correctly by the giks hook executor
type Plugin interface {
	Run(hook string, vars map[string]string, args []string) (bool, error)
	ID() string
}
