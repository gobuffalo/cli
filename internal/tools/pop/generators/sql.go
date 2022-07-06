package generators

import (
	"context"
	"errors"
	"flag"
	"io"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/pop/v6/genny/fizz/cempty"
	"github.com/gobuffalo/pop/v6/genny/fizz/ctable"
)

var SQL = &sqlGenerator{}

type sqlGenerator struct {
	flagSet *flag.FlagSet

	path string
	env  string
}

func (c sqlGenerator) Name() string {
	return "sql"
}

func (c sqlGenerator) Usage() string {
	return "buffalo db generate sql <name> [attributes]"
}

func (c sqlGenerator) HelpText() string {
	return "Generates Up/Down migrations for your database using sql."
}

func (c *sqlGenerator) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(io.Discard)
	}

	c.flagSet.StringVar(&c.path, "path", "migrations", "Path to generate migrations in.")
	c.flagSet.StringVar(&c.env, "env", "development", "Environment to use for connection.")
	_ = c.flagSet.Parse(args)

	return c.flagSet, nil
}

func (c *sqlGenerator) PopGenerate(ctx context.Context, pwd string, args []string) error {
	name := ""
	if len(args) > 0 {
		name = args[0]
	}

	var (
		atts attrs.Attrs
		err  error
	)
	if len(args) > 1 {
		atts, err = attrs.ParseArgs(args[1:]...)
		if err != nil {
			return err
		}
	}

	run := genny.WetRunner(context.Background())

	// Ensure the generator is as verbose as the old one.
	run.Logger = logger.New(logger.DebugLevel)

	type nameable interface {
		Name() string
	}

	var translator nameable
	db, err := pop.Connect(c.env)
	if err != nil {
		return err
	}

	t := db.Dialect.FizzTranslator()
	if tn, ok := t.(nameable); ok {
		translator = tn
	} else {
		return errors.New("invalid fizz translator")
	}

	var g *genny.Generator
	if len(atts) == 0 {
		g, err = cempty.New(&cempty.Options{
			Name:       name,
			Path:       c.path,
			Type:       "sql",
			Translator: translator,
		})

	} else {
		g, err = ctable.New(&ctable.Options{
			TableName:  name,
			Path:       c.path,
			Type:       "sql",
			Attrs:      atts,
			Translator: t,
		})
	}

	if err != nil {
		return err
	}

	err = run.With(g)
	if err != nil {
		return err
	}

	return run.Run()
}
