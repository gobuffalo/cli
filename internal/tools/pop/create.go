package pop

import (
	"context"

	"github.com/gobuffalo/pop/v6"
)

var Create = &create{}

type create struct{}

func (c create) Name() string {
	return "create"
}

func (c create) HelpText() string {
	return "Creates databases for you"
}

func (c create) Run(context.Context, []pop.ConnectionDetails, []string) error {
	return nil
}
