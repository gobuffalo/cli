package pop

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

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

func (c command) Usage() string {
	return "buffalo db <subcommand> [flags] [options]"
}

func (c command) HelpText() string {
	return "A tasty treat for all your database needs"
}

func (c *command) LongHelpText() string {
	buf := bytes.NewBuffer([]byte{})
	w := tabwriter.NewWriter(buf, 0, 0, 3, ' ', 0)

	w.Write([]byte("Subcommands\n"))
	for _, gg := range c.subcommands {
		fmt.Fprintf(w, "  %s\t%s\n", gg.Name(), gg.HelpText())
	}

	w.Flush()

	return buf.String()
}

func (c *command) Receive(pls plugin.Plugins) {
	for _, v := range pls {
		if sc, ok := v.(Subcommand); ok {
			c.subcommands = append(c.subcommands, sc)
		}
	}
}

func (c *command) Main(ctx context.Context, pwd string, args []string) error {
	return nil
}
