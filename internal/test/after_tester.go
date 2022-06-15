package test

import "context"

type AfterTester interface {
	// Name of the AfterTester, useful for debugging, and logging.
	Name() string

	// The actual function that will run AfterTest, returns an error.
	AfterTest(context.Context, string, []string) error
}
