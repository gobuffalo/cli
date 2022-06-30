package pop

import (
	"context"

	"github.com/gobuffalo/pop/v6"
)

var Generate = &generate{}

type generate struct{}

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
