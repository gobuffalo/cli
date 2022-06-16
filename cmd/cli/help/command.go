package help

import (
	"context"
	"flag"
	"fmt"

	"github.com/gobuffalo/cli/cmd/cli/clio"
	"github.com/gobuffalo/cli/cmd/cli/plugin"
)

// Command is in charge of printing the help text for a given command.
// its flags and any other information available to make it easy for the user.
type Command struct {
	clio.IO

	Commands plugin.Commands
}

func (c Command) Name() string {
	return "help"
}

func (c Command) HelpText() string {
	return "Provides help for a given command, p.e. buffalo help list."
}

func (c Command) Main(ctx context.Context, pwd string, args []string) error {
	var command plugin.Command
	if len(args) > 0 {
		command = c.Commands.Find(args[0])
		if command == nil {
			fmt.Fprintf(c.Stdout(), "Error: did not find `%v` command\n", args[0])
		}
	}

	if len(args) == 0 || command == nil {
		return c.General()
	}

	return c.Specific(command)
}

// ReceivePlugins and keep the commands for the help
// command to use.
func (c *Command) Receive(plugins plugin.Plugins) {
	c.Commands = plugin.CommandsFrom(plugins)
}

// General method prints help text for all of the commands.
func (c Command) General() error {
	fmt.Fprint(c.Stdout(), "Usage: buffalo [command] [options]\n\n")

	// If there are no commands it just prints the usage.
	if len(c.Commands) == 0 {
		return nil
	}

	fmt.Fprintln(c.Stdout(), "Available Commands:")
	for _, v := range c.Commands {
		if ht, ok := v.(HelpTexter); ok {
			fmt.Fprintf(c.Stdout(), "%v\t\t%v\n", v.Name(), ht.HelpText())

			continue
		}

		fmt.Fprintf(c.Stdout(), "%v\t (runs the %[1]v command)\n", v.Name())
	}

	fmt.Fprintln(c.Stdout(), "\nFor command specific information use the help command, p.e.")
	fmt.Fprintln(c.Stdout(), "$ buffalo help [command]")

	return nil
}

// Specific help text for a command passed.
// if the command is a FlagParser it uses the flagSet to
// print help text for the flags. Also if the command implements
// LongHelpTexter it prints the long help text.
func (c Command) Specific(cm plugin.Command) error {
	fmt.Fprintf(c.Stdout(), "Usage: buffalo %v [options]\n\n", cm.Name())

	if ht, ok := cm.(HelpTexter); ok {
		fmt.Fprintf(c.Stdout(), ht.HelpText()+"\n\n")
	}

	if ht, ok := cm.(LongHelpTexter); ok {
		fmt.Fprintf(c.Stdout(), ht.LongHelpText())
	}

	if fl, ok := cm.(clio.FlagParser); ok {
		fl, _ := fl.ParseFlags([]string{})

		fmt.Fprintf(c.Stdout(), "Flags:\n")
		fl.VisitAll(func(ff *flag.Flag) {
			fmt.Fprintf(c.Stdout(), "--%v\t%v\n", ff.Name, ff.Usage)
		})
	}

	return nil
}
