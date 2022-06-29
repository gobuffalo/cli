package pop

import (
	"context"

	"github.com/gobuffalo/pop/v6"
)

var Fix = &fix{}

type fix struct{}

func (c fix) Name() string {
	return "fix"
}

func (c fix) HelpText() string {
	return "Brings pop, soda, and fizz files in line with the latest APIs"
}

func (c fix) Run(context.Context, []pop.ConnectionDetails, []string) error {
	return nil
}
