package pop

<<<<<<< HEAD
import "context"

var Command = &command{}

type command struct{}

func (c command) Name() string {
	return "pop"
}

func (c command) Aliases() []string {
	return []string{"db"}
}

func (c command) Usage() string {
	return "buffalo pop [subcommand] [flags] [...]"
}

func (c command) HelpText() string {
	return "Manage your database with Pop"
}

// func (c command) LongHelpText() string {
// 	return "Manage your database with Pop [TODO: Subcommands]"
// }

func (c command) Main(ctx context.Context, pwd string, args []string) error {
=======
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
>>>>>>> bdd723c64171956e9ae1272aaddbc7c1d3f9bdda
	return nil
}
