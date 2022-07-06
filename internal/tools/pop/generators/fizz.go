package generators

import (
	"context"
	"flag"
	"io"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/pop/v6/genny/fizz/cempty"
	"github.com/gobuffalo/pop/v6/genny/fizz/ctable"
)

var Fizz = &fizzGenerator{}

type fizzGenerator struct {
	flagSet *flag.FlagSet

	path string
}

func (c fizzGenerator) Name() string {
	return "fizz"
}

func (c fizzGenerator) Usage() string {
	return "buffalo db generate fizz <name> [attributes]"
}

func (c fizzGenerator) Aliases() []string {
	return []string{"migration"}
}

func (c fizzGenerator) HelpText() string {
	return "Generates Up/Down migrations for your database using fizz."
}

func (c *fizzGenerator) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(io.Discard)
	}

	c.flagSet.StringVar(&c.path, "path", "migrations", "Path to generate migrations in.")
	_ = c.flagSet.Parse(args)

	return c.flagSet, nil
}

func (c *fizzGenerator) PopGenerate(ctx context.Context, pwd string, args []string) error {
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

	var g *genny.Generator

	if len(atts) == 0 {
		g, err = cempty.New(&cempty.Options{
			Name: name,
			Path: c.path,
			Type: "fizz",
		})
	} else {
		g, err = ctable.New(&ctable.Options{
			TableName: name,
			Path:      c.path,
			Type:      "fizz",
			Attrs:     atts,
		})
	}

	if err != nil {
		return err
	}

	run.With(g)
	return run.Run()
}
