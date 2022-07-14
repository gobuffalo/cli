package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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
}

// Main entry point for the application. This method finds the passed command
// and executes it with the passed arguments. If there is no command passed
// it will print the usage.
func (app *App) Main(ctx context.Context, pwd string, args []string) error {
	if app == nil {
		return fmt.Errorf("app is nil")
	}

	// Seek for custom command, which means the user is running
	// a codebase or user specific CLI.
	if cmd, p := app.CustomCommand(ctx, pwd, args); cmd != nil {
		fmt.Fprintf(app.Stdout(), "[Info] Running CLI in `%v`", p)

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

// CustomCommand returns an exec.Cmd if the user is overriding the CLI. This
// is used to allow users to add their own CLI plugins to the Buffalo CLI, when
// the CLI determines this is the case it runs the users CLI instead of the default
// CLI binary.
func (app *App) CustomCommand(ctx context.Context, pwd string, args []string) (*exec.Cmd, string) {
	// Here we take care of looking for CLI overriders
	// Overriders are go files to run instead of the CLI,
	// The two main use cases are:
	//  1. Running codebase specific CLI (PWD/cmd/buffalo/main.go) ✅
	//  2. Running user specific CLI ($HOME/buffalo/cmd/main.go)
	if _, err := os.Stat(filepath.Join(pwd, "cmd", "buffalo", "main.go")); err != nil {
		return nil, ""
	}

	cmd := exec.CommandContext(ctx, "go")
	cmd.Args = append(cmd.Args, "run", filepath.Join(pwd, "cmd", "buffalo", "main.go"))
	cmd.Args = append(cmd.Args, args[1:]...)
	cmd.Stdout = app.Stdout()
	cmd.Stderr = app.Stderr()
	cmd.Stdin = app.Stdin()

	return cmd, "cmd/buffalo"
}

// NewApp creates a CLI app with the given plugins.
// It prepends the `help` and `plugins` commands
// to the list of plugins.
func NewApp(plugins ...plugin.Plugin) *App {
	return &App{
		plugins: append(basePlugins, plugins...),
		help:    help.Command,
	}
}
