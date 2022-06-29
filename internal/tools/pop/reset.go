package pop

import (
	"context"

	"github.com/gobuffalo/pop/v6"
)

var Reset = &reset{}

type reset struct{}

func (c reset) Name() string {
	return "reset"
}

func (c reset) HelpText() string {
	return "Drop, then recreate databases"
}

func (c reset) Run(context.Context, []pop.ConnectionDetails, []string) error {
	return nil
}
