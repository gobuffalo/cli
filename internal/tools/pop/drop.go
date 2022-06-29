package pop

import (
	"context"

	"github.com/gobuffalo/pop/v6"
)

var Drop = &drop{}

type drop struct{}

func (c drop) Name() string {
	return "drop"
}

func (c drop) HelpText() string {
	return "Drops databases for you"
}

func (c drop) Run(context.Context, []pop.ConnectionDetails, []string) error {
	return nil
}
