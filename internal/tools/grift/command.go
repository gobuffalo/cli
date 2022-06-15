package grift

import (
	"context"

	grifts "github.com/markbates/grift/cmd"
)

// Command is the shared instance of the
// Grift command
var Command = command("Grift")

type command string

func (tc command) Aliases() []string {
	return []string{"t", "tasks"}
}

func (tc command) Name() string {
	return "task"
}

func (tc command) HelpText() string {
	return "Runs grift tasks with the passed arguments."
}

func (tc command) Main(ctx context.Context, pwd string, args []string) error {
	return grifts.Run("buffalo task", args)
}
