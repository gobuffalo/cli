package pop

import (
	"context"

	"github.com/gobuffalo/pop/v6"
)

var Schema = &schema{}

type schema struct{}

func (c schema) Name() string {
	return "schema"
}

func (c schema) HelpText() string {
	return "Tools for working with your database schema"
}

func (c schema) Run(ctx context.Context, conn *pop.Connection) error {
	return nil
}
