package cli

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/gobuffalo/cli/cmd/cli/clio"
	"github.com/gobuffalo/cli/cmd/cli/help"
	"github.com/gobuffalo/cli/cmd/cli/plugin"

	"github.com/gobuffalo/cli/internal/build"
	"github.com/gobuffalo/cli/internal/destroy"
	"github.com/gobuffalo/cli/internal/dev"
	"github.com/gobuffalo/cli/internal/fix"
	"github.com/gobuffalo/cli/internal/generate"
	"github.com/gobuffalo/cli/internal/info"
	"github.com/gobuffalo/cli/internal/new"
	"github.com/gobuffalo/cli/internal/routes"
	"github.com/gobuffalo/cli/internal/setup"
	"github.com/gobuffalo/cli/internal/test"
	"github.com/gobuffalo/cli/internal/tools/bzr"
	"github.com/gobuffalo/cli/internal/tools/frontend"
	"github.com/gobuffalo/cli/internal/tools/git"
	"github.com/gobuffalo/cli/internal/tools/grift"
	"github.com/gobuffalo/cli/internal/tools/pop"
	"github.com/gobuffalo/cli/internal/version"
)

var (
	basePlugins = plugin.Plugins{
		help.Command,
		plugin.List,
	}

	defaultPlugins = plugin.Plugins{
		// Top level commands
		test.Command,
		version.Command,
		grift.Command,
		routes.Command,
		setup.Command, // TODO: DOCS: Document how to wire setuppers here.
		info.Command,
		generate.Command, // TODO: DOCS: Document how to add generators
		fix.Command,
		new.Command,
		destroy.Command,
		build.Command,
		dev.Command, // TODO: DOCS: Document how to add dev plugins
		pop.Command, // TODO: This needs to live outside of the CLI package and into the pop/buffalo-pop package.

		pop.Create,
		pop.Drop,
		pop.Fix,
		pop.Migrate,
		pop.Reset,
		pop.Generate,
		pop.Schema,

		// Generators
		generate.ActionGenerator,
		generate.MailerGenerator,
		generate.ResourceGenerator,
		grift.Generator,

		//Destroyers
		destroy.ActionDestroyer,
		destroy.MailerDestroyer,
		destroy.ResourceDestroyer,

		// Setup plugins
		grift.SetupSeedDatabase,
		pop.Setup,
		frontend.Setup,
		setup.Test,

		// Development plugins
		dev.SetupDevelopment,
		dev.StartServer,
		dev.StartFrontend,

		git.VersionRunner,
		bzr.VersionRunner,
	}
)

type App struct {
	clio.IO

	help    help.GeneralHelper
	plugins plugin.Plugins

	overriders []overrider
}

// overrider allows to override the default plugins by running
// custom main.go.
type overrider func(pwd string, args []string) (*exec.Cmd, string)

// Main entry point for the application. This method finds the passed command
// and executes it with the passed arguments. If there is no command passed
// it will print the usage.
func (app *App) Main(ctx context.Context, pwd string, args []string) error {
	if app == nil {
		return fmt.Errorf("app is nil")
	}

	// Seek for an overrider that provides a command to execute instead
	// of the default flow.
	for _, v := range app.overriders {
		cmd, p := v(pwd, args)
		if cmd == nil {
			continue
		}

		fmt.Fprintf(app.Stdout(), "[Info] Running CLI in `%v`\n\n", p)
		return cmd.Run()
	}

	// Pass all of the plugins to the PluginsReceivers in the
	// list of plugins so that they can keep copy of these.
	for _, v := range app.plugins {
		pr, ok := v.(plugin.Receiver)
		if !ok {
			continue
		}

		pr.Receive(app.plugins)
	}

	if len(args) == 0 {
		return app.help.General()
	}

	// Find the command from the list of commands
	// to determine what to show to the user.
	command := plugin.CommandsFrom(app.plugins).Find(args[0])
	if command == nil {
		// Print out general help if no command is passed.
		return app.help.General()
	}

	args = args[1:]
	if fp, ok := command.(clio.FlagParser); ok {
		// Pass the args to the command, it should take care  of passing
		// the args to subcommands in case it applies so that
		// these are prepared to be executed.
		_, err := fp.ParseFlags(args)
		if err != nil {
			return err
		}
	}

	if ist, ok := command.(clio.Setter); ok {
		ist.SetIO(app.Stdout(), app.Stderr(), app.Stdin())
	}

	if wdv, ok := command.(plugin.WorkDirValidator); ok {
		valid, err := wdv.ValidateWorkDir(pwd)
		if err != nil {
			return err
		}

		if !valid {
			return fmt.Errorf("'%v' command cannot run in '%v'", command.Name(), pwd)
		}
	}

	return command.Main(ctx, pwd, args)
}

// NewApp creates a CLI app with the given plugins.
// It prepends the `help` and `plugins` commands
// to the list of plugins. The new app will not
// have the
func NewApp(plugins ...plugin.Plugin) *App {
	return &App{
		plugins: append(basePlugins, plugins...),
		help:    help.Command,
	}
}
