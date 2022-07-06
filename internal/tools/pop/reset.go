package pop

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/gobuffalo/pop/v6"
)

var Reset = &reset{}

type reset struct {
	flagSet *flag.FlagSet

	all   bool
	input string
	env   string
}

func (c reset) Usage() string {
	return "buffalo db reset [flags]"
}

func (c reset) Name() string {
	return "reset"
}

func (c *reset) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name(), flag.ExitOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(io.Discard)

		c.flagSet.BoolVar(&c.all, "all", false, "reset all databases")
		c.flagSet.StringVar(&c.env, "env", "development", "environment to be reset")
		c.flagSet.StringVar(&c.input, "input", "schema.sql", "The path to the schema file you want to load")
	}

	_ = c.flagSet.Parse(args)

	return c.flagSet, nil
}

func (c reset) HelpText() string {
	return "Drop, then recreate databases"
}

func (c *reset) PopMain(ctx context.Context, pwd string, args []string) error {
	// Fallback to migrations
	// if input cannot be opened.
	useMigrations := true
	var schema *os.File

	if _, err := os.Stat(c.input); err == nil {
		schema, err = os.Open(c.input)
		if err != nil {
			return err
		}

		useMigrations = false
		defer schema.Close()
	}

	conn := pop.Connections[c.env]
	if conn == nil {
		return fmt.Errorf("no connection named %s", c.env)
	}

	conns := []*pop.Connection{conn}
	if c.all {
		for _, v := range pop.Connections {
			conns = append(conns, v)
		}
	}

	for _, cx := range conns {
		if err := pop.DropDB(cx); err != nil {
			return err
		}

		if err := pop.CreateDB(cx); err != nil {
			return err
		}

		mig, err := pop.NewFileMigrator("migrations", cx)
		if err != nil {
			return err
		}

		// Apply the migrations directly
		if useMigrations {
			err = mig.Up()
			if err != nil {
				return err
			}

			continue
		}

		// Otherwise, use schema instead
		if err := cx.Dialect.LoadSchema(schema); err != nil {
			return err
		}

		// Then load migrations entries, without applying them
		err = mig.UpLogOnly()
		if err != nil {
			return err
		}
	}

	return nil
}
