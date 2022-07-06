package pop

import (
	"context"
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/gobuffalo/pop/v6"
)

var Schema = &schema{}

type schema struct {
	flagSet *flag.FlagSet

	env    string
	output string
}

func (c schema) Name() string {
	return "schema"
}

func (c schema) Usage() string {
	return "buffalo db schema [dump|load]"
}

func (c schema) HelpText() string {
	return "Tools for working with your database schema"
}

func (c *schema) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name(), flag.ExitOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(io.Discard)
	}

	c.flagSet.StringVar(&c.env, "env", "development", "environment to be reset")
	c.flagSet.StringVar(&c.output, "input", "./migrations/schema.sql", "The path to the schema file you want to load")

	_ = c.flagSet.Parse(args)

	return c.flagSet, nil
}

func (c schema) PopMain(ctx context.Context, pwd string, args []string) error {
	conn, err := pop.Connect(c.env)
	if err != nil {
		return err
	}

	out := os.Stdout
	rollback := func() {}

	if c.output != "-" {
		err = os.MkdirAll(filepath.Dir(c.output), 0755)
		if err != nil {
			return err
		}

		out, err = os.Create(c.output)
		if err != nil {
			return err
		}

		rollback = func() {
			os.RemoveAll(c.output)
		}
	}

	if err := conn.Dialect.DumpSchema(out); err != nil {
		rollback()

		return err
	}

	return nil
}
