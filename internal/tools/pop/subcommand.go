package pop

import (
	"context"

	"github.com/gobuffalo/cli/cmd/cli/help"
	"github.com/gobuffalo/cli/cmd/cli/plugin"
	"github.com/gobuffalo/pop/v6"
)

// Subcommand is the interface of the plugins that could be hooked into
// the pop command.
type Subcommand interface {
	plugin.Plugin
	help.HelpTexter

	Run(context.Context, *pop.Connection) error
}

// Subcommands is a convenient type for a list
// of Subcommand.
type Subcommands []Subcommand

// Find a command from the list given his name
// or aliases if the command is Aliaser.
func (cc Subcommands) Find(name string) Subcommand {
	for _, v := range cc {
		if v.Name() == name {
			return v
		}

		al, ok := v.(plugin.Aliaser)
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
