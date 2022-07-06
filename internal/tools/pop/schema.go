package pop

import (
	"context"
)

var Schema = &schema{}

type schema struct{}

func (c schema) Name() string {
	return "schema"
}

func (c schema) HelpText() string {
	return "Tools for working with your database schema"
}

func (c schema) PopMain(ctx context.Context, pwd string, args []string) error {
	return nil
}
