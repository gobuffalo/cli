package routes

import (
	"context"

	"github.com/gobuffalo/cli/cmd/cli/plugin"
	grifts "github.com/markbates/grift/cmd"
)

const Command = command("routes")

type command string

func (c command) Name() string {
	return "routes"
}

func (c command) HelpText() string {
	return "Prints a list of all defined routes"
}

func (c command) ValidateWorkDir(wd string) (bool, error) {
	return plugin.ValidateBuffaloRoot(wd)
}

func (c command) Main(ctx context.Context, pwd string, args []string) error {
	return grifts.Run("buffalo task", []string{"routes"})
}
