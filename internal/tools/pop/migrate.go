package pop

import (
	"context"

	"github.com/gobuffalo/pop/v6"
)

var Migrate = &migrate{}

type migrate struct{}

func (c migrate) Name() string {
	return "migrate"
}

func (c migrate) HelpText() string {
	return "Runs migrations against your database."
}

func (c migrate) Run(context.Context, []pop.ConnectionDetails, []string) error {
	return nil
}
