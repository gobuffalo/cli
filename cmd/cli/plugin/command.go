package plugin

import (
	"context"
)

// Command is a struct that can be invoked by the main app
// to do so it would use the Name method to identify it.
type Command interface {
	Plugin

	// Main will be called to execute the command
	// and passed with the context, working directory and the args.
	Main(ctx context.Context, pwd string, args []string) error
}

// Commands holds a set of useful methods for working with
// a group of commands.
type Commands []Command

func (cc Commands) Contains(name string) bool {
	return cc.Find(name) != nil
}

// Find a command from the list given his name
// or aliases if the command is Aliaser.
func (cc Commands) Find(name string) Command {
	for _, v := range cc {
		if v.Name() == name {
			return v
		}

		al, ok := v.(Aliaser)
		if !ok {
			continue
		}

		// If the command is an alias we need to check if one the
		// alias is the one we are looking for.
		for _, a := range al.Aliases() {
			if a == name {
				return v
			}
		}
	}

	return nil
}

// Commands returns the list of commands within
// a list of plugins.
func CommandsFrom(pls Plugins) Commands {
	cmds := Commands{}
	for _, v := range pls {
		c, ok := v.(Command)
		if !ok {
			continue
		}

		cmds = append(cmds, c)
	}

	return cmds
}
