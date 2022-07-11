package dev

import (
	"context"
	"os"
)

var SetupDevelopment = setupEnv("setup")

type setupEnv string

func (se setupEnv) Name() string {
	return "dev/setup-env"
}

func (se setupEnv) HelpText() string {
	return "Sets the GO_ENV to development before running the app with the `dev` command."
}

func (se setupEnv) SetupDevelopment(ctx context.Context, pwd string, args []string) error {
	return os.Setenv("GO_ENV", "development")
}
