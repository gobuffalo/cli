package destroy

import (
	"context"

	"github.com/gobuffalo/cli/cmd/cli/plugin"
)

// Destroyer is a plugin that will destroy existing files
// by using passed args. Destroyers are not flag parsers
// and to keep it simple will only be marked as Preconfirmed
// with the PreConfirm method.
type Destroyer interface {
	plugin.Plugin

	// Sets destroyer not to confirm dangerous actions.
	PreConfirm()

	// Destroys generated files depending on the type of
	// destroyer.
	Destroy(context.Context, string, []string) error
}

// Generators is a convenient type for a list
// of Generator.
type Destroyers []Destroyer

// Find a command from the list given his name
// or aliases if the command is Aliaser.
func (cc Destroyers) Find(name string) Destroyer {
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
