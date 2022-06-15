package cli

import (
	"context"
	"fmt"

	"github.com/gobuffalo/cli/cmd/cli/clio"
	"github.com/gobuffalo/cli/cmd/cli/help"
	"github.com/gobuffalo/cli/cmd/cli/plugin"
)

type App struct {
	clio.Container

	help    *help.Command
	plugins plugin.Plugins
}

// Main entry point for the application. This method finds the passed command
// and executes it with the passed arguments. If there is no command passed
// it will print the usage.
func (app *App) Main(ctx context.Context, pwd string, args []string) error {
	if app == nil {
		return fmt.Errorf("app is nil")
	}

	if len(args) == 0 {
		return app.help.General()
	}

	command := plugin.CommandsFrom(app.plugins).Find(args[0])
	if command == nil {
		// Print out general help if no command is passed.
		return app.help.General()
	}

	args = args[1:]
	if fp, ok := command.(clio.FlagParser); ok {
		fs, err := fp.ParseFlags(args)
		if err != nil {
			return err
		}

		// We update the args to remove the parsed flags.
		args = fs.Args()
	}

	if ist, ok := command.(clio.Setter); ok {
		ist.SetIO(app.Stdout(), app.Stderr(), app.Stdin())
	}

	return command.Main(ctx, pwd, args)
}

// NewApp creates a CLI app with the given plugins.
// It prepends the `help` command to the list of plugins.
func NewApp(plugins ...plugin.Plugin) *App {
	// Initializing the Help command and prepending it to
	// the list of plugins passed.
	help := &help.Command{}
	plugins = append(plugin.Plugins{help}, plugins...)

	// Pass all of the plugins to the PluginsReceivers in the
	// list of plugins.
	for _, v := range plugins {
		pr, ok := v.(plugin.Receiver)
		if !ok {
			continue
		}

		pr.Receive(plugins)
	}

	return &App{
		plugins: plugins,
		help:    help,
	}
}
