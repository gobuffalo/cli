package generators

import (
	"context"
	"flag"
	"io"
	"path/filepath"

	"github.com/gobuffalo/cli/internal/defaults"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/pop/v6/genny/config"
)

var ConfigGenerator = &configGenerator{}

type configGenerator struct {
	flagSet *flag.FlagSet

	configFile string
	dialect    string
}

func (c configGenerator) Name() string {
	return "config"
}

func (c configGenerator) HelpText() string {
	return "Generates a database.yml file for your project."
}

func (c *configGenerator) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(io.Discard)
	}

	c.flagSet.StringVar(&c.configFile, "config", "database.yml", "The name of the config file to generate.")
	c.flagSet.StringVar(&c.dialect, "type", "postgres", "The dialect to use for the config file.")

	_ = c.flagSet.Parse(args)

	return c.flagSet, nil
}

func (c *configGenerator) Generate(ctx context.Context, pwd string, args []string) error {
	cfgFile := defaults.String(c.configFile, "database.yml")
	run := genny.WetRunner(ctx)

	g, err := config.New(&config.Options{
		Root:     pwd,
		Prefix:   filepath.Base(pwd),
		FileName: cfgFile,
		Dialect:  c.dialect,
	})

	if err != nil {
		return err
	}

	err = run.With(g)
	if err != nil {
		return err
	}

	return run.Run()
}
