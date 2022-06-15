package test

import "context"

// BeforeTester plugins are those that will be run before the
// test command, this interface is useful for later CLI extensions
// and addons.
type BeforeTester interface {
	// Name of the BeforeTester, useful for debugging, and logging.
	Name() string

	// BeforeTest receives the context, the working directory
	// and the command line arguments.
	BeforeTest(context.Context, string, []string) error
}
