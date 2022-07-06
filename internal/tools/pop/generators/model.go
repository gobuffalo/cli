package generators

import (
	"context"
	"flag"
	"io"
	"os"
	"os/exec"

	"github.com/gobuffalo/attrs"
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/pop/v6/genny/fizz/ctable"
	gmodel "github.com/gobuffalo/pop/v6/genny/model"
)

var Model = &modelGenerator{}

type modelGenerator struct {
	flagSet *flag.FlagSet

	structTag     string
	migrationType string
	modelsPath    string

	skipMigration bool
	env           string
	migrationPath string
}

func (c modelGenerator) Name() string {
	return "model"
}

func (c modelGenerator) Aliases() []string {
	return []string{"m"}
}

func (c modelGenerator) Usage() string {
	return "buffalo db generate model <name> [attributes]"
}

func (c modelGenerator) HelpText() string {
	return "Generates a model for your database"
}

func (c *modelGenerator) ParseFlags(args []string) (*flag.FlagSet, error) {
	if c.flagSet == nil {
		c.flagSet = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
		c.flagSet.Usage = func() {}
		c.flagSet.SetOutput(io.Discard)
	}

	c.flagSet.StringVar(&c.structTag, "struct-tag", "json", "sets the struct tags for model (xml/json/jsonapi)")
	c.flagSet.StringVar(&c.migrationType, "migration-type", "fizz", "sets the type of migration files for model (sql or fizz)")
	c.flagSet.StringVar(&c.modelsPath, "models-path", "models", "the path the model will be created in")
	c.flagSet.StringVar(&c.migrationPath, "migrations-path", "migrations", "the path the migrations will be created in")
	c.flagSet.StringVar(&c.env, "env", "development", "Environment to use for connection.")

	c.flagSet.BoolVar(&c.skipMigration, "skip-migrations", false, "Skip creating a new fizz migration for this model.")

	_ = c.flagSet.Parse(args)

	return c.flagSet, nil
}

func (c *modelGenerator) PopGenerate(ctx context.Context, pwd string, args []string) error {
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

	// Mount models generator
	g, err := gmodel.New(&gmodel.Options{
		Name:                   name,
		Attrs:                  atts,
		Path:                   c.modelsPath,
		Encoding:               c.structTag,
		ForceDefaultID:         true,
		ForceDefaultTimestamps: true,
	})

	if err != nil {
		return err
	}

	run.With(g)

	// format generated go files
	g, err = gogen.Fmt(pwd)
	if err != nil {
		return err
	}

	run.With(g)

	// generated modules may have new dependencies
	if _, err := os.Stat("go.mod"); err == nil {
		g = genny.New()
		g.Command(exec.Command("go", "mod", "tidy"))
		run.With(g)
	}

	// Mount migrations generator
	if c.skipMigration {
		return run.Run()
	}

	var translator fizz.Translator
	if c.migrationType == "sql" {
		db, err := pop.Connect(c.env)
		if err != nil {
			return err
		}
		translator = db.Dialect.FizzTranslator()
	}

	g, err = ctable.New(&ctable.Options{
		TableName:              name,
		Attrs:                  atts,
		Path:                   c.migrationPath,
		Type:                   c.migrationType,
		Translator:             translator,
		ForceDefaultID:         true,
		ForceDefaultTimestamps: true,
	})

	if err != nil {
		return err
	}

	run.With(g)

	return run.Run()
}
