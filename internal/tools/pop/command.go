package pop

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"text/tabwriter"

	"github.com/gobuffalo/cli/cmd/cli/clio"
	"github.com/gobuffalo/cli/cmd/cli/plugin"
	"github.com/gobuffalo/pop/v6"
)

var Command = &command{}

type command struct {
	flagSet *flag.FlagSet

	env    string
	config string

	subcommands Subcommands
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

func (c *command) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name(), flag.ExitOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(io.Discard)

		c.flagSet.StringVar(&c.env, "env", "development", "the environment to use")
		c.flagSet.StringVar(&c.config, "config", "database.yml", "the path to the config file")
	}

	_ = c.flagSet.Parse(args)

	// Takes care of calling its subcommands and
	// passing args.
	for _, v := range c.subcommands {
		fp, ok := v.(clio.FlagParser)
		if !ok {
			continue
		}

		if len(args) == 0 {
			continue
		}

		// Remove the first argument
		// as it is the name of the subcommand
		ax := args[1:]
		fp.ParseFlags(ax)
	}

	return c.flagSet, nil
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
	if len(args) == 0 {
		return fmt.Errorf("please specify the subcommand to use")
	}

	// TODO: Only do when needed
	if c.config != "" {
		abs, err := filepath.Abs(c.config)
		if err != nil {
			return err
		}

		dir, file := filepath.Split(abs)

		pop.AddLookupPaths(dir)
		pop.ConfigName = file
	}

	pop.LoadConfigFile()

	conn := pop.Connections[c.env]
	if conn == nil {
		return fmt.Errorf("There is no connection named '%s' defined!\n", c.env)
	}

	cx := c.subcommands.Find(args[0])
	if cx == nil {
		return fmt.Errorf("no subcommand found for '%v'", args[0])
	}

	return cx.Run(ctx, conn)
}
