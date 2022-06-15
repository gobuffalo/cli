package cli

import (
	"context"
	"fmt"
)

type App struct {
	IO

	commands Commands
	usage    func() error
}

// Main entry point for the application. This method finds the passed command
// and executes it with the passed arguments. If there is no command passed
// it will print the usage.
func (app *App) Main(ctx context.Context, pwd string, args []string) error {
	if app == nil {
		return fmt.Errorf("app is nil")
	}

	if len(args) == 0 {
		return app.usage()
	}

	command := app.commands.Find(args[0])
	if command == nil {
		return app.usage()
	}

	args = args[1:]
	if fp, ok := command.(FlagParser); ok {
		fs, err := fp.ParseFlags(args)
		if err != nil {
			return err
		}

		// We update the args to remove the parsed flags.
		args = fs.Args()
	}

	if ist, ok := command.(IOSetter); ok {
		ist.SetIO(app.Stdout(), app.Stderr(), app.Stdin())
	}

	return command.Main(ctx, pwd, args)
}

// NewApp creates a CLI app with the given commands.
// It prepends the `help` command to the list of commands.
func NewApp(commands ...Command) *App {
	// An instance of the help command, to be able to reference
	// its general function.
	help := &HelpCommand{}

	// Adding the help command always.
	cmms := append(Commands{help}, commands...)
	help.Commands = cmms

	return &App{
		commands: cmms,
		usage:    help.general,
	}
}
