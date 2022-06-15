package pop

import (
	"bytes"
	"context"
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/gobuffalo/pop/v6"
	"github.com/sirupsen/logrus"
)

// SetupDatabase is a before tester that will load the schema
// from the migrations folder if found. This is useful to speed
// tests from running migrations on each iteration.
type SetupDatabase struct {
	forceMigrations bool
}

func (ls SetupDatabase) Name() string {
	return "pop/setup-database"
}

func (ls SetupDatabase) HelpText() string {
	return "Sets the database up before the tests run."
}

func (ls *SetupDatabase) ParseFlags(args []string) (*flag.FlagSet, error) {
	fls := flag.NewFlagSet("pop/setup-database", flag.ContinueOnError)
	fls.BoolVar(&ls.forceMigrations, "force-migrations", false, "force migrations to run")
	_ = fls.Parse(args)

	return fls, nil
}

func (ls *SetupDatabase) BeforeTest(ctx context.Context, pwd string, args []string) error {
	// there's a database
	test, err := pop.Connect("test")
	if err != nil {
		return err
	}

	// drop the test db:
	if err := test.Dialect.DropDB(); err != nil {
		// not an error, since the database will be created in the next step anyway
		logrus.Info("no test database to drop")
	}

	// create the test db:
	err = test.Dialect.CreateDB()
	if err != nil {
		return err
	}

	if ls.forceMigrations {
		fm, err := pop.NewFileMigrator("./migrations", test)
		if err != nil {
			return err
		}

		if err := fm.Up(); err != nil {
			return err
		}

		return nil
	}

	schema := findSchema()
	if schema == nil {
		return nil
	}

	err = test.Dialect.LoadSchema(schema)
	if err != nil {
		return err
	}

	return nil
}

func findSchema() io.Reader {
	if f, err := os.Open(filepath.Join("migrations", "schema.sql")); err == nil {
		return f
	}

	if dev, err := pop.Connect("development"); err == nil {
		schema := &bytes.Buffer{}
		if err = dev.Dialect.DumpSchema(schema); err == nil {
			return schema
		}
	}

	if test, err := pop.Connect("test"); err == nil {
		fm, err := pop.NewFileMigrator("./migrations", test)
		if err != nil {
			return nil
		}

		if err := fm.Up(); err == nil {
			return nil
		}
	}
	return nil
}
