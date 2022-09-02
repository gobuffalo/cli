package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/gobuffalo/cli/cmd/cli/clio"
	"github.com/gobuffalo/cli/cmd/cli/help"
	"github.com/gobuffalo/cli/cmd/cli/plugin"
)

// overrider allows to override the default plugins by running
// custom main.go.
type overrider func(pwd string) (*exec.Cmd, string)

type app struct {
	clio.IO

	help    help.GeneralHelper
	plugins plugin.Plugins

	overriders []overrider
}

// Main entry point for the application. This method finds the passed command
// and executes it with the passed arguments. If there is no command passed
// it will print the usage.
func (a *app) Main(ctx context.Context, pwd string, args []string) error {
	// Setting the IO for the help command in case its clio.Setter.
	if ist, ok := a.help.(clio.Setter); ok {
		ist.SetIO(a.Stdout(), a.Stderr(), a.Stdin())
	}

	// Seek for an overrider that provides a command to execute instead
	// of the default flow.
	for _, v := range a.overriders {
		cmd, p := v(pwd)
		if cmd == nil {
			continue
		}

		fmt.Fprintf(a.Stdout(), "[Info] Running CLI in `%v`\n\n", p)
		cmd.Stdout = a.Stdout()
		cmd.Stderr = a.Stderr()
		cmd.Stdin = a.Stdin()

		return cmd.Run()
	}

	// Pass all of the plugins to the PluginsReceivers in the
	// list of plugins so that they can keep copy of these.
	for _, v := range a.plugins {
		pr, ok := v.(plugin.Receiver)
		if !ok {
			continue
		}

		pr.Receive(a.plugins)
	}

	if len(args) == 0 {
		return a.help.General()
	}

	// Find the command from the list of commands
	// to determine what to show to the user.
	command := plugin.CommandsFrom(a.plugins).Find(args[0])
	if command == nil {
		// Print out general help if no command is passed.
		return a.help.General()
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
		ist.SetIO(a.Stdout(), a.Stderr(), a.Stdin())
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

// Run starts the CLI by tracking the PWD, creating a context
// and running the Main method.
func (a *app) Run() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	// get the present working directory. (PWD)
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	err = a.Main(ctx, pwd, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		os.Exit(1)
	}

	<-ctx.Done()
}

// NewApp creates a CLI app with the given plugins.
// It prepends the `help` and `plugins` commands
// to the list of plugins. The new app will not
// have the
func NewApp(plugins ...plugin.Plugin) *app {
	return &app{
		plugins: append(basePlugins, plugins...),
		help:    help.Command,
	}
}

// NewWithDefaults creates a new CLI app with the
// default plugins and adds the extra plugins.
func NewWithDefaults(extra ...plugin.Plugin) *app {
	initial := append(basePlugins, defaultPlugins...)
	return &app{
		plugins: append(initial, extra...),
		help:    help.Command,
	}
}

// NewWithOverriders returns a pointer to a CLI that
// will override the CLI with project and user overriders
// this is for the default Buffalo CLI and not intended to be
// used in custom CLI's.
func NewWithOverriders() *app {
	return &app{
		help:    help.Command,
		plugins: append(basePlugins, defaultPlugins...),

		overriders: []overrider{
			projectOverrider,
			userOverrider,
		},
	}
}
