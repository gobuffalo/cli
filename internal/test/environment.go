package test

import (
	"context"
	"os"
)

// TestEnvironment Plugin is a BeforeTester that sets
// the GO_ENV environment to be test.
const TestEnvironment = testEnvironment("test")

type testEnvironment string

func (e testEnvironment) Name() string {
	return "test_environment"
}

func (e testEnvironment) HelpText() string {
	return "Sets the GO_ENV environment to be 'test'."
}

func (e testEnvironment) BeforeTest(ctx context.Context, pwd string, args []string) error {
	return os.Setenv("GO_ENV", "test")
}
