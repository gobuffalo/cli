package pop

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
	return nil
}
