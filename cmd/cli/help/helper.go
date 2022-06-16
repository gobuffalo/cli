package help

import "context"

// Helper is a plugin that can provide help
// for given args, this is useful for subcommands
// and other plugins that want to provide help for
// its parts.
type Helper interface {
	// Help receives the args and prints the help
	// depending on plugin logic. Things like
	// looking for the subcommand details should happen
	// here.
	Help(context.Context, []string) error
}
