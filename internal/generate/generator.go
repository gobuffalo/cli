package generate

import (
	"context"

	"github.com/gobuffalo/cli/cmd/cli/plugin"
)

// Generator is a plugin that generates code block
type Generator interface {
	plugin.Plugin

	// HelpText is a short description of the
	// generator.
	HelpText() string

	// Generate the desired code block or return
	// and error if something goes wrong.
	Generate(context.Context, string, []string) error
}

// Generators is a convenient type for a list
// of Generator.
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
