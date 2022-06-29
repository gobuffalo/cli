package pop

import (
	"context"

	"github.com/gobuffalo/cli/cmd/cli/plugin"
)

var Command = &command{}

type command struct {
	subcommands []Subcommand
}

func (c command) Name() string {
	return "db"
}

func (c command) Aliases() []string {
	return []string{"pop", "database"}
}

func (c command) HelpText() string {
	return "A tasty treat for all your database needs"
}

func (c *command) LongHelpText() string {
	return ""
}

func (c *command) Receive(pls plugin.Plugins) {
	return
}

func (c *command) Main(ctx context.Context, pwd string, args []string) error {
	return nil
}
