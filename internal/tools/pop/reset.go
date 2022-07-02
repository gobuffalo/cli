package pop

import (
	"context"
	"flag"
	"io"
	"os"

	"github.com/gobuffalo/pop/v6"
)

var Reset = &reset{}

type reset struct {
	flagSet *flag.FlagSet

	all   bool
	input string
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
		c.flagSet.StringVar(&c.input, "input", "schema.sql", "The path to the schema file you want to load")
	}

	_ = c.flagSet.Parse(args)

	return c.flagSet, nil
}

func (c reset) HelpText() string {
	return "Drop, then recreate databases"
}

func (c *reset) Run(ctx context.Context, conn *pop.Connection) error {
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
