package pop

import (
	"context"

	"github.com/gobuffalo/cli/cmd/cli/help"
	"github.com/gobuffalo/cli/cmd/cli/plugin"
)

type Generator interface {
	plugin.Plugin
	help.HelpTexter

	PopGenerate(context.Context, string, []string) error
}

type Generators []Generator

// Find a command from the list given his name
// or aliases if the command is Aliaser.
func (cc Generators) Find(name string) Generator {
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
