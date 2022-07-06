package pop

import (
	"context"

	"github.com/gobuffalo/cli/internal/tools/pop/generators"
	"github.com/gobuffalo/pop/v6"
)

var Generate = &generate{
	generators: Generators{
		generators.Config,
		generators.Fizz,
		generators.SQL,
		generators.Model,
	},
}

type generate struct {
	generators Generators
}

func (g generate) Name() string {
	return "generate"
}

func (g generate) Aliases() []string {
	return []string{"g"}
}

func (g generate) HelpText() string {
	return "Generates config, model, and migrations files."
}

func (g generate) Run(context.Context, *pop.Connection) error {
	return nil
}
