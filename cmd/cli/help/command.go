package help

import (
	"context"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/gobuffalo/cli/cmd/cli/clio"
	"github.com/gobuffalo/cli/cmd/cli/plugin"
)

// Command to be used outside of the Help package.
var Command = &command{
	IO: &clio.IO{},
}

// command is in charge of printing the help text for a given command.
// its flags and any other information available to make it easy for the user.
type command struct {
	*clio.IO

	commands plugin.Commands
}

func (c command) Name() string {
	return "help"
}

func (c command) HelpText() string {
	return "Provides help for a given command, p.e. buffalo help list."
}

func (c command) Main(ctx context.Context, pwd string, args []string) error {
	var command plugin.Command
	if len(args) > 0 {
		command = c.commands.Find(args[0])
		if command == nil {
			fmt.Fprintf(c.Stdout(), "Error: did not find `%v` command\n", args[0])
		}
	}

	if len(args) == 0 || command == nil {
		return c.General()
	}

	hh, ok := command.(Helper)
	if !ok || len(args) == 1 {
		return c.Specific(command)
	}

	// If the command implements Helper
	// the command itself will take care
	// of printing the help with the args.
	return hh.Help(ctx, args[1:])
}

func (c *command) SetIO(stdout io.Writer, stderr io.Writer, stdin io.Reader) {
	c.IO.Out = stdout
	c.IO.Err = stderr
	c.IO.In = stdin
}

// ReceivePlugins and keep the commands for the help
// command to use.
func (c *command) Receive(plugins plugin.Plugins) {
	c.commands = plugin.CommandsFrom(plugins)
}

// General method prints help text for all of the commands.
func (c command) General() error {
	fmt.Fprint(c.Stdout(), "Usage: buffalo [command] [flags] [...]\n\n")

	// If there are no commands it just prints the usage.
	if len(c.commands) == 0 {
		return nil
	}

	fmt.Fprintln(c.Stdout(), "Available Commands:")
	w := tabwriter.NewWriter(c.Stdout(), 0, 0, 3, ' ', 0)

	for _, v := range c.commands {
		if ht, ok := v.(HelpTexter); ok {
			fmt.Fprintf(w, "%v\t\t%v\n", v.Name(), ht.HelpText())

			continue
		}

		fmt.Fprintf(w, "%v\t (runs the %[1]v command)\n", v.Name())
	}

	w.Flush()

	fmt.Fprintln(c.Stdout(), "\nFor command specific information use the help command, p.e.")
	fmt.Fprintln(c.Stdout(), "$ buffalo help [command]")

	return nil
}

// Specific help text for a command passed.
// if the command is a FlagParser it uses the flagSet to
// print help text for the flags. Also if the command implements
// LongHelpTexter it prints the long help text.
func (c command) Specific(cm plugin.Command) error {
	return Specific(c.IO.Stdout(), cm)
}
