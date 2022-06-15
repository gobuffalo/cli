package test

import (
	"context"
	"fmt"

	"github.com/gobuffalo/cli/internal/tools/pop"
)

// The test command instance with default before and
// after test plugins.
var Command = &command{
	before: []BeforeTester{
		TestEnvironment,
		&pop.SetupDatabase{},
	},
}

type command struct {
	before  []BeforeTester
	testers []Tester
}

func (c command) Name() string {
	return "test"
}

func (c command) HelpText() string {
	// TODO: List before and after plugins
	return "Runs application tests by invoking before and after test plugins."
}

func (c *command) Main(ctx context.Context, pwd string, args []string) error {
	// Iterate over the BeforeTesters and run each of them
	// in case of an error halt the testing process by returning the error.
	for _, v := range c.before {
		err := v.BeforeTest(ctx, pwd, args)
		if err == nil {
			continue
		}

		return fmt.Errorf("error running `%s` before test: %w", v.Name(), err)
	}

	return nil
}
