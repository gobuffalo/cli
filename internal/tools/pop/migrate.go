package pop

import (
	"context"
)

var Migrate = &migrate{}

type migrate struct{}

func (c migrate) Name() string {
	return "migrate"
}

func (c migrate) HelpText() string {
	return "Runs migrations against your database."
}

func (c migrate) PopMain(ctx context.Context, pwd string, args []string) error {
	return nil
}
